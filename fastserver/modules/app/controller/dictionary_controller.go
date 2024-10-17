package controller

import (
	"fastgin/common/httpz"
	"fastgin/database"
	"fastgin/modules/app/model"
	"fastgin/modules/app/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type DictionaryController struct {
	service *service.DictionaryService
}

func NewDictionaryController() *DictionaryController {
	return &DictionaryController{service: service.NewDictionaryService()}
}

func (ctrl *DictionaryController) Create(c *gin.Context) {
	var entity model.Dictionary
	if err := c.ShouldBindJSON(&entity); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	if err := ctrl.service.Create(&entity); err != nil {
		httpz.ServerError(c, err.Error())
		return
	}
	httpz.Success(c, entity)
}

func (ctrl *DictionaryController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httpz.BadRequest(c, "Invalid ID")
		return
	}
	entity, err := ctrl.service.GetByID(uint(id))
	if err != nil {
		httpz.ServerError(c, err.Error())
		return
	}
	httpz.Success(c, entity)
}

func (ctrl *DictionaryController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httpz.BadRequest(c, "Invalid ID")
		return
	}
	var entity model.Dictionary
	if err := c.ShouldBindJSON(&entity); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	entity.ID = int32(uint(id))
	if err := ctrl.service.Update(&entity); err != nil {
		httpz.ServerError(c, err.Error())
		return
	}
	httpz.Success(c, entity)
}

func (ctrl *DictionaryController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httpz.BadRequest(c, "Invalid ID")
		return
	}
	if err := ctrl.service.Delete(uint(id)); err != nil {
		httpz.ServerError(c, err.Error())
		return
	}
	httpz.Success(c, nil)
}

func (ctrl *DictionaryController) List(c *gin.Context) {
	params, e := httpz.GetFormData(c)
	if e != nil {
		httpz.BadRequest(c, e.Error())
		return
	}
	sr := httpz.NewSearchRequest(params)
	data, total, err := database.SearchTable[model.Dictionary](sr)
	if err != nil {
		httpz.ServerError(c, err.Error())
		return
	}
	httpz.Success(c, gin.H{"Data": data, "Total": total, "PageNum": sr.PageNum, "PageSize": sr.PageSize})
}
func (ctrl *DictionaryController) DeleteBatch(c *gin.Context) {
	var req httpz.IdListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	err := ctrl.service.DeleteBatch(req.Ids)
	if err != nil {
		httpz.ServerError(c, err.Error())
		return
	}
	httpz.Success(c, nil)
}
