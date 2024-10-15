package controller

import (
	"fastgin/config"
	"fastgin/sys/dto"
	"fastgin/sys/service"
	"fastgin/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// OperationLogController handles operation log-related requests
type OperationLogController struct {
	logDao *service.OperationLogService
}

// NewOperationLogController creates a new OperationLogController
func NewOperationLogController() *OperationLogController {
	return &OperationLogController{logDao: service.NewLogService()}
}

// GetOperationLogs retrieves a list of operation logs
// @Summary Get operation log list
// @Description Get a list of operation logs
// @Tags OperationLog
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param method query string false "Request method"
// @Param path query string false "Request path"
// @Param category query string false "Category"
// @Param creator query string false "Creator"
// @Param pageNum query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/operation_logs [get]
func (oc *OperationLogController) GetOperationLogs(c *gin.Context) {
	var req dto.OperationLogListRequest
	// 绑定参数
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		util.Fail(c, nil, errStr)
		return
	}
	// 获取
	logs, total, err := oc.logDao.GetOperationLogs(&req)
	if err != nil {
		util.Fail(c, nil, "获取操作日志列表失败: "+err.Error())
		return
	}
	util.Success(c, gin.H{"logs": logs, "total": total}, "获取操作日志列表成功")
}

// BatchDeleteOperationLogByIds deletes multiple operation logs by their IDs
// @Summary Batch delete operation logs
// @Description Delete multiple operation logs by their IDs
// @Tags OperationLog
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param operationLogIds body dto.DeleteOperationLogRequest true "Delete operation log request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/operation_logs/batch_delete [delete]
func (oc *OperationLogController) BatchDeleteOperationLogByIds(c *gin.Context) {
	var req dto.DeleteOperationLogRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		util.Fail(c, nil, errStr)
		return
	}

	// 删除接口
	err := oc.logDao.BatchDeleteOperationLogByIds(req.OperationLogIds)
	if err != nil {
		util.Fail(c, nil, "删除日志失败: "+err.Error())
		return
	}

	util.Success(c, nil, "删除日志成功")
}
