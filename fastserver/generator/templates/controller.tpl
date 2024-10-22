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

// New{{.ModelName}}Controller creates a new {{.ModelName}}Controller
func New{{.ModelName}}Controller() *{{.ModelName}}Controller {
 return &{{.ModelName}}Controller{service: service.New{{.ModelName}}Service()}
}


// Create godoc
// @Summary Create a new {{.LowModelName}} entry
// @Description Create a new {{.LowModelName}} entry
// @Tags {{.ModelName}}
// @Accept json
// @Produce json
// @Param {{.LowModelName}} body model.{{.ModelName}} true "{{.ModelName}} entry"
// @Success 200 {object} model.{{.ModelName}}
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth/{{.LowModelName}}/index [post]
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


// GetByID godoc
// @Summary Get a {{.LowModelName}} entry by ID
// @Description Get a {{.LowModelName}} entry by ID
// @Tags {{.ModelName}}
// @Produce json
// @Param id path int true "{{.ModelName}} ID"
// @Success 200 {object} model.{{.ModelName}}
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth//{{.LowModelName}}/index/{id} [get]
func (ctrl *{{.ModelName}}Controller) GetByID(c *gin.Context) {
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
// @Summary Update a {{.LowModelName}} entry
// @Description Update a {{.LowModelName}} entry
// @Tags {{.ModelName}}
// @Accept json
// @Produce json
// @Param id path int true "{{.ModelName}} ID"
// @Param {{.LowModelName}} body model.{{.ModelName}} true "{{.ModelName}} entry"
// @Success 200 {object} model.{{.ModelName}}
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth//{{.LowModelName}}/index/{id} [patch]
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

// Delete godoc
// @Summary Delete a {{.LowModelName}} entry
// @Description Delete a {{.LowModelName}} entry
// @Tags {{.ModelName}}
// @Produce json
// @Param id path int true "{{.ModelName}} ID"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth//{{.LowModelName}}/index/{id} [delete]
func (ctrl *{{.ModelName}}Controller) Delete(c *gin.Context) {
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
// @Summary List {{.LowModelName}} entries
// @Description List {{.LowModelName}} entries
// @Tags {{.ModelName}}
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth//{{.LowModelName}}/index [get]
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
 httpz.Success(c, gin.H{"data": data, "total": total, "page_num": sr.PageNum, "page_size": sr.PageSize})
}

// DeleteBatch godoc
// @Summary Batch delete {{.LowModelName}} entries
// @Description Batch delete {{.LowModelName}} entries
// @Tags {{.ModelName}}
// @Accept json
// @Produce json
// @Param ids body httpz.IdListRequest true "List of {{.LowModelName}} IDs"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Failure 500 {object} httpz.ResponseBody
// @Router /api/auth//{{.LowModelName}}/index [delete]
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
