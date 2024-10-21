package controller

import (
 "fastgin/common/httpz"
 "fastgin/database"
 "{{.Module}}/model"
 "{{.Module}}/service"
 "github.com/gin-gonic/gin"
 "strconv"
)

type {{.ModelName}}Controller struct {
 service *service.{{.ModelName}}Service
}

func New{{.ModelName}}Controller() *{{.ModelName}}Controller {
 return &{{.ModelName}}Controller{service: service.New{{.ModelName}}Service()}
}

func (ctrl *{{.ModelName}}Controller) Create(c *gin.Context) {
 var entity model.{{.ModelName}}
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

func (ctrl *{{.ModelName}}Controller) GetByID(c *gin.Context) {
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

func (ctrl *{{.ModelName}}Controller) Update(c *gin.Context) {
 id, err := strconv.Atoi(c.Param("id"))
 if err != nil {
  httpz.BadRequest(c, "Invalid ID")
  return
 }
 var entity model.{{.ModelName}}
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

func (ctrl *{{.ModelName}}Controller) Delete(c *gin.Context) {
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

func (ctrl *{{.ModelName}}Controller) List(c *gin.Context) {
 params, e := httpz.GetFormData(c)
 if e != nil {
  httpz.BadRequest(c, e.Error())
  return
 }
 sr := httpz.NewSearchRequest(params)
 data, total, err := database.SearchTable[model.{{.ModelName}}](sr)
 if err != nil {
  httpz.ServerError(c, err.Error())
  return
 }
 httpz.Success(c, gin.H{"Data": data, "Total": total, "PageNum": sr.PageNum, "PageSize": sr.PageSize})
}
func (ctrl *{{.ModelName}}Controller) DeleteBatch(c *gin.Context) {
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
