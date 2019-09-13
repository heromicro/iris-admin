package ctl

import (
	"github.com/wanhello/iris-admin/internal/app/bll"
	"github.com/wanhello/iris-admin/internal/app/errors"
	"github.com/wanhello/iris-admin/internal/app/schema"
	"github.com/wanhello/iris-admin/internal/app/irisplus"

	"github.com/kataras/iris"

)

// NewRole 创建角色管理控制器
func NewRole(bRole bll.IRole) *Role {
	return &Role{
		RoleBll: bRole,
	}
}

// Role 角色管理
// @Name Role
// @Description 角色管理接口
type Role struct {
	RoleBll bll.IRole
}

// Query 查询数据
func (a *Role) Query(c *iris.Context) {
	switch c.URLParam("q") {
	case "page":
		a.QueryPage(c)
	case "select":
		a.QuerySelect(c)
	default:
		irisplus.ResError(c, errors.ErrUnknownQuery)
	}
}

// QueryPage 查询分页数据
// @Summary 查询分页数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" 1
// @Param pageSize query int true "分页大小" 10
// @Param name query string false "角色名称(模糊查询)"
// @Success 200 []schema.Role "分页查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/roles?q=page
func (a *Role) QueryPage(c *iris.Context) {
	var params schema.RoleQueryParam
	params.LikeName = c.URLParam("name")

	result, err := a.RoleBll.Query(irisplus.NewContext(c), params, schema.RoleQueryOptions{
		PageParam: irisplus.GetPaginationParam(c),
	})
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResPage(c, result.Data, result.PageResult)
}

// QuerySelect 查询选择数据
// @Summary 查询选择数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 []schema.Role "{list:角色列表}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/roles?q=select
func (a *Role) QuerySelect(c *iris.Context) {
	result, err := a.RoleBll.Query(irisplus.NewContext(c), schema.RoleQueryParam{})
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResList(c, result.Data)
}

// Get 查询指定数据
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.Role
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/roles/{id}
func (a *Role) Get(c *iris.Context) {
	item, err := a.RoleBll.Get(irisplus.NewContext(c), c.Param("id"), schema.RoleQueryOptions{
		IncludeMenus: true,
	})
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResSuccess(c, item)
}

// Create 创建数据
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.Role true
// @Success 200 schema.Role
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/roles
func (a *Role) Create(c *iris.Context) {
	var item schema.Role
	if err := irisplus.ParseJSON(c, &item); err != nil {
		irisplus.ResError(c, err)
		return
	}

	item.Creator = irisplus.GetUserID(c)
	nitem, err := a.RoleBll.Create(irisplus.NewContext(c), item)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}

	irisplus.ResSuccess(c, nitem)
}

// Update 更新数据
// @Summary 更新数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Param body body schema.Role true
// @Success 200 schema.Role
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/roles/{id}
func (a *Role) Update(c *iris.Context) {
	var item schema.Role
	if err := irisplus.ParseJSON(c, &item); err != nil {
		irisplus.ResError(c, err)
		return
	}

	nitem, err := a.RoleBll.Update(irisplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResSuccess(c, nitem)
}

// Delete 删除数据
// @Summary 删除数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.HTTPStatus "{status:OK}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router DELETE /api/v1/roles/{id}
func (a *Role) Delete(c *iris.Context) {
	err := a.RoleBll.Delete(irisplus.NewContext(c), c.Param("id"))
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResOK(c)
}
