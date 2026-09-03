package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/service"
	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/response"
)

// MenuHandler 菜单管理处理器
type MenuHandler struct {
	menuService *service.MenuService
}

// NewMenuHandler 创建菜单管理处理器实例
func NewMenuHandler(menuService *service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

// List 获取菜单树
// @Summary 获取全部菜单树
// @Tags RBAC权限
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/menu/list [get]
func (h *MenuHandler) List(c *gin.Context) {
	tree, err := h.menuService.List()
	if err != nil {
		response.Fail(c, 200, errors.CodeServerError, errors.GetMsg(errors.CodeServerError))
		return
	}

	response.SuccessData(c, tree)
}

// Create 创建菜单
// @Summary 创建菜单
// @Tags RBAC权限
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateMenuReq true "菜单信息"
// @Success 200 {object} response.Response
// @Router /api/menu/create [post]
func (h *MenuHandler) Create(c *gin.Context) {
	var req dto.CreateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.menuService.Create(&req); err != nil {
		errCode, _ := strconv.Atoi(err.Error())
		if errCode == 0 {
			errCode = errors.CodeServerError
		}
		response.Fail(c, 200, errCode, errors.GetMsg(errCode))
		return
	}

	response.Success(c)
}

// Update 更新菜单
// @Summary 更新菜单
// @Tags RBAC权限
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateMenuReq true "菜单信息"
// @Success 200 {object} response.Response
// @Router /api/menu/update [put]
func (h *MenuHandler) Update(c *gin.Context) {
	var req dto.UpdateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.menuService.Update(&req); err != nil {
		errCode, _ := strconv.Atoi(err.Error())
		if errCode == 0 {
			errCode = errors.CodeServerError
		}
		response.Fail(c, 200, errCode, errors.GetMsg(errCode))
		return
	}

	response.Success(c)
}

// Delete 删除菜单
// @Summary 删除菜单
// @Tags RBAC权限
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response
// @Router /api/menu/{id} [delete]
func (h *MenuHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, errors.CodeBadRequest, errors.GetMsg(errors.CodeBadRequest))
		return
	}

	if err := h.menuService.Delete(id); err != nil {
		errCode, _ := strconv.Atoi(err.Error())
		if errCode == 0 {
			errCode = errors.CodeServerError
		}
		response.Fail(c, 200, errCode, errors.GetMsg(errCode))
		return
	}

	response.Success(c)
}
