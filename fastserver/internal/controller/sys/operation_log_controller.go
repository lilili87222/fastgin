package sys

import (
	"fastgin/config"
	"fastgin/internal/bean"
	"fastgin/internal/controller"
	"fastgin/internal/dao/sys"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// OperationLogController handles operation log-related requests
type OperationLogController struct {
	operationLogRepository sys.OperationLogRepository
}

// NewOperationLogController creates a new OperationLogController
func NewOperationLogController() OperationLogController {
	operationLogRepository := sys.NewOperationLogRepository()
	operationLogController := OperationLogController{operationLogRepository: operationLogRepository}
	return operationLogController
}

// GetOperationLogs retrieves a list of operation logs
// @Summary Get operation log list
// @Description Get a list of operation logs
// @Tags OperationLog
// @Accept json
// @Produce json
// @Param method query string false "Request method"
// @Param path query string false "Request path"
// @Param category query string false "Category"
// @Param creator query string false "Creator"
// @Param pageNum query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /operation_logs [get]
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

// BatchDeleteOperationLogByIds deletes multiple operation logs by their IDs
// @Summary Batch delete operation logs
// @Description Delete multiple operation logs by their IDs
// @Tags OperationLog
// @Accept json
// @Produce json
// @Param operationLogIds body bean.DeleteOperationLogRequest true "Delete operation log request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /operation_logs/batch_delete [delete]
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
