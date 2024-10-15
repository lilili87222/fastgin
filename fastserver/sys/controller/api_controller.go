package controller

import (
	"fastgin/config"
	"fastgin/sys/dto"
	"fastgin/sys/model"
	"fastgin/sys/service"
	"fastgin/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
)

// ApiController handles API requests
type ApiController struct {
	apiService  *service.ApiService
	userService *service.UserService
}

// NewApiController creates a new ApiController
func NewApiController() *ApiController {
	return &ApiController{apiService: service.NewApiService(), userService: service.NewUserService()}
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
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/apis [get]
func (ac *ApiController) GetApis(c *gin.Context) {
	var req dto.ApiListRequest
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		util.Fail(c, nil, errStr)
		return
	}
	apis, total, err := ac.apiService.GetApis(&req)
	if err != nil {
		util.Fail(c, nil, "获取接口列表失败")
		return
	}
	util.Success(c, gin.H{
		"Apis":  apis,
		"Total": total,
	}, "获取接口列表成功")
}

// GetApiTree retrieves the API tree
// @Summary Get API tree
// @Description Get the API tree categorized by the Category field
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/api/tree [get]
func (ac *ApiController) GetApiTree(c *gin.Context) {
	tree, err := ac.apiService.GetApiTree()
	if err != nil {
		util.Fail(c, nil, "获取接口树失败")
		return
	}
	util.Success(c, gin.H{
		"ApiTree": tree,
	}, "获取接口树成功")
}

// CreateApi creates a new API
// @Summary Create API
// @Description Create a new API
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param api body dto.CreateApiRequest true "Create API request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/api [post]
func (ac *ApiController) CreateApi(c *gin.Context) {
	var req dto.CreateApiRequest
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		util.Fail(c, nil, errStr)
		return
	}
	//ur := service.NewUserService()
	ctxUser, err := ac.userService.GetCurrentUser(c)
	if err != nil {
		util.Fail(c, nil, "获取当前用户信息失败")
		return
	}
	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.UserName,
	}
	err = ac.apiService.CreateApi(&api)
	if err != nil {
		util.Fail(c, nil, "创建接口失败: "+err.Error())
		return
	}
	util.Success(c, nil, "创建接口成功")
}

// UpdateApiById updates an existing API by Id
// @Summary Update API
// @Description Update an existing API by Id
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param apiId path int true "API Id"
// @Param api body dto.UpdateApiRequest true "Update API request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/api/{apiId} [put]
func (ac *ApiController) UpdateApiById(c *gin.Context) {
	var req dto.UpdateApiRequest
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		util.Fail(c, nil, errStr)
		return
	}
	apiId, _ := strconv.Atoi(c.Param("apiId"))
	if apiId <= 0 {
		util.Fail(c, nil, "接口ID不正确")
		return
	}
	//ur := service.NewUserService()
	ctxUser, err := ac.userService.GetCurrentUser(c)
	if err != nil {
		util.Fail(c, nil, "获取当前用户信息失败")
		return
	}
	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.UserName,
	}
	err = ac.apiService.UpdateApiById(uint(apiId), &api)
	if err != nil {
		util.Fail(c, nil, "更新接口失败: "+err.Error())
		return
	}
	util.Success(c, nil, "更新接口成功")
}

// BatchDeleteApiByIds deletes multiple APIs by their Ids
// @Summary Batch delete APIs
// @Description Delete multiple APIs by their Ids
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param apiIds body dto.DeleteApiRequest true "Delete API request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/api/batch_delete [delete]
func (ac *ApiController) BatchDeleteApiByIds(c *gin.Context) {
	var req dto.DeleteApiRequest
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		util.Fail(c, nil, errStr)
		return
	}
	err := ac.apiService.BatchDeleteApiByIds(req.ApiIds)
	if err != nil {
		util.Fail(c, nil, "删除接口失败: "+err.Error())
		return
	}
	util.Success(c, nil, "删除接口成功")
}
