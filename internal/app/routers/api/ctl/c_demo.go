package ctl

import (
	"github.com/wanhello/iris-admin/internal/app/bll"
	"github.com/wanhello/iris-admin/internal/app/errors"
	"github.com/wanhello/iris-admin/internal/app/schema"

	"github.com/wanhello/iris-admin/internal/app/irisplus"
	
	"github.com/wanhello/iris-admin/pkg/util"

	"github.com/kataras/iris"

)

// NewDemo 创建demo控制器
func NewDemo(bDemo bll.IDemo) *Demo {
	return &Demo{
		DemoBll: bDemo,
	}
}

// Demo demo
// @Name Demo
// @Description 示例接口
type Demo struct {
	DemoBll bll.IDemo
}

// Query 查询数据
func (a *Demo) Query(c iris.Context) {
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
// @Param code query string false "编号"
// @Param name query string false "名称"
// @Param status query int false "状态(1:启用 2:停用)"
// @Success 200 []schema.Demo "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/demos?q=page
func (a *Demo) QueryPage(c iris.Context) {
	var params schema.DemoQueryParam
	params.LikeCode = c.URLParam("code")
	params.LikeName = c.URLParam("name")
	params.Status = util.S(c.URLParam("status")).DefaultInt(0)

	result, err := a.DemoBll.Query(irisplus.NewContext(c), params, schema.DemoQueryOptions{
		PageParam: irisplus.GetPaginationParam(c),
	})
	if err != nil {
		irisplus.ResError(c, err)
		return
	}

	irisplus.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.Demo
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/demos/{id}
func (a *Demo) Get(c iris.Context) {
	item, err := a.DemoBll.Get(irisplus.NewContext(c), c.URLParam("id"))
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResSuccess(c, item)
}

// Create 创建数据
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.Demo true
// @Success 200 schema.Demo
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/demos
func (a *Demo) Create(c iris.Context) {
	var item schema.Demo
	if err := irisplus.ParseJSON(c, &item); err != nil {
		irisplus.ResError(c, err)
		return
	}

	item.Creator = irisplus.GetUserID(c)
	nitem, err := a.DemoBll.Create(irisplus.NewContext(c), item)
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
// @Param body body schema.Demo true
// @Success 200 schema.Demo
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/demos/{id}
func (a *Demo) Update(c iris.Context) {
	var item schema.Demo
	if err := irisplus.ParseJSON(c, &item); err != nil {
		irisplus.ResError(c, err)
		return
	}

	nitem, err := a.DemoBll.Update(irisplus.NewContext(c), c.URLParam("id"), item)
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
// @Router DELETE /api/v1/demos/{id}
func (a *Demo) Delete(c iris.Context) {
	err := a.DemoBll.Delete(irisplus.NewContext(c), c.URLParam("id"))
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
// @Router PATCH /api/v1/demos/{id}/enable
func (a *Demo) Enable(c iris.Context) {
	err := a.DemoBll.UpdateStatus(irisplus.NewContext(c), c.URLParam("id"), 1)
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
// @Router PATCH /api/v1/demos/{id}/disable
func (a *Demo) Disable(c iris.Context) {
	err := a.DemoBll.UpdateStatus(irisplus.NewContext(c), c.URLParam("id"), 2)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResOK(c)
}

