package minio

import (
	"context"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

// 判断存储桶是否存在
func existBulck(bulkName string) bool {
	bool_num, err := minioClient.BucketExists(context.Background(), bulkName)
	if err != nil {
		log.Fatalf("is exist bulk eror " + err.Error())
		return false
	}
	return bool_num
}

// CreateBucket 创建桶
func CreateBucket(bucketName string) error {
	if len(bucketName) <= 0 {
		log.Fatal("bucketName invalid")
	}

	location := "beijing"
	ctx := context.Background()

	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("bucket %s already exists", bucketName)
			return nil
		} else {
			return err
		}
	} else {
		log.Printf("bucket %s create successfully", bucketName)
	}
	return nil
}

// UploadFile 上传文件（提供reader）至 minio
func UploadFile(bucketName string, objectName string, reader io.Reader, contentType string, objectSize int64) error {
	ctx := context.Background()
	n, err := minioClient.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Printf("upload %s of size %d failed, %s", bucketName, objectSize, err)
		return err
	}
	log.Printf("upload %s of bytes %d successfully", objectName, n.Size)
	return nil
}

// GetFileUrl 从 minio 获取文件Url
func GetFileUrl(bucketName string, fileName string, expires time.Duration) (*url.URL, error) {
	ctx := context.Background()
	reqParams := make(url.Values)
	if expires <= 0 {
		expires = time.Second * 60 * 60 * 24
	}
	presignedUrl, err := minioClient.PresignedGetObject(ctx, bucketName, fileName, expires, reqParams)
	if err != nil {
		log.Printf("get url of file %s from bucket %s failed, %s", fileName, bucketName, err)
		return nil, err
	}
	return presignedUrl, nil
}
