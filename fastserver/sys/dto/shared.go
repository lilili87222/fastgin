package dto

import "strconv"

type IdListRequest struct {
	Ids []uint `json:"Ids" form:"Ids"`
}
type SearchRequest struct {
	PageNum   int
	PageSize  int
	KeyValues map[string]string
}

func NewSearchRequest(params map[string]string) *SearchRequest {
	req := &SearchRequest{}
	if pageNum, err := strconv.Atoi(params["PageNum"]); err == nil {
		req.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(params["PageSize"]); err == nil {
		req.PageSize = pageSize
	}
	req.KeyValues = params
	return req
}
