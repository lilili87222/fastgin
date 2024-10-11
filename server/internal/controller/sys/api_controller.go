package sys

import (
	"fastgin/config"
	"fastgin/internal/bean"
	"fastgin/internal/controller"
	sys2 "fastgin/internal/dao/sys"
	"fastgin/internal/model/sys"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
)

type ApiController struct {
	ApiRepository sys2.ApiRepository
}

func NewApiController() ApiController {
	apiRepository := sys2.NewApiRepository()
	apiController := ApiController{ApiRepository: apiRepository}
	return apiController
}

// 获取接口列表
func (ac ApiController) GetApis(c *gin.Context) {
	var req bean.ApiListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	// 获取
	apis, total, err := ac.ApiRepository.GetApis(&req)
	if err != nil {
		controller.Fail(c, nil, "获取接口列表失败")
		return
	}
	controller.Success(c, gin.H{
		"apis": apis, "total": total,
	}, "获取接口列表成功")
}

// 获取接口树(按接口Category字段分类)
func (ac ApiController) GetApiTree(c *gin.Context) {
	tree, err := ac.ApiRepository.GetApiTree()
	if err != nil {
		controller.Fail(c, nil, "获取接口树失败")
		return
	}
	controller.Success(c, gin.H{
		"apiTree": tree,
	}, "获取接口树成功")
}

// 创建接口
func (ac ApiController) CreateApi(c *gin.Context) {
	var req bean.CreateApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}

	// 获取当前用户
	ur := sys2.NewUserRepository()
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

	// 创建接口
	err = ac.ApiRepository.CreateApi(&api)
	if err != nil {
		controller.Fail(c, nil, "创建接口失败: "+err.Error())
		return
	}

	controller.Success(c, nil, "创建接口成功")
	return
}

// 更新接口
func (ac ApiController) UpdateApiById(c *gin.Context) {
	var req bean.UpdateApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}

	// 获取路径中的apiId
	apiId, _ := strconv.Atoi(c.Param("apiId"))
	if apiId <= 0 {
		controller.Fail(c, nil, "接口ID不正确")
		return
	}

	// 获取当前用户
	ur := sys2.NewUserRepository()
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

	err = ac.ApiRepository.UpdateApiById(uint(apiId), &api)
	if err != nil {
		controller.Fail(c, nil, "更新接口失败: "+err.Error())
		return
	}

	controller.Success(c, nil, "更新接口成功")
}

// 批量删除接口
func (ac ApiController) BatchDeleteApiByIds(c *gin.Context) {
	var req bean.DeleteApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}

	// 删除接口
	err := ac.ApiRepository.BatchDeleteApiByIds(req.ApiIds)
	if err != nil {
		controller.Fail(c, nil, "删除接口失败: "+err.Error())
		return
	}

	controller.Success(c, nil, "删除接口成功")
}
