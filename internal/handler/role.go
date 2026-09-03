package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/service"
	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/response"
)

// RoleHandler 角色管理处理器
type RoleHandler struct {
	roleService *service.RoleService
}

// NewRoleHandler 创建角色管理处理器实例
func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// List 角色列表
// @Summary 角色列表
// @Tags RBAC权限
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/role/list [get]
func (h *RoleHandler) List(c *gin.Context) {
	list, err := h.roleService.List()
	if err != nil {
		response.Fail(c, 200, errors.CodeServerError, errors.GetMsg(errors.CodeServerError))
		return
	}

	response.SuccessData(c, list)
}

// Create 创建角色
// @Summary 创建角色
// @Tags RBAC权限
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateRoleReq true "角色信息"
// @Success 200 {object} response.Response
// @Router /api/role/create [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.roleService.Create(req); err != nil {
		response.Fail(c, 200, errors.CodeServerError, errors.GetMsg(errors.CodeServerError))
		return
	}

	response.Success(c)
}

// Update 更新角色
// @Summary 更新角色
// @Tags RBAC权限
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateRoleReq true "角色信息"
// @Success 200 {object} response.Response
// @Router /api/role/update [put]
func (h *RoleHandler) Update(c *gin.Context) {
	var req dto.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.roleService.Update(req); err != nil {
		errCode, _ := strconv.Atoi(err.Error())
		if errCode == 0 {
			errCode = errors.CodeServerError
		}
		response.Fail(c, 200, errCode, errors.GetMsg(errCode))
		return
	}

	response.Success(c)
}

// Delete 删除角色
// @Summary 删除角色
// @Tags RBAC权限
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/role/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.roleService.Delete(id); err != nil {
		errCode, _ := strconv.Atoi(err.Error())
		if errCode == 0 {
			errCode = errors.CodeServerError
		}
		response.Fail(c, 200, errCode, errors.GetMsg(errCode))
		return
	}

	response.Success(c)
}
