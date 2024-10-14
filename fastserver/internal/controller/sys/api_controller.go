package sys

import (
	"fastgin/config"
	"fastgin/internal/controller"
	sys2 "fastgin/internal/dao/sys"
	"fastgin/internal/model/sys"
	"fastgin/internal/model/sys/request"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
)

// ApiController handles API requests
type ApiController struct {
	ApiDao sys2.ApiDao
}

// NewApiController creates a new ApiController
func NewApiController() ApiController {
	apiDao := sys2.NewApiDao()
	apiController := ApiController{ApiDao: apiDao}
	return apiController
}

// GetApis retrieves a list of APIs
// @Summary Get API list
// @Description Get a list of APIs
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param method query string false "Request method"
// @Param path query string false "Request path"
// @Param category query string false "Category"
// @Param creator query string false "Creator"
// @Param pageNum query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/apis [get]
func (ac ApiController) GetApis(c *gin.Context) {
	var req request.ApiListRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	apis, total, err := ac.ApiDao.GetApis(&req)
	if err != nil {
		controller.Fail(c, nil, "获取接口列表失败")
		return
	}
	controller.Success(c, gin.H{
		"apis":  apis,
		"total": total,
	}, "获取接口列表成功")
}

// GetApiTree retrieves the API tree
// @Summary Get API tree
// @Description Get the API tree categorized by the Category field
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/api/tree [get]
func (ac ApiController) GetApiTree(c *gin.Context) {
	tree, err := ac.ApiDao.GetApiTree()
	if err != nil {
		controller.Fail(c, nil, "获取接口树失败")
		return
	}
	controller.Success(c, gin.H{
		"apiTree": tree,
	}, "获取接口树成功")
}

// CreateApi creates a new API
// @Summary Create API
// @Description Create a new API
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param api body request.CreateApiRequest true "Create API request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/api [post]
func (ac ApiController) CreateApi(c *gin.Context) {
	var req request.CreateApiRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	ur := sys2.NewUserDao()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		controller.Fail(c, nil, "获取当前用户信息失败")
		return
	}
	api := sys.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}
	err = ac.ApiDao.CreateApi(&api)
	if err != nil {
		controller.Fail(c, nil, "创建接口失败: "+err.Error())
		return
	}
	controller.Success(c, nil, "创建接口成功")
}

// UpdateApiById updates an existing API by ID
// @Summary Update API
// @Description Update an existing API by ID
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param apiId path int true "API ID"
// @Param api body request.UpdateApiRequest true "Update API request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/api/{apiId} [put]
func (ac ApiController) UpdateApiById(c *gin.Context) {
	var req request.UpdateApiRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	apiId, _ := strconv.Atoi(c.Param("apiId"))
	if apiId <= 0 {
		controller.Fail(c, nil, "接口ID不正确")
		return
	}
	ur := sys2.NewUserDao()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		controller.Fail(c, nil, "获取当前用户信息失败")
		return
	}
	api := sys.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}
	err = ac.ApiDao.UpdateApiById(uint(apiId), &api)
	if err != nil {
		controller.Fail(c, nil, "更新接口失败: "+err.Error())
		return
	}
	controller.Success(c, nil, "更新接口成功")
}

// BatchDeleteApiByIds deletes multiple APIs by their IDs
// @Summary Batch delete APIs
// @Description Delete multiple APIs by their IDs
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param apiIds body request.DeleteApiRequest true "Delete API request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/api/batch_delete [delete]
func (ac ApiController) BatchDeleteApiByIds(c *gin.Context) {
	var req request.DeleteApiRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	err := ac.ApiDao.BatchDeleteApiByIds(req.ApiIds)
	if err != nil {
		controller.Fail(c, nil, "删除接口失败: "+err.Error())
		return
	}
	controller.Success(c, nil, "删除接口成功")
}
