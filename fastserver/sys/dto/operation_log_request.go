package dto

// 操作日志请求结构体
type OperationLogListRequest struct {
	Username string `json:"UserName" form:"UserName"`
	Ip       string `json:"Ip" form:"Ip"`
	Path     string `json:"Path" form:"Path"`
	Status   int    `json:"Status" form:"Status"`
	PageNum  int    `json:"PageNum" form:"PageNum"`
	PageSize int    `json:"PageSize" form:"PageSize"`
}

// 批量删除操作日志结构体
type DeleteOperationLogRequest struct {
	OperationLogIds []uint `json:"OperationLogIds" form:"OperationLogIds"`
}
