package ctl

import (
	"github.com/LyricTian/captcha"
	"github.com/wanhello/iris-admin/internal/app/bll"
	"github.com/wanhello/iris-admin/internal/app/config"
	"github.com/wanhello/iris-admin/internal/app/errors"
	"github.com/wanhello/iris-admin/internal/app/schema"
	"github.com/wanhello/iris-admin/internal/app/irisplus"
	"github.com/wanhello/iris-admin/pkg/logger"

	"github.com/kataras/iris"

)

// NewLogin 创建登录管理控制器
func NewLogin(bLogin bll.ILogin) *Login {
	return &Login{
		LoginBll: bLogin,
	}
}

// Login 登录管理
// @Name Login
// @Description 登录管理接口
type Login struct {
	LoginBll bll.ILogin
}

// GetCaptcha 获取验证码信息
// @Summary 获取验证码信息
// @Success 200 schema.LoginCaptcha
// @Router GET /api/v1/pub/login/captchaid
func (a *Login) GetCaptcha(c iris.Context) {
	item, err := a.LoginBll.GetCaptcha(irisplus.NewContext(c), config.GetGlobalConfig().Captcha.Length)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResSuccess(c, item)
}

// ResCaptcha 响应图形验证码
// @Summary 响应图形验证码
// @Param id query string true "验证码ID"
// @Param reload query string false "重新加载"
// @Success 200 file "图形验证码"
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/pub/login/captcha
func (a *Login) ResCaptcha(c iris.Context) {
	captchaID := c.URLParam("id")
	if captchaID == "" {
		irisplus.ResError(c, errors.ErrInvalidRequestParameter)
		return
	}

	if c.URLParam("reload") != "" {
		if !captcha.Reload(captchaID) {
			irisplus.ResError(c, errors.ErrInvalidRequestParameter)
			return
		}
	}

	cfg := config.GetGlobalConfig().Captcha
	err := a.LoginBll.ResCaptcha(irisplus.NewContext(c), c.ResponseWriter(), captchaID, cfg.Width, cfg.Height)
	if err != nil {
		irisplus.ResError(c, err)
	}
}

// Login 用户登录
// @Summary 用户登录
// @Param body body schema.LoginParam true
// @Success 200 schema.LoginTokenInfo
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/pub/login
func (a *Login) Login(c iris.Context) {
	var item schema.LoginParam
	if err := irisplus.ParseJSON(c, &item); err != nil {
		irisplus.ResError(c, err)
		return
	}

	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		irisplus.ResError(c, errors.ErrLoginInvalidVerifyCode)
		return
	}

	user, err := a.LoginBll.Verify(irisplus.NewContext(c), item.UserName, item.Password)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}

	userID := user.RecordID
	// 将用户ID放入上下文
	irisplus.SetUserID(c, userID)

	tokenInfo, err := a.LoginBll.GenerateToken(irisplus.NewContext(c), userID)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}

	logger.StartSpan(irisplus.NewContext(c), logger.SetSpanTitle("用户登录"), logger.SetSpanFuncName("Login")).Infof("登入系统")
	irisplus.ResSuccess(c, tokenInfo)
}

// Logout 用户登出
// @Summary 用户登出
// @Success 200 schema.HTTPStatus "{status:OK}"
// @Router POST /api/v1/pub/login/exit
func (a *Login) Logout(c iris.Context) {
	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := irisplus.GetUserID(c)
	if userID != "" {
		ctx := irisplus.NewContext(c)
		err := a.LoginBll.DestroyToken(ctx, irisplus.GetToken(c))
		if err != nil {
			logger.Errorf(ctx, err.Error())
		}
		logger.StartSpan(irisplus.NewContext(c), logger.SetSpanTitle("用户登出"), logger.SetSpanFuncName("Logout")).Infof("登出系统")
	}
	irisplus.ResOK(c)
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 schema.LoginTokenInfo "{access_token:访问令牌,token_type:令牌类型,expires_in:过期时长(单位秒)}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/pub/refresh_token
func (a *Login) RefreshToken(c iris.Context) {
	tokenInfo, err := a.LoginBll.GenerateToken(irisplus.NewContext(c), irisplus.GetUserID(c))
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResSuccess(c, tokenInfo)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 schema.UserLoginInfo
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/pub/current/user
func (a *Login) GetUserInfo(c iris.Context) {
	info, err := a.LoginBll.GetLoginInfo(irisplus.NewContext(c), irisplus.GetUserID(c))
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResSuccess(c, info)
}

// QueryUserMenuTree 查询当前用户菜单树
// @Summary 查询当前用户菜单树
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 schema.Menu "查询结果：{list:菜单树}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/pub/current/menutree
func (a *Login) QueryUserMenuTree(c iris.Context) {
	menus, err := a.LoginBll.QueryUserMenuTree(irisplus.NewContext(c), irisplus.GetUserID(c))
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResList(c, menus)
}

// UpdatePassword 更新个人密码
// @Summary 更新个人密码
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.UpdatePasswordParam true
// @Success 200 schema.HTTPStatus "{status:OK}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/pub/current/password
func (a *Login) UpdatePassword(c iris.Context) {
	var item schema.UpdatePasswordParam
	if err := irisplus.ParseJSON(c, &item); err != nil {
		irisplus.ResError(c, err)
		return
	}

	err := a.LoginBll.UpdatePassword(irisplus.NewContext(c), irisplus.GetUserID(c), item)
	if err != nil {
		irisplus.ResError(c, err)
		return
	}
	irisplus.ResOK(c)
}
