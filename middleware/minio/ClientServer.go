package minio

import (
	"github.com/minio/minio-go/v6"
	"log"
	"strconv"
	"strings"
	"tiktok_demo/config"
	"tiktok_demo/util"
)

type Minio struct {
	MinioClient  *minio.Client
	Endpoint     string
	Port         string
	VideoBuckets string
	PicBuckets   string
}

var client Minio

func InitMinio() {
	conf := config.GetConfig()
	endpoint := conf.Minio.Host
	port := conf.Minio.Port
	endpoint = endpoint + ":" + port
	rootUser := conf.Minio.RootUser
	rootPassword := conf.Minio.RootPassword
	videoBucket := conf.Minio.VideoBuckets
	picBucket := conf.Minio.PicBuckets

	// Minio init
	minioClient, err := minio.New(endpoint, rootUser, rootPassword, false)
	if err != nil {
		panic("failed to init minio, err:" + err.Error())
	}
	log.Printf("Init Minio Client succeed")
	CreatBucket(minioClient, videoBucket)
	CreatBucket(minioClient, picBucket)
	client = Minio{minioClient, endpoint, port, videoBucket, picBucket}
}

func CreatBucket(m *minio.Client, bucket string) {
	found, err := m.BucketExists(bucket)
	if err != nil {
		log.Println("bucketExists: ", bucket, err.Error())
	}
	if !found {
		err := m.MakeBucket(bucket, "us-east-1")
		if err != nil {
			log.Println("MakeBucket failed: ", bucket, err.Error())
		}
	}
	policy := `{"Version": "2012-10-17",
				"Statement": 
					[{
						"Action":["s3:GetObject"],
						"Effect": "Allow",
						"Principal": {"AWS": ["*"]},
						"Resource": ["arn:aws:s3:::` + bucket + `/*"],
						"Sid": ""
					}]
				}`
	err = m.SetBucketPolicy(bucket, policy)
	if err != nil {
		log.Println("SetBucketPolicy err: ", bucket, err.Error())
	}
}

func GetMinio() Minio {
	return client
}

func (m *Minio) UploadFile(fileType, file, userId string) (string, error) {
	var fileName strings.Builder
	var contentType, Suffix, bucket string
	if fileType == "video" {
		contentType = "video/mp4"
		Suffix = ".mp4"
		bucket = m.VideoBuckets
	} else {
		contentType = "image/jpeg"
		Suffix = ".jpg"
		bucket = m.PicBuckets
	}
	fileName.WriteString(userId)
	fileName.WriteString("_")
	fileName.WriteString(strconv.FormatInt(util.GetCurrentTimeMillisecond(), 10))
	fileName.WriteString(Suffix)
	n, err := m.MinioClient.FPutObject(bucket, fileName.String(), file, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Println("upload file error: ", err.Error())
		return "", err
	}
	log.Println("upload file success, fileName:", n, fileName, " bytes")
	url := "http:" + "//" + m.Endpoint + "/" + bucket + "/" + fileName.String()
	return url, nil
}
