package controller

import (
	"fastgin/common/httpz"
	"fastgin/config"
	"fastgin/database"
	"fastgin/modules/sys/dto"
	"fastgin/modules/sys/model"
	"fastgin/modules/sys/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
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

// List retrieves a list of APIs
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
// @Param PageNum query int false "Page number"
// @Param PageSize query int false "Page size"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/api/index [get]
func (ac *ApiController) List(c *gin.Context) {
	params, e := httpz.GetFormData(c)
	if e != nil {
		httpz.BadRequest(c, e.Error())
		return
	}
	data, total, err := database.SearchTable[model.Api](httpz.NewSearchRequest(params))
	if err != nil {
		httpz.ServerError(c, "获取接口列表失败: "+err.Error())
		return
	}
	httpz.Success(c, gin.H{"Data": data, "Total": total})
}

// GetApiTree retrieves the API tree
// @Summary Get API tree
// @Description Get the API tree categorized by the Category field
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/api/tree [get]
func (ac *ApiController) GetApiTree(c *gin.Context) {
	tree, err := ac.apiService.GetApiTree()
	if err != nil {
		httpz.ServerError(c, "获取接口树失败")
		return
	}
	httpz.Success(c, tree)
}

// Create creates a new API
// @Summary Create API
// @Description Create a new API
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param api body dto.CreateApiRequest true "Create API request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/api/index [post]
func (ac *ApiController) Create(c *gin.Context) {
	var req dto.CreateApiRequest
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		httpz.BadRequest(c, errStr)
		return
	}
	//ur := service.NewUserService()
	ctxUser, err := ac.userService.GetCurrentUser(c)
	if err != nil {
		httpz.ServerError(c, "获取当前用户信息失败")
		return
	}
	api := model.Api{
		//Method:   req.Method,
		//Path:     req.Path,
		//Category: req.Category,
		//Desc:     req.Desc,
		Creator: ctxUser.UserName,
	}
	copier.Copy(&api, &req)
	//api.Creator = ctxUser.UserName
	err = ac.apiService.CreateApi(&api)
	if err != nil {
		httpz.ServerError(c, "创建接口失败: "+err.Error())
		return
	}
	httpz.Success(c, nil)
}

// Update updates an existing API by ID
// @Summary Update API
// @Description Update an existing API by ID
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param apiId path int true "API ID"
// @Param api body dto.CreateApiRequest true "Update API request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/api/index/{apiId} [put]
func (ac *ApiController) Update(c *gin.Context) {
	//type UpdateApiRequest struct {
	//	Method   string `json:"Method" form:"Method" validate:"min=1,max=20"`
	//	Path     string `json:"Path" form:"Path" validate:"min=1,max=100"`
	//	Category string `json:"Category" form:"Category" validate:"min=1,max=50"`
	//	Desc     string `json:"Desc" form:"Desc" validate:"min=0,max=100"`
	//}
	var req dto.CreateApiRequest
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		httpz.BadRequest(c, errStr)
		return
	}
	apiId, _ := strconv.Atoi(c.Param("apiId"))
	if apiId <= 0 {
		httpz.BadRequest(c, "接口ID不正确")
		return
	}
	ctxUser, err := ac.userService.GetCurrentUser(c)
	if err != nil {
		httpz.ServerError(c, "获取当前用户信息失败")
		return
	}
	api := model.Api{
		//Method:   req.Method,
		//Path:     req.Path,
		//Category: req.Category,
		//Desc:     req.Desc,
		Creator: ctxUser.UserName,
	}
	copier.Copy(&api, &req)
	api.ID = uint64(apiId)
	err = ac.apiService.UpdateApiById(&api)
	if err != nil {
		httpz.ServerError(c, "更新接口失败: "+err.Error())
		return
	}
	httpz.Success(c, nil)
}

// BatchDelete deletes multiple APIs by their Ids
// @Summary Batch delete APIs
// @Description BatchDelete multiple APIs by their Ids
// @Tags API
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param apiIds body httpz.IdListRequest true "BatchDelete API request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/api/index [delete]
func (ac *ApiController) BatchDelete(c *gin.Context) {
	var req httpz.IdListRequest
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	err := database.DeleteByIds[model.Api](req.Ids)
	if err != nil {
		httpz.ServerError(c, "删除接口失败: "+err.Error())
		return
	}
	httpz.Success(c, nil)
}
