package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wuzhispace.com/internal/model"
	"wuzhispace.com/internal/service"
	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/logger"
	"wuzhispace.com/pkg/response"
	"wuzhispace.com/pkg/storage"
)

// SiteHandler 站点配置处理器
type SiteHandler struct {
	siteService *service.SiteConfigService
}

// NewSiteHandler 创建站点配置处理器实例
func NewSiteHandler(siteService *service.SiteConfigService) *SiteHandler {
	return &SiteHandler{siteService: siteService}
}

// SaveConfigReq 保存站点配置请求体
type SaveConfigReq struct {
	SiteName     string `json:"site_name"`
	SiteDesc     string `json:"site_desc"`
	Copyright    string `json:"copyright"`
	FoundedAt    string `json:"founded_at"`
	Icp          string `json:"icp"`
	PoliceIcp    string `json:"police_icp"`
	OssAccessKey string `json:"oss_access_key"`
	OssSecretKey string `json:"oss_secret_key"`
	OssBucket    string `json:"oss_bucket"`
	OssEndpoint  string `json:"oss_endpoint"`
	OssDomain    string `json:"oss_domain"`
	OssProvider  string `json:"oss_provider"`
	OssRegion    string `json:"oss_region"`
	OssInsecure  int    `json:"oss_insecure"`
}

// GetConfig 获取站点配置
// @Summary 获取站点配置
// @Tags 站点管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/site/config [get]
func (h *SiteHandler) GetConfig(c *gin.Context) {
	cfg, err := h.siteService.GetConfig()
	if err != nil {
		response.Fail(c, http.StatusOK, errors.CodeServerError, "获取站点配置失败")
		return
	}
	response.SuccessData(c, cfg)
}

// PublicConfig 获取站点公开配置（仅返回展示字段，不返回 OSS 配置）
// @Summary 获取站点公开配置
// @Description 博客前台公开接口，仅返回 site_name、site_desc、copyright、founded_at、icp、police_icp
// @Tags 博客公开接口
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/public/site/config [get]
func (h *SiteHandler) PublicConfig(c *gin.Context) {
	cfg, err := h.siteService.GetConfig()
	if err != nil {
		response.Fail(c, http.StatusOK, errors.CodeServerError, "获取站点配置失败")
		return
	}

	public := gin.H{
		"site_name":  "",
		"site_desc":  "",
		"copyright":  "",
		"founded_at": "",
		"icp":        "",
		"police_icp": "",
	}
	if cfg != nil {
		public["site_name"] = cfg.SiteName
		public["site_desc"] = cfg.SiteDesc
		public["copyright"] = cfg.Copyright
		public["founded_at"] = cfg.FoundedAt
		public["icp"] = cfg.Icp
		public["police_icp"] = cfg.PoliceIcp
	}

	response.SuccessData(c, public)
}

// SaveConfig 保存站点配置
// @Summary 保存站点配置
// @Tags 站点管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body handler.SaveConfigReq true "站点配置"
// @Success 200 {object} response.Response
// @Router /api/site/config [put]
func (h *SiteHandler) SaveConfig(c *gin.Context) {
	var req SaveConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errors.CodeBadRequest, "参数错误")
		return
	}

	// 先尝试获取已有配置（用于获取 ID）
	// 直接读 DB 而非缓存：Redis 缓存可能是无 provider 字段的旧结构，
	// 会导致 InvalidateInstance 用错误的 key 清理不到旧实例，引发厂商路由残留
	existing, _ := h.siteService.GetOssConfigSource()

	cfg := &model.SiteConfig{
		SiteName:     req.SiteName,
		SiteDesc:     req.SiteDesc,
		Copyright:    req.Copyright,
		FoundedAt:    req.FoundedAt,
		Icp:          req.Icp,
		PoliceIcp:    req.PoliceIcp,
		OssAccessKey: req.OssAccessKey,
		OssSecretKey: req.OssSecretKey,
		OssBucket:    req.OssBucket,
		OssEndpoint:  req.OssEndpoint,
		OssDomain:    req.OssDomain,
		OssProvider:  req.OssProvider,
		OssRegion:    req.OssRegion,
		OssInsecure:  req.OssInsecure,
	}

	if existing != nil {
		cfg.ID = existing.ID
	}

	if err := h.siteService.SaveConfig(cfg); err != nil {
		logger.Error("保存站点配置失败", "userID", c.GetInt64("userID"), "error", err)
		response.Fail(c, http.StatusOK, errors.CodeServerError, "保存站点配置失败")
		return
	}

	// 配置变更后清理旧的存储客户端实例缓存（新旧配置均清理，避免残留）
	if existing != nil {
		storage.InvalidateInstance(existing.OssProvider, existing.OssEndpoint, existing.OssBucket, existing.OssRegion, existing.OssInsecure != 0)
	}
	storage.InvalidateInstance(req.OssProvider, req.OssEndpoint, req.OssBucket, req.OssRegion, req.OssInsecure != 0)

	logger.Info("保存站点配置成功", "userID", c.GetInt64("userID"), "ossProvider", req.OssProvider, "ossRegion", req.OssRegion, "ossInsecure", req.OssInsecure)
	response.Success(c)
}

// Stats 获取站点运行统计信息
// @Summary 获取站点运行统计信息
// @Description 返回建站日期、运行天数、文章总数、分类数、用户数，运行天数由后端根据 founded_at 计算
// @Tags 站点管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.SiteStatsResp}
// @Router /api/site/stats [get]
func (h *SiteHandler) Stats(c *gin.Context) {
	stats, err := h.siteService.GetStats()
	if err != nil {
		response.Fail(c, http.StatusOK, errors.CodeServerError, "获取站点统计信息失败")
		return
	}
	response.SuccessData(c, stats)
}

// TestOss OSS 连通测试
// @Summary OSS连通测试
// @Tags 站点管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/site/testOss [post]
func (h *SiteHandler) TestOss(c *gin.Context) {
	ossCfg, err := h.siteService.GetOssConfig()
	if err != nil {
		response.Fail(c, http.StatusOK, errors.ErrOSSNotConfigured, errors.GetMsg(errors.ErrOSSNotConfigured))
		return
	}

	// 仅纯 AWS S3（provider == "s3"）不强制 endpoint；RustFS 等自建存储必填服务地址，与其他厂商同等校验
	if ossCfg.AccessKey == "" || ossCfg.SecretKey == "" || ossCfg.Bucket == "" || (ossCfg.Endpoint == "" && ossCfg.Provider != "s3") {
		logger.Warn("OSS连通测试拒绝，配置不完整", "userID", c.GetInt64("userID"), "provider", ossCfg.Provider)
		response.Fail(c, http.StatusOK, errors.ErrOSSNotConfigured, errors.GetMsg(errors.ErrOSSNotConfigured))
		return
	}

	ossClient, err := storage.NewStorage(ossCfg.Provider, *ossCfg)
	if err != nil {
		logger.Error("OSS连通测试创建客户端失败", "provider", ossCfg.Provider, "error", err)
		response.Fail(c, http.StatusOK, errors.ErrOSSUploadFailed, "创建OSS客户端失败: "+err.Error())
		return
	}

	if err := ossClient.TestConnectivity(c.Request.Context()); err != nil {
		response.Fail(c, http.StatusOK, errors.ErrOSSUploadFailed, "OSS连通测试失败: "+err.Error())
		return
	}

	logger.Info("OSS连通测试成功", "userID", c.GetInt64("userID"), "provider", ossCfg.Provider)
	response.SuccessData(c, gin.H{"message": "OSS连通测试成功"})
}
