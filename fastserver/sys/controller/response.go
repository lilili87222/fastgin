package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseBody struct {
	Code int            `json:"code"`
	Msg  string         `json:"message"`
	Data map[string]any `json:"data"`
}

// 返回前端
func Response(c *gin.Context, httpStatus int, code int, data gin.H, message string) {
	c.JSON(httpStatus, gin.H{"code": code, "data": data, "message": message})
}

// 返回前端-成功
func Success(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusOK, 200, data, message)
}

// 返回前端-失败
func Fail(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusBadRequest, 400, data, message)
}