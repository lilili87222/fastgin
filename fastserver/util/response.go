package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseBody struct {
	Code    int            `json:"Code"`
	Message string         `json:"Message"`
	Data    map[string]any `json:"Data"`
}

// 返回前端
func Response(c *gin.Context, httpStatus int, code int, data gin.H, message string) {
	c.JSON(httpStatus, ResponseBody{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// 返回前端-成功
func Success(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusOK, 200, data, message)
}

// 返回前端-失败
func Fail(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusBadRequest, 400, data, message)
}
func GetFormData(c *gin.Context) (map[string]string, error) {
	err := c.Request.ParseForm()
	if err != nil {
		return nil, err
	}
	params := make(map[string]string)
	for key, values := range c.Request.Form {
		// 只取第一个值
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	return params, nil
}
