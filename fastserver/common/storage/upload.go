package storage

import (
	"gopkg.in/yaml.v3"
	"mime/multipart"
	"os"
)

type IStorage interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

func NewStorage(c *StorageConfig) (IStorage, error) {
	switch c.StorageType {
	case "local":
		return NewLocal(c.Local)
	case "qiniu":
		return NewQiniu(c.Qiniu)
	case "tencent-cos":
		return NewTencentCOS(c.TencentCOS)
	case "aliyun-oss":
		return NewAliyunOSS(c.AliyunOSS)
	case "huawei-obs":
		return NewObs(c.HuaWeiObs)
	case "aws-s3":
		return NewAwsS3(c.AwsS3)
	case "cloudflare-r2":
		return NewCloudflareR2(c.CloudflareR2)
	default:
		return NewLocal(c.Local)
	}
}

var Storage IStorage

func InitStorage(filePath string) {
	data, e := os.ReadFile(filePath)
	if e != nil {
		panic(e)
	}
	s := &StorageConfig{}
	e = yaml.Unmarshal(data, s)
	if e != nil {
		panic(e)
	}
	Storage, e = NewStorage(s)
	if e != nil {
		panic(e)
	}
}
