package permession

import (
	"time"
	"transfer/api/types/schema"
	"transfer/internal/app/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/wire"
)

var PermessionSet = wire.NewSet(wire.Struct(new(PermessionRepo), "*"))

type Permession struct {
	gorm.Model
	Fileid     string `gorm:"fileid;size:64;not null"`
	Username   string `gorm:"username;size:32;unique_index:username_shorturl"`
	Shorturl   int64  `gorm:"shorturl;unique_index:username_shorturl;default:1"` // 默认值设置为1，未指定默认值为0插入报错
	Expiredate time.Time
}

type PermessionRepo struct {
	DB *gorm.DB
}

func (p *PermessionRepo) Create(permession *Permession) error {
	if err := p.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "shorturl"}, {Name: "username"}},
		DoUpdates: clause.AssignmentColumns([]string{"fileid", "expiredate", "deleted_at"}),
	}).Create(&permession).Error; err != nil {
		logger.Errorf("create auth failed. %s", err.Error())
		return err
	}
	return nil
}

func (p *PermessionRepo) Delete(permession *Permession) error {
	if err := p.DB.Where("fileid=?", permession.Fileid).Delete(permession).Error; err != nil {
		logger.Warnf("delete file permession failed. %s", err.Error())
		return err
	}
	return nil
}

func (p *PermessionRepo) Info(permession *Permession) *[]schema.FilePermession {
	filePermession := &[]schema.FilePermession{}
	nowTime := time.Now()
	if err := p.DB.Model(permession).Where(permession).Where("expiredate>?", nowTime).Scan(filePermession).Error; err != nil {
		logger.Debugf("get file permession failed or permession is not exist. %s", err.Error())
		return nil
	}
	logger.Debugf("file permession %s: %s ", permession.Fileid, filePermession)
	return filePermession
}

func (p *PermessionRepo) Update(permession *Permession) bool {
	result := p.DB.Model(permession).Where("fileid=? and username=?", permession.Fileid, permession.Username).Updates(permession)
	if result.Error != nil || result.RowsAffected == 0 {
		logger.Warnf("update permession error: %s", result.Error.Error())
		return false
	}
	return true
}
