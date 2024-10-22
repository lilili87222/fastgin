package storage

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
	"time"
)

type AliyunOSSConfig struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyId     string `mapstructure:"access-key-id" json:"access-key-id" yaml:"access-key-id"`
	AccessKeySecret string `mapstructure:"access-key-secret" json:"access-key-secret" yaml:"access-key-secret"`
	BucketName      string `mapstructure:"bucket-name" json:"bucket-name" yaml:"bucket-name"`
	BucketUrl       string `mapstructure:"bucket-url" json:"bucket-url" yaml:"bucket-url"`
	BasePath        string `mapstructure:"base-path" json:"base-path" yaml:"base-path"`
}
type AliyunOSS struct {
	Config AliyunOSSConfig
	bucket *oss.Bucket
}

func (oss *AliyunOSS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// 读取本地文件。
	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() Failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭
	// 上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
	// yunFileTmpPath := filepath.Join("uploads", time.Now().Format("2006-01-02")) + "/" + file.Filename
	yunFileTmpPath := oss.Config.BasePath + "/" + "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + file.Filename
	// 上传文件流。
	err := oss.bucket.PutObject(yunFileTmpPath, f)
	if err != nil {
		return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}
	return oss.Config.BucketUrl + "/" + yunFileTmpPath, yunFileTmpPath, nil
}

func (oss *AliyunOSS) DeleteFile(key string) error {
	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
	return oss.bucket.DeleteObject(key)
}

func NewAliyunOSS(config AliyunOSSConfig) (*AliyunOSS, error) {
	// 创建OSSClient实例。
	client, err := oss.New(config.Endpoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, err
	}
	return &AliyunOSS{
		Config: config,
		bucket: bucket,
	}, nil
}
