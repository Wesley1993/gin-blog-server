package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/service"
	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login 用户登录
// @Summary 用户登录
// @Description 账号密码登录，返回 JWT token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.LoginReq true "登录参数"
// @Success 200 {object} response.Response{data=dto.LoginResp}
// @Failure 401 {object} response.Response
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		errCode, _ := strconv.Atoi(err.Error())
		if errCode == 0 {
			errCode = errors.CodeServerError
		}
		response.Fail(c, 200, errCode, errors.GetMsg(errCode))
		return
	}

	response.SuccessData(c, resp)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 清除 Redis 中的 token
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Fail(c, 401, errors.CodeUnauthorized, errors.GetMsg(errors.CodeUnauthorized))
		return
	}

	_ = h.authService.Logout(userID.(int64))
	response.Success(c)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取用户信息
// @Description 返回当前用户信息、菜单树、权限标识
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserInfoResp}
// @Router /api/auth/userinfo [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Fail(c, 401, errors.CodeUnauthorized, errors.GetMsg(errors.CodeUnauthorized))
		return
	}

	resp, err := h.authService.GetUserInfo(userID.(int64))
	if err != nil {
		response.Fail(c, 200, errors.CodeServerError, errors.GetMsg(errors.CodeServerError))
		return
	}

	response.SuccessData(c, resp)
}
