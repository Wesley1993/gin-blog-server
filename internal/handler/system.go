package handler

import (
	"github.com/gin-gonic/gin"

	"wuzhispace.com/internal/service"
	"wuzhispace.com/pkg/response"
)

// SystemHandler 系统运行信息处理器
type SystemHandler struct {
	systemService *service.SystemService
}

// NewSystemServiceHandler 创建系统处理器实例
func NewSystemServiceHandler(svc *service.SystemService) *SystemHandler {
	return &SystemHandler{systemService: svc}
}

// Overview 获取系统运行概览
// @Summary 获取系统运行概览
// @Description 返回应用版本、Go 运行时、内存、运行时长及 DB/Redis/ES 健康状态
// @Tags 系统管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=service.SystemOverview}
// @Router /api/system/overview [get]
func (h *SystemHandler) Overview(c *gin.Context) {
	overview := h.systemService.GetOverview()
	response.SuccessData(c, overview)
}
