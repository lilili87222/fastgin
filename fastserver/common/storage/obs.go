package storage

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/pkg/errors"
	"mime/multipart"
)

type HuaWeiObsConfig struct {
	Path      string `mapstructure:"path" json:"path" yaml:"path"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Endpoint  string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKey string `mapstructure:"access-key" json:"access-key" yaml:"access-key"`
	SecretKey string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
}

type Obs struct {
	config HuaWeiObsConfig
	client *obs.ObsClient
}

func NewObs(config HuaWeiObsConfig) (*Obs, error) {
	client, e := obs.New(config.AccessKey, config.SecretKey, config.Endpoint)
	if e != nil {
		return nil, e
	}
	return &Obs{config: config, client: client}, nil
}

func (o *Obs) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// var open multipart.File
	open, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer open.Close()
	filename := file.Filename
	input := &obs.PutObjectInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket: o.config.Bucket,
				Key:    filename,
			},
			HttpHeader: obs.HttpHeader{
				ContentType: file.Header.Get("content-type"),
			},
		},
		Body: open,
	}

	_, err = o.client.PutObject(input)
	if err != nil {
		return "", "", errors.Wrap(err, "文件上传失败!")
	}
	filepath := o.config.Path + "/" + filename
	return filepath, filename, err
}

func (o *Obs) DeleteFile(key string) error {
	input := &obs.DeleteObjectInput{
		Bucket: o.config.Bucket,
		Key:    key,
	}
	output, err := o.client.DeleteObject(input)
	if err != nil {
		return errors.Wrapf(err, "删除对象(%s)失败!, output: %v", key, output)
	}
	return nil
}
