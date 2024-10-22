package storage

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

type QiniuConfig struct {
	Zone          string `mapstructure:"zone" json:"zone" yaml:"zone"`                                  // 存储区域
	Bucket        string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`                            // 空间名称
	ImgPath       string `mapstructure:"img-path" json:"img-path" yaml:"img-path"`                      // CDN加速域名
	AccessKey     string `mapstructure:"access-key" json:"access-key" yaml:"access-key"`                // 秘钥AK
	SecretKey     string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`                // 秘钥SK
	UseHTTPS      bool   `mapstructure:"use-https" json:"use-https" yaml:"use-https"`                   // 是否使用https
	UseCdnDomains bool   `mapstructure:"use-cdn-domains" json:"use-cdn-domains" yaml:"use-cdn-domains"` // 上传是否使用CDN上传加速
}

type Qiniu struct {
	Config   QiniuConfig
	qiConfig *storage.Config
}

func (q *Qiniu) UploadFile(file *multipart.FileHeader) (string, string, error) {
	putPolicy := storage.PutPolicy{Scope: q.Config.Bucket}
	mac := qbox.NewMac(q.Config.AccessKey, q.Config.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	formUploader := storage.NewFormUploader(q.qiConfig)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}

	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close()                                                  // 创建文件 defer 关闭
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
	if putErr != nil {
		return "", "", errors.New("function formUploader.Put() failed, err:" + putErr.Error())
	}
	return q.Config.ImgPath + "/" + ret.Key, ret.Key, nil
}

func (q *Qiniu) DeleteFile(key string) error {
	mac := qbox.NewMac(q.Config.AccessKey, q.Config.SecretKey)
	bucketManager := storage.NewBucketManager(mac, q.qiConfig)
	if err := bucketManager.Delete(q.Config.Bucket, key); err != nil {
		return errors.New("function bucketManager.Delete() failed, err:" + err.Error())
	}
	return nil
}

func NewQiniu(config QiniuConfig) (*Qiniu, error) {
	cfg := storage.Config{
		UseHTTPS:      config.UseHTTPS,
		UseCdnDomains: config.UseCdnDomains,
	}
	switch config.Zone { // 根据配置文件进行初始化空间对应的机房
	case "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	}
	return &Qiniu{Config: config, qiConfig: &cfg}, nil
}
