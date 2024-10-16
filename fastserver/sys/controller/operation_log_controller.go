package controller

import (
	"fastgin/sys/dto"
	"fastgin/sys/service"
	"fastgin/util"
	"github.com/gin-gonic/gin"
)

// OperationLogController handles operation log-related requests
type OperationLogController struct {
	logService *service.OperationLogService
}

// NewOperationLogController creates a new OperationLogController
func NewOperationLogController() *OperationLogController {
	return &OperationLogController{logService: service.NewLogService()}
}

// List retrieves a list of operation logs
// @Summary Get operation log list
// @Description Get a list of operation logs
// @Tags OperationLog
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param path query string false "Request path"
// @Param user_name query string false "user_name"
// @Param ip query string false "ip"
// @Param status query int false "status"
// @Param PageNum query int false "Page number"
// @Param PageSize query int false "Page size"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/log/index [get]
func (oc *OperationLogController) List(c *gin.Context) {
	params, e := util.GetFormData(c)
	if e != nil {
		util.Fail(c, nil, e.Error())
		return
	}
	data, total, err := oc.logService.Search(dto.NewSearchRequest(params))
	if err != nil {
		util.Fail(c, nil, "获取操作日志列表失败: "+err.Error())
		return
	}
	util.Success(c, gin.H{"Data": data, "Total": total}, "获取操作日志列表成功")
}

// BatchDeleteByIds deletes multiple operation logs by their Ids
// @Summary Batch delete operation logs
// @Description Delete multiple operation logs by their Ids
// @Tags OperationLog
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param operationLogIds body dto.IdListRequest true "Delete operation log request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/log/index [delete]
func (oc *OperationLogController) BatchDeleteByIds(c *gin.Context) {
	var req dto.IdListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	err := oc.logService.BatchDelete(req.Ids)
	if err != nil {
		util.Fail(c, nil, "删除日志失败: "+err.Error())
		return
	}
	util.Success(c, nil, "删除日志成功")
}
