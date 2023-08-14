package app

import (
	"fmt"
	"transfer/internal/app/model/file"
	"transfer/internal/app/model/permession"
	"transfer/internal/app/model/user"
	"transfer/internal/app/pkg/logger"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	var db *gorm.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", viper.GetViper().GetString("MySQL.User"), viper.GetViper().GetString("MySQL.Password"),
		viper.GetViper().GetString("MySQL.Host"), viper.GetViper().GetInt("MySQL.Port"), viper.GetViper().GetString("MySQL.DB"))

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if viper.GetViper().GetBool("Global.Debug") {
		db.AutoMigrate(&file.File{})
		db.AutoMigrate(&user.User{})
		db.AutoMigrate(&user.Role{})
		db.AutoMigrate(&permession.Permession{})
		db = db.Debug()
	}
	if err != nil {
		logger.Fatalf("%s", err)
		return nil
	}
	return db
}
