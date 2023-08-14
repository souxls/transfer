package file

import (
	"context"
	"io"
	"time"
	"transfer/api/types/schema"
	fileModel "transfer/internal/app/model/file"
	permessionModel "transfer/internal/app/model/permession"
	"transfer/internal/app/pkg/logger"
	mURL "transfer/internal/app/pkg/url"
	permessionService "transfer/internal/app/service/permession"
	fileStorage "transfer/internal/app/storage/file"

	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

var ServiceSet = wire.NewSet(wire.Struct(new(Service), "*"))

type Service struct {
	FileRepo          *fileModel.FileRepo
	PermessionRepo    *permessionModel.PermessionRepo
	PermessionSrvRepo *permessionService.Service
	FileStorage       *fileStorage.Storage
}

func (s *Service) Create(ctx context.Context, file *fileModel.File, r io.Reader) error {

	if err := s.FileStorage.PostObjectFile(ctx, file.Bucketname, file.Filename, r); err != nil {
		logger.Errorf("upload file to MinIO failed. %s", err.Error())
		return err
	}

	if err := s.FileRepo.Create(file); err != nil {
		logger.Errorf("save file infomation to db failed. %s", err.Error())
		return err
	}
	userPermession := &permessionModel.Permession{
		Fileid:     file.Fileid,
		Username:   file.Owner,
		Shorturl:   file.ID,
		Expiredate: time.Now().Add(time.Minute * viper.GetViper().GetDuration("MinIO.UserExpired")),
	}
	if err := s.PermessionSrvRepo.Create(ctx, userPermession); err != nil {
		logger.Errorf("set auth for %s failed. %s", userPermession.Fileid, err.Error())
		return err
	}

	return nil

}

func (s *Service) Delete(ctx context.Context, file *fileModel.File) error {
	if err := s.FileRepo.Delete(file); err != nil {
		return err
	}
	userPermession := &permessionModel.Permession{
		Fileid: file.Fileid,
	}
	if err := s.PermessionRepo.Delete(userPermession); err != nil {
		return err
	}
	return nil
}

func (s *Service) QueryShow(ctx context.Context, file *fileModel.File, pp schema.PaginationParam) (*schema.FileQueryResult, error) {
	return s.FileRepo.Query(ctx, file, pp)
}

func (s *Service) ForDownload(ctx context.Context, fileListType string, filePermession *[]schema.FilePermession) schema.Files {
	var files schema.Files

	for _, v := range *filePermession {
		fileInfo := s.FileRepo.Info(&fileModel.File{
			Fileid: v.Fileid,
		})
		if fileInfo != nil {
			files = append(files, &schema.FileInfo{
				Fileid:     fileInfo.Fileid,
				Filename:   fileInfo.Filename,
				Owner:      fileInfo.Owner,
				Createtime: fileInfo.Createtime,
				Expired:    v.Expiredate,
				URL:        mURL.ShortUrlencode(fileInfo.ID),
			})
		}
	}
	return files
}

func (s *Service) Info(ctx context.Context, file *fileModel.File) *schema.File {
	info := s.FileRepo.Info(file)

	if info != nil {
		return info
	}
	return nil
}

func (s *Service) ByFileid(ctx context.Context, file *fileModel.File) *schema.File {
	fileInfo := s.FileRepo.ByFileid(file)

	if fileInfo != nil {
		return fileInfo
	}
	return nil
}

func (s *Service) Sign(ctx context.Context, file *schema.File) *schema.FileSign {
	fileSign := s.FileStorage.Sign(ctx, file.Bucketname, file.Filename)
	if fileSign != "" {
		return &schema.FileSign{
			Filename:  file.Filename,
			AuthParam: fileSign,
		}
	}
	return nil
}

func (s *Service) CheckSign(ctx context.Context, file *schema.File, fileSign string) bool {
	return s.FileStorage.CheckSign(ctx, file.Bucketname, file.Filename, fileSign)
}

func (s *Service) Download(ctx context.Context, bucketName string, fileName string) *minio.Object {
	// 使用下载object文件
	return s.FileStorage.ObjectFile(ctx, bucketName, fileName)
}
