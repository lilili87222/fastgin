package storage

type StorageConfig struct {
	StorageType  string             `mapstructure:"storage_type" json:"storage_type" yaml:"storage_type"`
	Local        LocalConfig        `mapstructure:"local" json:"local" yaml:"local"`
	Qiniu        QiniuConfig        `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`
	AliyunOSS    AliyunOSSConfig    `mapstructure:"aliyun-oss" json:"aliyun-oss" yaml:"aliyun-oss"`
	HuaWeiObs    HuaWeiObsConfig    `mapstructure:"hua-wei-obs" json:"hua-wei-obs" yaml:"hua-wei-obs"`
	TencentCOS   TencentCOSConfig   `mapstructure:"tencent-cos" json:"tencent-cos" yaml:"tencent-cos"`
	AwsS3        AwsS3Config        `mapstructure:"aws-s3" json:"aws-s3" yaml:"aws-s3"`
	CloudflareR2 CloudflareR2Config `mapstructure:"cloudflare-r2" json:"cloudflare-r2" yaml:"cloudflare-r2"`
}
