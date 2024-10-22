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

// NewDictionaryController creates a new DictionaryController
func NewDictionaryController() *DictionaryController {
 return &DictionaryController{service: service.NewDictionaryService()}
}


// Create godoc
// @Summary Create a new dictionary entry
// @Description Create a new dictionary entry
// @Tags Dictionary
// @Accept json
// @Produce json
// @Param dictionary body model.Dictionary true "Dictionary entry"
// @Success 200 {object} model.Dictionary
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth/dictionary/index [post]
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


// GetByID godoc
// @Summary Get a dictionary entry by ID
// @Description Get a dictionary entry by ID
// @Tags Dictionary
// @Produce json
// @Param id path int true "Dictionary ID"
// @Success 200 {object} model.Dictionary
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth//dictionary/index/{id} [get]
func (ctrl *DictionaryController) GetByID(c *gin.Context) {
 id, err := strconv.Atoi(c.Param("id"))
 if err != nil {
  httpz.BadRequest(c, "Invalid ID")
  return
 }
 entity, err := ctrl.service.GetByID(uint64(id))
 if err != nil {
  httpz.ServerError(c, err.Error())
  return
 }
 httpz.Success(c, entity)
}

// Update godoc
// @Summary Update a dictionary entry
// @Description Update a dictionary entry
// @Tags Dictionary
// @Accept json
// @Produce json
// @Param id path int true "Dictionary ID"
// @Param dictionary body model.Dictionary true "Dictionary entry"
// @Success 200 {object} model.Dictionary
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth//dictionary/index/{id} [patch]
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

// Delete godoc
// @Summary Delete a dictionary entry
// @Description Delete a dictionary entry
// @Tags Dictionary
// @Produce json
// @Param id path int true "Dictionary ID"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth//dictionary/index/{id} [delete]
func (ctrl *DictionaryController) Delete(c *gin.Context) {
 id, err := strconv.Atoi(c.Param("id"))
 if err != nil {
  httpz.BadRequest(c, "Invalid ID")
  return
 }
 if err := ctrl.service.Delete(uint64(id)); err != nil {
  httpz.ServerError(c, err.Error())
  return
 }
 httpz.Success(c, nil)
}

// List godoc
// @Summary List dictionary entries
// @Description List dictionary entries
// @Tags Dictionary
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth//dictionary/index [get]
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
 httpz.Success(c, gin.H{"data": data, "total": total, "page_num": sr.PageNum, "page_size": sr.PageSize})
}

// DeleteBatch godoc
// @Summary Batch delete dictionary entries
// @Description Batch delete dictionary entries
// @Tags Dictionary
// @Accept json
// @Produce json
// @Param ids body httpz.IdListRequest true "List of dictionary IDs"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth//dictionary/index [delete]
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
