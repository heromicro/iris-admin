package ctl

import (
	"github.com/wanhello/iris-admin/internal/app/bll"
	"github.com/wanhello/iris-admin/internal/app/errors"
	"github.com/wanhello/iris-admin/internal/app/schema"
	"github.com/wanhello/iris-admin/internal/app/irisplus"
	"github.com/wanhello/iris-admin/pkg/util"

	"github.com/kataras/iris"

)

// NewMenu 创建菜单管理控制器
func NewMenu(bMenu bll.IMenu) *Menu {
	return &Menu{
		MenuBll: bMenu,
	}
}

// Menu 菜单管理
// @Name Menu
// @Description 菜单管理接口
type Menu struct {
	MenuBll bll.IMenu
}

// Query 查询数据
func (a *Menu) Query(c iris.Context) {
	switch c.URLParam("q") {
	case "page":
		a.QueryPage(c)
	case "tree":
		a.QueryTree(c)
	default:
		irisplus.ResError(c, errors.ErrUnknownQuery)
	}
}

// QueryPage 查询分页数据
// @Summary 查询分页数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" 1
// @Param pageSize query int true "分页大小" 10
// @Param name query string false "名称"
// @Param hidden query int false "隐藏菜单(0:不隐藏 1:隐藏)"
// @Param parent_id query string false "父级ID"
// @Success 200 []schema.Menu "分页查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/menus?q=page
func (a *Menu) QueryPage(c iris.Context) {
	params := schema.MenuQueryParam{
		LikeName: c.URLParam("name"),
	}

	if v := c.URLParam("parent_id"); v != "" {
		params.ParentID = &v
	}

	if v := c.URLParam("hidden"); v != "" {
		if hidden := util.S(v).DefaultInt(0); hidden > -1 {
			params.Hidden = &hidden
		}
	}

	result, err := a.MenuBll.Query(irisplus.NewContext(c), params, schema.MenuQueryOptions{
		PageParam: irisplus.GetPaginationParam(c),
	})
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResPage(c, result.Data, result.PageResult)
}

// QueryTree 查询菜单树
// @Summary 查询菜单树
// @Param Authorization header string false "Bearer 用户令牌"
// @Param include_actions query int false "是否包含动作数据(1是)"
// @Param include_resources query int false "是否包含资源数据(1是)"
// @Success 200 option.Interface "查询结果：{list:菜单树}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/menus?q=tree
func (a *Menu) QueryTree(c iris.Context) {
	result, err := a.MenuBll.Query(irisplus.NewContext(c), schema.MenuQueryParam{}, schema.MenuQueryOptions{
		IncludeActions:   c.URLParam("include_actions") == "1",
		IncludeResources: c.URLParam("include_resources") == "1",
	})
	if err != nil {
		irisplus.ResError(c, err)
		return
	}

	irisplus.ResList(c, result.Data.ToTrees().ToTree())
}

// Get 查询指定数据
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.Menu
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/menus/{id}
func (a *Menu) Get(c iris.Context) {
	item, err := a.MenuBll.Get(irisplus.NewContext(c), c.URLParam("id"), schema.MenuQueryOptions{
		IncludeActions:   true,
		IncludeResources: true,
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
// @Param body body schema.Menu true
// @Success 200 schema.Menu
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/menus
func (a *Menu) Create(c iris.Context) {
	var item schema.Menu
	if err := irisplus.ParseJSON(c, &item); err != nil {
		irisplus.ResError(c, err)
		return
	}

	item.Creator = irisplus.GetUserID(c)
	nitem, err := a.MenuBll.Create(irisplus.NewContext(c), item)
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
// @Param body body schema.Menu true
// @Success 200 schema.Menu
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/menus/{id}
func (a *Menu) Update(c iris.Context) {
	var item schema.Menu
	if err := irisplus.ParseJSON(c, &item); err != nil {
		irisplus.ResError(c, err)
		return
	}

	nitem, err := a.MenuBll.Update(irisplus.NewContext(c), c.URLParam("id"), item)
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
// @Router DELETE /api/v1/menus/{id}
func (a *Menu) Delete(c iris.Context) {
	err := a.MenuBll.Delete(irisplus.NewContext(c), c.URLParam("id"))
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResOK(c)
}


