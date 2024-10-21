package httpz

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 返回前端
func Response(c *gin.Context, httpStatus int, data any, message string) {
	c.JSON(httpStatus, ResponseBody{
		Code:    httpStatus,
		Message: message,
		Data:    data,
	})
}

// 返回前端-成功
func Success(c *gin.Context, data any) {
	Response(c, http.StatusOK, data, "Operation is successful")
}
func BadRequest(c *gin.Context, message string) {
	Response(c, http.StatusBadRequest, nil, message)
}

// 返回前端-失败
func ServerError(c *gin.Context, message string) {
	Response(c, http.StatusInternalServerError, nil, message)
}
