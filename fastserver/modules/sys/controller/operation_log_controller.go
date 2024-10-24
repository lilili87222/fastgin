package controller

import (
	"fastgin/common/httpz"
	"fastgin/modules/sys/service"
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
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/log/index [get]
func (oc *OperationLogController) List(c *gin.Context) {
	params, e := httpz.GetFormData(c)
	if e != nil {
		httpz.BadRequest(c, e.Error())
		return
	}
	data, total, err := oc.logService.Search(httpz.NewSearchRequest(params))
	if err != nil {
		httpz.ServerError(c, "获取操作日志列表失败: "+err.Error())
		return
	}
	httpz.Success(c, gin.H{"data": data, "total": total})
}

// BatchDeleteByIds deletes multiple operation logs by their Ids
// @Summary Batch delete operation logs
// @Description BatchDelete multiple operation logs by their Ids
// @Tags OperationLog
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param operationLogIds body httpz.IdListRequest true "BatchDelete operation log request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/log/index [delete]
func (oc *OperationLogController) BatchDeleteByIds(c *gin.Context) {
	var req httpz.IdListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	err := oc.logService.BatchDelete(req.Ids)
	if err != nil {
		httpz.ServerError(c, "删除日志失败: "+err.Error())
		return
	}
	httpz.Success(c, nil)
}
