package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/service"
	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/response"
)

// ProfileHandler 个人信息处理器
type ProfileHandler struct {
	profileService *service.ProfileService
}

// NewProfileHandler 创建个人信息处理器实例
func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

// GetInfo 获取当前用户个人信息
// @Summary 获取个人信息
// @Description 返回当前登录用户的昵称、头像、个人简介等信息
// @Tags 个人信息
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.ProfileDTO}
// @Router /api/profile/info [get]
func (h *ProfileHandler) GetInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, errors.CodeUnauthorized, errors.GetMsg(errors.CodeUnauthorized))
		return
	}

	profile, err := h.profileService.GetProfile(userID.(int64))
	if err != nil {
		errCode, _ := strconv.Atoi(err.Error())
		if errCode == 0 {
			errCode = errors.CodeServerError
		}
		response.Fail(c, http.StatusOK, errCode, errors.GetMsg(errors.CodeServerError))
		return
	}

	response.SuccessData(c, profile)
}

// Update 更新个人信息（昵称、头像、简介）
// @Summary 更新个人信息
// @Description 更新当前登录用户的昵称、头像、个人简介，更新后清除用户缓存
// @Tags 个人信息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileReq true "个人信息"
// @Success 200 {object} response.Response{data=dto.ProfileDTO}
// @Router /api/profile/update [put]
func (h *ProfileHandler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, errors.CodeUnauthorized, errors.GetMsg(errors.CodeUnauthorized))
		return
	}

	var req dto.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errors.CodeBadRequest, "参数错误")
		return
	}

	profile, err := h.profileService.UpdateProfile(userID.(int64), req)
	if err != nil {
		errCode, _ := strconv.Atoi(err.Error())
		if errCode == 0 {
			errCode = errors.CodeServerError
		}
		response.Fail(c, http.StatusOK, errCode, errors.GetMsg(errors.CodeServerError))
		return
	}

	response.SuccessData(c, profile)
}
