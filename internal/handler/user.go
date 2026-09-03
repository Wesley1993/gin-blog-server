package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/service"
	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/response"
)

// UserHandler 用户管理处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户管理处理器实例
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Page 用户分页列表
// @Summary 用户分页列表
// @Tags 用户管理
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(10)
// @Success 200 {object} response.Response
// @Router /api/user/page [get]
func (h *UserHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	list, total, err := h.userService.Page(page, pageSize)
	if err != nil {
		response.Fail(c, 200, errors.CodeServerError, errors.GetMsg(errors.CodeServerError))
		return
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// Create 创建用户
// @Summary 创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserReq true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/user/create [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.userService.Create(req); err != nil {
		errCode, _ := strconv.Atoi(err.Error())
		if errCode == 0 {
			errCode = errors.CodeServerError
		}
		response.Fail(c, 200, errCode, errors.GetMsg(errCode))
		return
	}

	response.Success(c)
}

// Update 更新用户
// @Summary 编辑用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserReq true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/user/update [put]
func (h *UserHandler) Update(c *gin.Context) {
	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.userService.Update(req); err != nil {
		response.Fail(c, 200, errors.CodeServerError, errors.GetMsg(errors.CodeServerError))
		return
	}

	response.Success(c)
}

// UpdateStatus 启用/禁用用户（独立接口，对应 user:status 权限）
// @Summary 启停用用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body dto.UpdateUserStatusReq true "目标状态"
// @Success 200 {object} response.Response
// @Router /api/user/status/{id} [put]
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	var req dto.UpdateUserStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.userService.UpdateStatus(userID, req.Status); err != nil {
		response.Fail(c, 200, errors.CodeServerError, errors.GetMsg(errors.CodeServerError))
		return
	}

	response.Success(c)
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body dto.ResetPwdReq true "新密码"
// @Success 200 {object} response.Response
// @Router /api/user/resetPwd/{id} [put]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	var req dto.ResetPwdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.userService.ResetPassword(userID, req.Password); err != nil {
		response.Fail(c, 200, errors.CodeServerError, errors.GetMsg(errors.CodeServerError))
		return
	}

	response.Success(c)
}

// Delete 删除用户
// @Summary 删除用户
// @Tags 用户管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/user/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.userService.Delete(userID); err != nil {
		errCode, _ := strconv.Atoi(err.Error())
		if errCode == 0 {
			errCode = errors.CodeServerError
		}
		response.Fail(c, 200, errCode, errors.GetMsg(errCode))
		return
	}

	response.Success(c)
}
