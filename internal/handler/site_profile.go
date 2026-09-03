package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/service"
	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/response"
)

// SiteProfileHandler 站长个人资料处理器（区别于登录用户个人信息 ProfileHandler）
type SiteProfileHandler struct {
	siteProfileService *service.SiteProfileService
}

// NewSiteProfileHandler 创建站长个人资料处理器实例
func NewSiteProfileHandler(siteProfileService *service.SiteProfileService) *SiteProfileHandler {
	return &SiteProfileHandler{siteProfileService: siteProfileService}
}

// PublicProfile 获取站长个人资料（公开接口）
// @Summary 获取站长个人资料
// @Description 博客前台公开接口，用于「关于页」与「联系站长」展示，表为空时返回空默认值
// @Tags 博客公开接口
// @Produce json
// @Success 200 {object} response.Response{data=dto.SiteProfileResp}
// @Router /api/public/profile [get]
func (h *SiteProfileHandler) PublicProfile(c *gin.Context) {
	resp, err := h.siteProfileService.GetProfile()
	if err != nil {
		response.Fail(c, http.StatusOK, errors.CodeServerError, "获取站长资料失败")
		return
	}
	response.SuccessData(c, resp)
}

// GetProfile 获取站长个人资料（管理接口）
// @Summary 获取站长个人资料
// @Tags 站点管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.SiteProfileResp}
// @Router /api/site/profile [get]
func (h *SiteProfileHandler) GetProfile(c *gin.Context) {
	resp, err := h.siteProfileService.GetProfile()
	if err != nil {
		response.Fail(c, http.StatusOK, errors.CodeServerError, "获取站长资料失败")
		return
	}
	response.SuccessData(c, resp)
}

// SaveProfile 保存站长个人资料（整体覆盖保存）
// @Summary 保存站长个人资料
// @Tags 站点管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SaveSiteProfileReq true "站长个人资料"
// @Success 200 {object} response.Response
// @Router /api/site/profile [put]
func (h *SiteProfileHandler) SaveProfile(c *gin.Context) {
	var req dto.SaveSiteProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errors.CodeBadRequest, "参数错误")
		return
	}

	if err := h.siteProfileService.SaveProfile(&req); err != nil {
		response.Fail(c, http.StatusOK, errors.CodeServerError, "保存站长资料失败")
		return
	}

	response.Success(c)
}
