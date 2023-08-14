package file

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"transfer/internal/app/pkg/logger"

	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

var StorageSet = wire.NewSet(wire.Struct(new(Storage), "*"))

type Storage struct {
	MinIOClient *minio.Client
}

func (s *Storage) PresignedURL(ctx context.Context, bucketName string, fileName string) (*url.URL, error) {

	reqParams := make(url.Values)
	logger.Infof("get minio file %s url", fileName)

	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	return s.MinIOClient.PresignedHeadObject(ctx, bucketName, fileName, viper.GetViper().GetDuration("MinIO.AuthExpired")*time.Second, reqParams)
}

func (s *Storage) Sign(ctx context.Context, bucketName string, fileName string) string {
	presignedURL, err := s.PresignedURL(ctx, bucketName, fileName)
	if err != nil {
		logger.Errorf("get presignedurl failed. %s", err.Error())
		return ""
	}
	presignDate, _ := time.Parse("20060102T150405Z", presignedURL.Query().Get("X-Amz-Date"))

	return fmt.Sprintf("%d_%s", presignDate.Unix(), presignedURL.Query().Get("X-Amz-Signature"))
}

func (s *Storage) CheckSign(ctx context.Context, bucketName string, fileName string, fileSign string) bool {
	presignedURL, err := s.PresignedURL(ctx, bucketName, fileName)
	if err != nil {
		logger.Errorf("get presignedurl failed. %s", err.Error())
		return false
	}
	q := presignedURL.Query()
	amzDate, _ := strconv.ParseInt(strings.Split(fileSign, "_")[0], 10, 64)
	q.Set("X-Amz-Date", time.Unix(amzDate, 0).UTC().Format("20060102T150405Z"))
	q.Set("X-Amz-Signature", strings.Split(fileSign, "_")[1])

	presignedURL.RawQuery = q.Encode()

	logger.Debugf("presigned url: %s", presignedURL)
	res, err := http.Head(presignedURL.String())
	if err != nil {
		logger.Warnf("presignedurl check failed: %s", err.Error())
		return false
	}

	if res.StatusCode != 200 {
		logger.Warnf("MinIO return %s, %s", res.StatusCode, res.Header.Get("X-Minio-Error-Desc"))
		return false
	}
	return true
}

func (s *Storage) PostObjectFile(ctx context.Context, bucketName string, fileName string, r io.Reader) error {
	// 使用用户名创建 bucket ，每个用户单独使用各自 bucket
	if err := s.CreateBucket(ctx, bucketName); err != nil {
		return err
	}

	if _, err := s.MinIOClient.PutObject(ctx, bucketName, fileName, r, -1, minio.PutObjectOptions{
		// 设置为 50M
		PartSize: 51200000,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Storage) ObjectFile(ctx context.Context, bucketName string, fileName string) *minio.Object {

	if _, err := s.MinIOClient.BucketExists(ctx, bucketName); err != nil {
		logger.Errorf("file is not exist or expired. %s", err.Error())
		return nil
	}

	logger.Debug("start download file from MinIO")
	file, err := s.MinIOClient.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		logger.Errorf("download file from MinIO failed. %s", err.Error())
		return nil
	}
	return file
}

func (s *Storage) CreateBucket(ctx context.Context, bucketName string) error {

	logger.Infof("start create buket %s", bucketName)
	exists, err := s.MinIOClient.BucketExists(ctx, bucketName)
	if err != nil {
		logger.Errorf("check bucket error: %s", err.Error())
		return err
	}
	if !exists {
		if mkerr := s.MinIOClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			logger.Errorf("create bucket error: %s", mkerr.Error())
		}
		logger.Debugf("bucket %s is exist.", bucketName)
	}

	logger.Infof("Successfully created %s", bucketName)
	return nil
}

func (s *Storage) CheckExist(ctx context.Context, bucketName string, fileName string) bool {

	exists, err := s.MinIOClient.BucketExists(ctx, bucketName)
	if err != nil && !exists {
		logger.Errorf("check bucket error: %s", err.Error())
		return false
	}
	for object := range s.MinIOClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{}) {
		if object.Err != nil {
			logger.Errorf("get file list from MinIO failed. %s", object.Err.Error())
		}
		if object.Key == fileName {
			return true
		}
	}
	return false
}
