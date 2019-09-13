package ctl

import (
	"strings"

	"github.com/wanhello/iris-admin/internal/app/bll"
	"github.com/wanhello/iris-admin/internal/app/errors"
	"github.com/wanhello/iris-admin/internal/app/schema"
	"github.com/wanhello/iris-admin/internal/app/irisplus"

	"github.com/wanhello/iris-admin/pkg/util"

	"github.com/kataras/iris"

)

// NewUser 创建用户管理控制器
func NewUser(bUser bll.IUser) *User {
	return &User{
		UserBll: bUser,
	}
}

// User 用户管理
// @Name User
// @Description 用户管理接口
type User struct {
	UserBll bll.IUser
}

// Query 查询数据
func (a *User) Query(c *iris.Context) {
	switch c.URLParam("q") {
	case "page":
		a.QueryPage(c)
	default:
		irisplus.ResError(c, errors.ErrUnknownQuery)
	}
}

// QueryPage 查询分页数据
// @Summary 查询分页数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" 1
// @Param pageSize query int true "分页大小" 10
// @Param user_name query string false "用户名(模糊查询)"
// @Param real_name query string false "真实姓名(模糊查询)"
// @Param role_ids query string false "角色ID(多个以英文逗号分隔)"
// @Param status query int false "状态(1:启用 2:停用)"
// @Success 200 []schema.UserShow "分页查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/users?q=page
func (a *User) QueryPage(c *iris.Context) {
	var params schema.UserQueryParam
	params.LikeUserName = c.URLParam("user_name")
	params.LikeRealName = c.URLParam("real_name")
	if v := util.S(c.URLParam("status")).DefaultInt(0); v > 0 {
		params.Status = v
	}

	if v := c.URLParam("role_ids"); v != "" {
		params.RoleIDs = strings.Split(v, ",")
	}

	result, err := a.UserBll.QueryShow(irisplus.NewContext(c), params, schema.UserQueryOptions{
		IncludeRoles: true,
		PageParam:    irisplus.GetPaginationParam(c),
	})
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
// Get 查询指定数据
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.User
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/users/{id}
func (a *User) Get(c *iris.Context) {
	item, err := a.UserBll.Get(irisplus.NewContext(c), c.Param("id"), schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResSuccess(c, item.CleanSecure())
}

// Create 创建数据
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.User true
// @Success 200 schema.User
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/users
func (a *User) Create(c *iris.Context) {
	var item schema.User
	if err := irisplus.ParseJSON(c, &item); err != nil {
		irisplus.ResError(c, err)
		return
	}

	item.Creator = irisplus.GetUserID(c)
	nitem, err := a.UserBll.Create(irisplus.NewContext(c), item)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResSuccess(c, nitem.CleanSecure())
}

// Update 更新数据
// @Summary 更新数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Param body body schema.User true
// @Success 200 schema.User
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/users/{id}
func (a *User) Update(c *iris.Context) {
	var item schema.User
	if err := irisplus.ParseJSON(c, &item); err != nil {
		irisplus.ResError(c, err)
		return
	}

	nitem, err := a.UserBll.Update(irisplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResSuccess(c, nitem.CleanSecure())
}

// Delete 删除数据
// @Summary 删除数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.HTTPStatus "{status:OK}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router DELETE /api/v1/users/{id}
func (a *User) Delete(c *iris.Context) {
	err := a.UserBll.Delete(irisplus.NewContext(c), c.Param("id"))
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResOK(c)
}

// Enable 启用数据
// @Summary 启用数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.HTTPStatus "{status:OK}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router PATCH /api/v1/users/{id}/enable
func (a *User) Enable(c *iris.Context) {
	err := a.UserBll.UpdateStatus(irisplus.NewContext(c), c.Param("id"), 1)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResOK(c)
}

// Disable 禁用数据
// @Summary 禁用数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.HTTPStatus "{status:OK}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router PATCH /api/v1/users/{id}/disable
func (a *User) Disable(c *iris.Context) {
	err := a.UserBll.UpdateStatus(irisplus.NewContext(c), c.Param("id"), 2)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResOK(c)
}
