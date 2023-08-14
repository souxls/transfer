package file

import (
	"context"
	"time"
	"transfer/api/types/errors"
	"transfer/api/types/schema"
	"transfer/internal/app/model/utils"
	"transfer/internal/app/pkg/logger"

	"github.com/google/wire"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var FileSet = wire.NewSet(wire.Struct(new(FileRepo), "*"))

type File struct {
	gorm.Model
	ID         int64     `gorm:"id;primaryKey;autIncrement"`
	Fileid     string    `gorm:"fileid;size:64;not null"`
	Filename   string    `gorm:"filname;size:128;not null;unique_index:filename_owmer"`
	Bucketname string    `gorm:"bucketname;size:32;not null"`
	Owner      string    `gorm:"owner; size:32;not null;unique_index:filename_owmer"`
	Createtime time.Time `gorm:"createtime"`
	Expired    time.Time `gorm:"expired"`
}

type FileRepo struct {
	DB *gorm.DB
}

func (f *FileRepo) Query(ctx context.Context, file *File, pp schema.PaginationParam) (*schema.FileQueryResult, error) {
	files := schema.Files{}
	nowTime := time.Now()
	pr, err := utils.WrapPageQuery(ctx, f.DB.Model(file).Where("owner=? and expired>?", file.Owner, nowTime), pp, &files)
	if err != nil {
		logger.Warnf("file list is empty or get list error. %s", err.Error())
		return nil, err
	}

	qr := &schema.FileQueryResult{
		PageResult: pr,
		PageData:   files,
	}
	return qr, nil

}

func (f *FileRepo) Create(file *File) error {
	if result := f.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "filename"}, {Name: "owner"}},
		UpdateAll: true,
	}).Create(&file); result.Error != nil || result.RowsAffected == 0 {
		logger.Warnf("create file failed. %s", result.Error.Error())
		return errors.ErrSaveFaild
	}
	return nil
}

func (f *FileRepo) Delete(file *File) error {

	if result := f.DB.Where(file).Delete(file); result.Error != nil || result.RowsAffected == 0 {
		logger.Warnf("delete file error. %s", result.Error.Error())
		return errors.ErrDelFaild
	}
	return nil

}

func (f *FileRepo) ByUnexpired(file *File) *[]schema.FileInfo {

	files := &[]schema.FileInfo{}
	nowTime := time.Now()
	if result := f.DB.Model(file).Where("owner=? and expired>?", file.Owner, nowTime).Scan(files); result.Error != nil || result.RowsAffected == 0 {
		logger.Debugf("get file info failed or file is not exist. %s", result.Error.Error())
		return nil
	}
	return files
}

func (f *FileRepo) Info(file *File) *schema.File {

	info := &schema.File{}
	if result := f.DB.Model(file).Where(file).Scan(info); result.Error != nil || result.RowsAffected == 0 {
		logger.Debugf("get file info failed or file is not exist. %s", result.Error.Error())
		return nil
	}
	logger.Debugf("file info: %s", info)
	return info
}

func (f *FileRepo) ByFileid(file *File) *schema.File {

	info := &schema.File{}
	if result := f.DB.Model(file).Where("fileid=?", file.Fileid).Scan(info); result.Error != nil || result.RowsAffected == 0 {
		logger.Errorf("get file info failed or file is not exist. %s", result.Error.Error())
		return nil
	}
	logger.Debugf("file info: %s", info)
	return info
}
