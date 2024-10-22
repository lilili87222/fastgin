package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

type TencentCOSConfig struct {
	Bucket     string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Region     string `mapstructure:"region" json:"region" yaml:"region"`
	SecretID   string `mapstructure:"secret-id" json:"secret-id" yaml:"secret-id"`
	SecretKey  string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	BaseURL    string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	PathPrefix string `mapstructure:"path-prefix" json:"path-prefix" yaml:"path-prefix"`
}

type TencentCOS struct {
	Config TencentCOSConfig
	client *cos.Client
}

func (t *TencentCOS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)

	_, err := t.client.Object.Put(context.Background(), t.Config.PathPrefix+"/"+fileKey, f, nil)
	if err != nil {
		panic(err)
	}
	return t.Config.BaseURL + "/" + t.Config.PathPrefix + "/" + fileKey, fileKey, nil
}

func (t *TencentCOS) DeleteFile(key string) error {
	name := t.Config.PathPrefix + "/" + key
	_, err := t.client.Object.Delete(context.Background(), name)
	if err != nil {
		return errors.New("function bucketManager.Delete() failed, err:" + err.Error())
	}
	return nil
}

func NewTencentCOS(config TencentCOSConfig) (*TencentCOS, error) {
	urlStr, _ := url.Parse("https://" + config.Bucket + ".cos." + config.Region + ".myqcloud.com")
	baseURL := &cos.BaseURL{BucketURL: urlStr}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
	})
	return &TencentCOS{
		Config: config,
		client: client,
	}, nil
}
