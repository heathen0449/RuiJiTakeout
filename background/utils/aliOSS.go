package utils

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
)

var Bucket *oss.Bucket

func BucketInit() error {
	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		return fmt.Errorf("oss client error: %v", err)
	}
	Bucket, err = client.Bucket("heathen-project")
	if err != nil {
		return fmt.Errorf("oss bucket error: %v", err)
	}
	return nil
}

func UploadImage(image multipart.File, imageName string) (string, error) {

	err := Bucket.PutObject(imageName, image)
	if err != nil {
		fmt.Println("上传函数失败 " + imageName)
		fmt.Println("Error:", err)
		return "", err
	}
	return fmt.Sprintf("https://%s.%s/%s", BucketName, Endpoint, imageName), nil
}
