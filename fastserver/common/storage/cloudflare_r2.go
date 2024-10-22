package storage

import (
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type CloudflareR2Config struct {
	Bucket          string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	BaseURL         string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	Path            string `mapstructure:"path" json:"path" yaml:"path"`
	AccountID       string `mapstructure:"account-id" json:"account-id" yaml:"account-id"`
	AccessKeyID     string `mapstructure:"access-key-id" json:"access-key-id" yaml:"access-key-id"`
	SecretAccessKey string `mapstructure:"secret-access-key" json:"secret-access-key" yaml:"secret-access-key"`
}

type CloudflareR2 struct {
	Config  CloudflareR2Config
	session *session.Session
}

func (c *CloudflareR2) UploadFile(file *multipart.FileHeader) (fileUrl string, fileName string, err error) {
	client := s3manager.NewUploader(c.session)
	fileKey := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	fileName = fmt.Sprintf("%s/%s", c.Config.Path, fileKey)
	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	input := &s3manager.UploadInput{
		Bucket: aws.String(c.Config.Bucket),
		Key:    aws.String(fileName),
		Body:   f,
	}

	_, err = client.Upload(input)
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%s/%s", c.Config.BaseURL,
			fileName),
		fileKey,
		nil
}

func (c *CloudflareR2) DeleteFile(key string) error {
	svc := s3.New(c.session)
	filename := c.Config.Path + "/" + key
	bucket := c.Config.Bucket

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return errors.New("function svc.DeleteObject() failed, err:" + err.Error())
	}

	return svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
}

func NewCloudflareR2(config CloudflareR2Config) (*CloudflareR2, error) {
	endpoint := fmt.Sprintf("%s.r2.cloudflarestorage.com", config.AccountID)
	s := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("auto"),
		Endpoint: aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyID,
			config.SecretAccessKey,
			"",
		),
	}))
	return &CloudflareR2{
		Config:  config,
		session: s,
	}, nil
}
