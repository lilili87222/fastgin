package model

//	type Api struct {
//		gorm.Model
//		Method   string `gorm:"type:varchar(20);comment:'请求方式'" json:"method"`
//		Path     string `gorm:"type:varchar(100);comment:'访问路径'" json:"path"`
//		Category string `gorm:"type:varchar(50);comment:'所属类别'" json:"category"`
//		Desc     string `gorm:"type:varchar(100);comment:'说明'" json:"desc"`
//		Creator  string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
//	}
type Api struct {
	Model
	Method   string `gorm:"type:varchar(20);comment:'请求方式'" `
	Path     string `gorm:"type:varchar(100);comment:'访问路径'" `
	Category string `gorm:"type:varchar(50);comment:'所属类别'" `
	Desc     string `gorm:"type:varchar(100);comment:'说明'" `
	Creator  string `gorm:"type:varchar(20);comment:'创建人'" `
}

func (*Api) TableName() string {
	return "sys_api"
}
