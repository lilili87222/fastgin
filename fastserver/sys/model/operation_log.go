package model

import (
	"gorm.io/gorm"
	"time"
)

//	type OperationLog struct {
//		gorm.Model
//		UserName   string    `gorm:"type:varchar(20);comment:'用户登录名'" json:"username"`
//		Ip         string    `gorm:"type:varchar(20);comment:'Ip地址'" json:"ip"`
//		IpLocation string    `gorm:"type:varchar(20);comment:'Ip所在地'" json:"ipLocation"`
//		Method     string    `gorm:"type:varchar(20);comment:'请求方式'" json:"method"`
//		Path       string    `gorm:"type:varchar(100);comment:'访问路径'" json:"path"`
//		Desc       string    `gorm:"type:varchar(100);comment:'说明'" json:"desc"`
//		Status     int       `gorm:"type:int(4);comment:'响应状态码'" json:"status"`
//		StartTime  time.Time `gorm:"type:datetime(3);comment:'发起时间'" json:"startTime"`
//		TimeCost   int64     `gorm:"type:int(6);comment:'请求耗时(ms)'" json:"timeCost"`
//		UserAgent  string    `gorm:"type:varchar(20);comment:'浏览器标识'" json:"userAgent"`
//	}
type OperationLog struct {
	gorm.Model
	UserName   string    `gorm:"type:varchar(20);comment:'用户登录名'" `
	Ip         string    `gorm:"type:varchar(20);comment:'Ip地址'" `
	IpLocation string    `gorm:"type:varchar(20);comment:'Ip所在地'" `
	Method     string    `gorm:"type:varchar(20);comment:'请求方式'" `
	Path       string    `gorm:"type:varchar(100);comment:'访问路径'" `
	Desc       string    `gorm:"type:varchar(100);comment:'说明'" `
	Status     int       `gorm:"type:int(4);comment:'响应状态码'" `
	StartTime  time.Time `gorm:"type:datetime(3);comment:'发起时间'" `
	TimeCost   int64     `gorm:"type:int(6);comment:'请求耗时(ms)'" `
	UserAgent  string    `gorm:"type:varchar(20);comment:'浏览器标识'"`
}

func (*OperationLog) TableName() string {
	return "sys_operation_log"
}
