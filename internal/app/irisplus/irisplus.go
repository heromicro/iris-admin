package irisplus

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	icontext "github.com/wanhello/iris-admin/internal/app/context"

	"github.com/wanhello/iris-admin/internal/app/errors"
	"github.com/wanhello/iris-admin/internal/app/schema"
	"github.com/wanhello/iris-admin/pkg/logger"
	"github.com/wanhello/iris-admin/pkg/util"

	"github.com/kataras/iris"

)

type HandlerFunc func(iris.Context)

// 定义上下文中的键
const (
	prefix = "irisadmin"
	// UserIDKey 存储上下文中的键(用户ID)
	UserIDKey = prefix + "/user_id"
	// TraceIDKey 存储上下文中的键(跟踪ID)
	TraceIDKey = prefix + "/trace_id"
	// ResBodyKey 存储上下文中的键(响应Body数据)
	ResBodyKey = prefix + "/res_body"
)

// NewContext 封装上线文入口
func NewContext(c iris.Context) context.Context {
	parent := context.Background()

	if v := GetTraceID(c); v != "" {
		parent = icontext.NewTraceID(parent, v)
		parent = logger.NewTraceIDContext(parent, GetTraceID(c))
	}

	if v := GetUserID(c); v != "" {
		parent = icontext.NewUserID(parent, v)
		parent = logger.NewUserIDContext(parent, v)
	}

	return parent
}

// GetToken 获取用户令牌
func GetToken(c iris.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// GetPageIndex 获取分页的页索引
func GetPageIndex(c iris.Context) int {
	defaultVal := 1
	if v := c.URLParam("current"); v != "" {
		if iv := util.S(v).DefaultInt(defaultVal); iv > 0 {
			return iv
		}
	}
	return defaultVal
}


// GetPageSize 获取分页的页大小(最大50)
func GetPageSize(c iris.Context) int {
	defaultVal := 10
	if v := c.Params().Get("pageSize"); v != "" {
		if iv := util.S(v).DefaultInt(defaultVal); iv > 0 {
			if iv > 50 {
				iv = 50
			}
			return iv
		}
	}
	return defaultVal
}


// GetPaginationParam 获取分页查询参数
func GetPaginationParam(c iris.Context) *schema.PaginationParam {
	return &schema.PaginationParam{
		PageIndex: GetPageIndex(c),
		PageSize:  GetPageSize(c),
	}
}

// GetTraceID 获取追踪ID
func GetTraceID(c iris.Context) string {
	//return c.GetString(TraceIDKey)
	return c.URLParam(TraceIDKey)
}

// GetUserID 获取用户ID
func GetUserID(c iris.Context) string {
	//return c.GetString(UserIDKey)
	return c.URLParam(UserIDKey)
}

// SetUserID 设定用户ID
func SetUserID(c iris.Context, userID string) {
	c.Params().Set(UserIDKey, userID)
}

// ParseJSON 解析请求JSON
func ParseJSON(c iris.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		logger.Warnf(NewContext(c), err.Error())
		return errors.ErrInvalidRequestParameter
	}
	return nil
}

// ResPage 响应分页数据
func ResPage(c iris.Context, v interface{}, pr *schema.PaginationResult) {
	list := schema.HTTPList{
		List: v,
		Pagination: &schema.HTTPPagination{
			Current:  GetPageIndex(c),
			PageSize: GetPageSize(c),
		},
	}
	if pr != nil {
		list.Pagination.Total = pr.Total
	}

	ResSuccess(c, list)
}

// ResList 响应列表数据
func ResList(c iris.Context, v interface{}) {
	ResSuccess(c, schema.HTTPList{List: v})
}

// ResOK 响应OK
func ResOK(c iris.Context) {
	ResSuccess(c, schema.HTTPStatus{Status: "OK"})
}

// ResSuccess 响应成功
func ResSuccess(c iris.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

// ResJSON 响应JSON数据
func ResJSON(c iris.Context, status int, v interface{}) {
	buf, err := util.JSONMarshal(v)
	if err != nil {
		panic(err)
	}
	c.Params().Set(ResBodyKey, string(buf) )
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

// ResError 响应错误
func ResError(c iris.Context, err error, status ...int) {
	statusCode := 500
	errItem := schema.HTTPErrorItem{
		Code:    500,
		Message: "服务器发生错误",
	}

	if errCode, ok := errors.FromErrorCode(err); ok {
		errItem.Code = errCode.Code
		errItem.Message = errCode.Message
		statusCode = errCode.HTTPStatusCode
	}

	if len(status) > 0 {
		statusCode = status[0]
	}

	if statusCode == 500 && err != nil {
		span := logger.StartSpan(NewContext(c))
		span = span.WithField("stack", fmt.Sprintf("%+v", err))
		span.Errorf(err.Error())
	}

	ResJSON(c, statusCode, schema.HTTPError{Error: errItem})
}

