package httpz

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type SearchRequest struct {
	PageNum   int
	PageSize  int
	KeyValues map[string]string
}
type IdListRequest struct {
	Ids []uint64 `json:"ids" form:"ids"`
}

func NewSearchRequest(params map[string]string) *SearchRequest {
	req := &SearchRequest{}
	if pageNum, err := strconv.Atoi(params["page_num"]); err == nil {
		req.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(params["page_size"]); err == nil {
		req.PageSize = pageSize
	}
	delete(params, "page_num")
	delete(params, "page_size")
	req.KeyValues = params
	return req
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
