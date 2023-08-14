package app

import (
	"os"
	"transfer/internal/app/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

func InitMinIO() *minio.Client {

	accessKeyID := viper.GetViper().GetString("MinIO.AccessKeyID")
	secretAccessKey := viper.GetViper().GetString("MinIO.SecretAccessKey")
	endPoint := viper.GetViper().GetString("MinIO.EndPoint")
	useSSL := viper.GetViper().GetBool("MinIO.UseSSL")

	minClient, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL})

	if err != nil {
		logger.Errorf("MinIO init error.", err.Error())
		return nil
	}

	if viper.GetViper().GetBool("Global.Debug") {
		minClient.TraceOn(os.Stdout)
	}

	return minClient
}
