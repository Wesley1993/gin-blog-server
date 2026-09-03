package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/service"
	pkgerrors "wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/logger"
	"wuzhispace.com/pkg/response"
)

// SiteLinkHandler 常用网站处理器
type SiteLinkHandler struct {
	siteLinkService *service.SiteLinkService
}

// NewSiteLinkHandler 创建常用网站处理器实例
func NewSiteLinkHandler(siteLinkService *service.SiteLinkService) *SiteLinkHandler {
	return &SiteLinkHandler{siteLinkService: siteLinkService}
}

// PublicLinks 获取启用中的常用网站列表（公开接口）
// @Summary 获取常用网站列表
// @Description 博客前台公开接口，仅返回启用中的记录，按 sort 升序、id 升序
// @Tags 博客公开接口
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.SiteLinkItem}
// @Router /api/public/links [get]
func (h *SiteLinkHandler) PublicLinks(c *gin.Context) {
	items, err := h.siteLinkService.PublicList()
	if err != nil {
		logger.Error("获取常用网站列表失败", "error", err)
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "获取常用网站列表失败")
		return
	}
	response.SuccessData(c, items)
}

// List 获取全部常用网站（含停用，管理接口）
// @Summary 获取全部常用网站
// @Tags 站点管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.SiteLink}
// @Router /api/site/links [get]
func (h *SiteLinkHandler) List(c *gin.Context) {
	links, err := h.siteLinkService.List()
	if err != nil {
		logger.Error("获取常用网站列表失败", "error", err)
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "获取常用网站列表失败")
		return
	}
	response.SuccessData(c, links)
}

// Create 新增常用网站
// @Summary 新增常用网站
// @Tags 站点管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateSiteLinkReq true "常用网站信息"
// @Success 200 {object} response.Response
// @Router /api/site/links [post]
func (h *SiteLinkHandler) Create(c *gin.Context) {
	var req dto.CreateSiteLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	if err := h.siteLinkService.Create(req); err != nil {
		logger.Error("新增常用网站失败", "name", req.Name, "error", err)
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "新增常用网站失败")
		return
	}

	response.Success(c)
}

// Update 更新常用网站
// @Summary 更新常用网站
// @Tags 站点管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "常用网站ID"
// @Param request body dto.UpdateSiteLinkReq true "常用网站信息"
// @Success 200 {object} response.Response
// @Router /api/site/links/{id} [put]
func (h *SiteLinkHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	var req dto.UpdateSiteLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	if err := h.siteLinkService.Update(id, req); err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			response.Fail(c, http.StatusOK, pkgerrors.ErrLinkNotFound, pkgerrors.GetMsg(pkgerrors.ErrLinkNotFound))
			return
		}
		logger.Error("更新常用网站失败", "id", id, "error", err)
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "更新常用网站失败")
		return
	}

	response.Success(c)
}

// Delete 删除常用网站
// @Summary 删除常用网站
// @Tags 站点管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "常用网站ID"
// @Success 200 {object} response.Response
// @Router /api/site/links/{id} [delete]
func (h *SiteLinkHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	if err := h.siteLinkService.Delete(id); err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			response.Fail(c, http.StatusOK, pkgerrors.ErrLinkNotFound, pkgerrors.GetMsg(pkgerrors.ErrLinkNotFound))
			return
		}
		logger.Error("删除常用网站失败", "id", id, "error", err)
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "删除常用网站失败")
		return
	}

	response.Success(c)
}
