package user

import (
	"context"
	"log"
	"transfer/internal/app/pkg/logger"

	"github.com/minio/madmin-go/v2"
	"github.com/spf13/viper"
)

func InitMadmin(ctx context.Context) *madmin.AdminClient {
	endPoint := viper.GetViper().GetString("MinIO.EndPoint")
	accessKeyID := viper.GetViper().GetString("MinIO.AccessKeyID")
	secretAccessKey := viper.GetViper().GetString("MinIO.SecretAccessKey")
	useSSL := viper.GetBool("MinIO.UseSSL")

	madmClnt, err := madmin.New(endPoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Println("min admin init failed", err)
	}
	return madmClnt
}

func CreateUser(ctx context.Context, userName string, password string) error {
	madmClnt := InitMadmin(ctx)
	if CheckUserExist(ctx, userName) {
		return nil
	}
	if err := madmClnt.AddUser(ctx, userName, password); err != nil {
		logger.Errorf("user add failed %s. %s", userName, err)
		return err
	}
	log.Println("user add success", userName)
	return nil

}

func CheckUserExist(ctx context.Context, userName string) bool {
	madmClnt := InitMadmin(ctx)
	users, err := madmClnt.ListUsers(ctx)
	if err != nil {
		logger.Fatal("获取用户列表失败")
		return false
	}
	if user, ok := users[userName]; !ok {

		logger.Warnf("用户不存在", user)
		return false
	}

	return true
}
