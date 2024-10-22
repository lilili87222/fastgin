package controller

import (
	"fastgin/boost/config"
	"fastgin/common/httpz"
	"fastgin/modules/sys/service"
	"github.com/gin-gonic/gin"
	"os"
)

type SystemController struct {
	systemService *service.SystemService
}

func NewSystemController() *SystemController {
	return &SystemController{systemService: &service.SystemService{}}
}

// GetSystemInformation godoc
// @Summary 获取系统信息
// @Description 获取系统的详细信息
// @Tags 系统
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/system/info [get]
func (oc *SystemController) GetSystemInformation(c *gin.Context) {
	//service := service.SystemService{}
	httpz.Success(c, oc.systemService.GetSystemInformation())
}

// GetStopServer godoc
// @Summary 停止服务器
// @Description 停止服务器的运行
// @Tags 系统
// @Param Authorization header string true "Bearer token"
// @Success 200 {string} string "停止服务成功"
// @Router /api/auth/system/stop [get]
func (oc *SystemController) GetStopServer(c *gin.Context) {
	config.Log.Info("停止服务")
	os.Exit(0)
	httpz.Success(c, nil)
}

// RestartServer godoc
// @Summary 重启服务器
// @Description 重启服务器的运行
// @Tags 系统
// @Param Authorization header string true "Bearer token"
// @Success 200 {string} string "重启服务成功"
// @Router /api/auth/system/restart [get]
func (oc *SystemController) RestartServer(c *gin.Context) {
	//service := service.SystemService{}
	e := oc.systemService.Restart()
	config.Log.Info("重启服务")
	if e != nil {
		config.Log.Info("重启服务失败 " + e.Error())
	} else {
		config.Log.Info(c, nil, "重启服务成功")
	}
	httpz.Success(c, nil)
}
