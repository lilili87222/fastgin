package sys

import (
	"fastgin/config"
	"fastgin/internal/bean"
	"fastgin/internal/controller"
	"fastgin/internal/dao/sys"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OperationLogController struct {
	operationLogRepository sys.OperationLogRepository
}

func NewOperationLogController() OperationLogController {
	operationLogRepository := sys.NewOperationLogRepository()
	operationLogController := OperationLogController{operationLogRepository: operationLogRepository}
	return operationLogController
}

// 获取操作日志列表
func (oc OperationLogController) GetOperationLogs(c *gin.Context) {
	var req bean.OperationLogListRequest
	// 绑定参数
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	// 获取
	logs, total, err := oc.operationLogRepository.GetOperationLogs(&req)
	if err != nil {
		controller.Fail(c, nil, "获取操作日志列表失败: "+err.Error())
		return
	}
	controller.Success(c, gin.H{"logs": logs, "total": total}, "获取操作日志列表成功")
}

// 批量删除操作日志
func (oc OperationLogController) BatchDeleteOperationLogByIds(c *gin.Context) {
	var req bean.DeleteOperationLogRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}

	// 删除接口
	err := oc.operationLogRepository.BatchDeleteOperationLogByIds(req.OperationLogIds)
	if err != nil {
		controller.Fail(c, nil, "删除日志失败: "+err.Error())
		return
	}

	controller.Success(c, nil, "删除日志成功")
}
