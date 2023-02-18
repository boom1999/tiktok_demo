package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"tiktok_demo/config"
)

var minioClient *minio.Client

// Minio 对象存储初始化
func InitMinio() {
	var (
		Conf                   = config.GetConfig()
		MinioHost              = Conf.Minio.Host
		MinioPort              = Conf.Minio.Port
		MinioUsername          = Conf.Minio.RootUser
		MinioPassword          = Conf.Minio.RootPassword
		MinioVideoBucketName   = Conf.Minio.VideoBuckets
		MinioPictureBucketName = Conf.Minio.PicBuckets
	)

	client, err := minio.New(MinioHost+":"+MinioPort, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioUsername, MinioPassword, ""),
		Secure: false,
	})
	if err != nil {
		panic("failed to connect minio, err:" + err.Error())
	}
	minioClient = client
	if !existBulck(MinioVideoBucketName) {
		if err = CreateBucket(MinioVideoBucketName); err != nil {
			panic("minio client init CreateBucket failed:" + err.Error())
		}
	}
	if !existBulck(MinioPictureBucketName) {
		if err = CreateBucket(MinioPictureBucketName); err != nil {
			panic("minio client init CreateBucket failed:" + err.Error())
		}
	}
}
