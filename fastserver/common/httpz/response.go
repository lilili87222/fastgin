package httpz

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseBody struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
	Data    any    `json:"Data"`
}

//func SuccessBody(data any) ResponseBody {
//	return ResponseBody{
//		Code:    http.StatusOK,
//		Message: "Operation is successful",
//		Data:    data,
//	}
//}
//func BadRequestBody(message string) ResponseBody {
//	return ResponseBody{
//		Code:    http.StatusBadRequest,
//		Message: message,
//	}
//}
//func ServerErrorBody(message string) ResponseBody {
//	return ResponseBody{
//		Code:    http.StatusInternalServerError,
//		Message: message,
//	}
//}

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
