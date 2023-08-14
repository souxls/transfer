package user

import (
	"transfer/api/types/schema"
	"transfer/internal/app/pkg/logger"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var UserSet = wire.NewSet(wire.Struct(new(UserRepo), "*"))

type User struct {
	gorm.Model
	Username string `gorm:"username;size:32;unique;not null"`
	Realname string `gorm:"realname;size:32;default:null"` // 真实姓名
	Password string `gorm:"password;size:32 not null"`     // 密码
	Roleid   int    `gorm:"roleid;default:2"`              // 角色: 1. admin 2.user
	Role     Role   `gorm:"foreignKey:Roleid"`
}

type Role struct {
	ID          int    `gorm:"id"`
	Roleid      int    `gorm:"roleid"`
	Rolename    string `gorm:"rolename;size:32"`
	Description string `gorm:"description;size:64"`
}

type UserRepo struct {
	DB *gorm.DB
}

func (u *UserRepo) CheckExist(user *User) bool {
	var userInfo *schema.UserInfo
	if err := u.DB.Model(user).Where("username=?", user.Username).Scan(&userInfo).Error; err != nil {
		logger.Debugf("get user error or user is not exist. %s", err.Error())
		return false
	}
	logger.Debugf("user info: %s", userInfo.Username)
	return userInfo.ID != 0
}

func (u *UserRepo) Create(user *User) error {
	if err := u.DB.Create(&user).Error; err != nil {
		logger.Errorf("create user failed. %s", err.Error())
		return err
	}
	return nil
}

func (u *UserRepo) Info(user *User) error {

	if err := u.DB.Where("username=?", user.Username).Find(&user).Error; err != nil {
		logger.Warnf("get user info failed or user is not exist. %s", err.Error())
		return err
	}
	return nil
}

func (u *UserRepo) Role(roleid int) string {
	var role *Role
	if result := u.DB.Where("roleid=?", roleid).Scan(role); result.Error != nil || result.RowsAffected == 0 {
		logger.Warnf("get role name for user %s error. %s", result.Error.Error())
		return "user"
	}
	return role.Rolename
}
