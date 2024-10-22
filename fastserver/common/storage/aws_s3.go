package storage

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AwsS3Config struct {
	Bucket           string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Region           string `mapstructure:"region" json:"region" yaml:"region"`
	Endpoint         string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	SecretID         string `mapstructure:"secret-id" json:"secret-id" yaml:"secret-id"`
	SecretKey        string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	BaseURL          string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	PathPrefix       string `mapstructure:"path-prefix" json:"path-prefix" yaml:"path-prefix"`
	S3ForcePathStyle bool   `mapstructure:"s3-force-path-style" json:"s3-force-path-style" yaml:"s3-force-path-style"`
	DisableSSL       bool   `mapstructure:"disable-ssl" json:"disable-ssl" yaml:"disable-ssl"`
}

type AwsS3 struct {
	awsS3Config AwsS3Config
	session     *session.Session
}

func (s *AwsS3) UploadFile(file *multipart.FileHeader) (string, string, error) {
	uploader := s3manager.NewUploader(s.session)
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	filename := s.awsS3Config.PathPrefix + "/" + fileKey
	f, openError := file.Open()
	if openError != nil {
		return "", "", openError
	}
	defer f.Close()
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.awsS3Config.Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		return "", "", err
	}
	return s.awsS3Config.BaseURL + "/" + filename, fileKey, nil
}

func (s *AwsS3) DeleteFile(key string) error {
	filename := s.awsS3Config.PathPrefix + "/" + key
	bucket := s.awsS3Config.Bucket
	svc := s3.New(s.session)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return err
	}
	return svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
}

func NewAwsS3(config AwsS3Config) (*AwsS3, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(config.Region),
		Endpoint:         aws.String(config.Endpoint), //minio在这里设置地址,可以兼容
		S3ForcePathStyle: aws.Bool(config.S3ForcePathStyle),
		DisableSSL:       aws.Bool(config.DisableSSL),
		Credentials: credentials.NewStaticCredentials(
			config.SecretID,
			config.SecretKey,
			"",
		),
	})
	return &AwsS3{
		awsS3Config: config,
		session:     sess,
	}, nil
}
