package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/service"
	pkgerrors "wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/response"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler 创建分类处理器实例
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// Tree 返回分类树
// @Summary 获取分类树
// @Tags 分类管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/category/tree [get]
func (h *CategoryHandler) Tree(c *gin.Context) {
	tree, err := h.categoryService.Tree()
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "获取分类树失败")
		return
	}
	response.SuccessData(c, tree)
}

// Create 新增分类
// @Summary 新增分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryReq true "分类信息"
// @Success 200 {object} response.Response
// @Router /api/category/create [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	if err := h.categoryService.Create(req); err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "创建分类失败")
		return
	}

	response.Success(c)
}

// Update 更新分类
// @Summary 更新分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateCategoryReq true "分类信息"
// @Success 200 {object} response.Response
// @Router /api/category/update [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	var req dto.UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	if err := h.categoryService.Update(req); err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			response.Fail(c, http.StatusOK, pkgerrors.ErrCategoryNotFound, pkgerrors.GetMsg(pkgerrors.ErrCategoryNotFound))
			return
		}
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "更新分类失败")
		return
	}

	response.Success(c)
}

// Delete 删除分类
// @Summary 删除分类
// @Tags 分类管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response "分类被引用禁止删除"
// @Router /api/category/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	if err := h.categoryService.Delete(id); err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			response.Fail(c, http.StatusOK, pkgerrors.ErrCategoryNotFound, pkgerrors.GetMsg(pkgerrors.ErrCategoryNotFound))
			return
		}
		if errors.Is(err, service.ErrCategoryInUse) {
			response.Fail(c, http.StatusOK, pkgerrors.ErrCategoryInUse, pkgerrors.GetMsg(pkgerrors.ErrCategoryInUse))
			return
		}
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "删除分类失败")
		return
	}

	response.Success(c)
}
