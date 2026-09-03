package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"wuzhispace.com/internal/repository"
	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/logger"
	"wuzhispace.com/pkg/response"
	"wuzhispace.com/pkg/storage"
)

// allowedMimeTypes 允许上传的文件 MIME 类型
var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// maxUploadSize 最大上传文件大小 10MB
const maxUploadSize = 10 * 1024 * 1024

// UploadHandler 文件上传处理器
type UploadHandler struct {
	siteConfigRepo *repository.SiteConfigRepository
}

// NewUploadHandler 创建上传处理器实例
func NewUploadHandler(siteConfigRepo *repository.SiteConfigRepository) *UploadHandler {
	return &UploadHandler{siteConfigRepo: siteConfigRepo}
}

// UploadOSS 上传文件到 OSS
// @Summary 上传文件到OSS
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "上传文件"
// @Success 200 {object} response.Response
// @Router /api/upload/oss [post]
func (h *UploadHandler) UploadOSS(c *gin.Context) {
	userID := c.GetInt64("userID")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logger.Warn("上传拒绝，未选择文件", "userID", userID, "error", err)
		response.Fail(c, http.StatusBadRequest, errors.CodeBadRequest, "请选择要上传的文件")
		return
	}
	defer file.Close()

	// 校验文件大小
	if header.Size > maxUploadSize {
		logger.Warn("上传拒绝，文件超过大小限制", "userID", userID, "filename", header.Filename, "size", header.Size)
		response.Fail(c, http.StatusOK, errors.ErrFileSizeExceed, errors.GetMsg(errors.ErrFileSizeExceed))
		return
	}

	// 校验文件类型
	contentType := header.Header.Get("Content-Type")
	if !allowedMimeTypes[contentType] {
		logger.Warn("上传拒绝，MIME类型非法", "userID", userID, "filename", header.Filename, "contentType", contentType)
		response.Fail(c, http.StatusOK, errors.ErrFileTypeNotSupport, errors.GetMsg(errors.ErrFileTypeNotSupport))
		return
	}

	// 从 DB 读取 OSS 配置
	siteCfg, err := h.siteConfigRepo.Get()
	if err != nil {
		logger.Error("上传读取站点配置失败", "userID", userID, "error", err)
		response.Fail(c, http.StatusOK, errors.ErrOSSNotConfigured, errors.GetMsg(errors.ErrOSSNotConfigured))
		return
	}

	ossCfg := storage.OssConfig{
		AccessKey:          siteCfg.OssAccessKey,
		SecretKey:          siteCfg.OssSecretKey,
		Bucket:             siteCfg.OssBucket,
		Endpoint:           siteCfg.OssEndpoint,
		Domain:             siteCfg.OssDomain,
		Provider:           siteCfg.OssProvider,
		Region:             siteCfg.OssRegion,
		InsecureSkipVerify: siteCfg.OssInsecure == 1,
	}

	// 校验 OSS 配置是否完整：仅纯 AWS S3（provider == "s3"）按 region 路由不强制 endpoint，
	// RustFS 等自建存储必须填写服务地址，与 aliyun/tencent/qiniu 同等校验
	if ossCfg.AccessKey == "" || ossCfg.SecretKey == "" || ossCfg.Bucket == "" || (ossCfg.Endpoint == "" && ossCfg.Provider != "s3") {
		logger.Error("上传拒绝，OSS配置不完整", "userID", userID, "provider", ossCfg.Provider)
		response.Fail(c, http.StatusOK, errors.ErrOSSNotConfigured, errors.GetMsg(errors.ErrOSSNotConfigured))
		return
	}

	// 根据配置的存储厂商创建对应的存储客户端（带实例缓存）
	ossClient, err := storage.NewStorage(siteCfg.OssProvider, ossCfg)
	if err != nil {
		logger.Error("上传创建OSS客户端失败", "userID", userID, "provider", siteCfg.OssProvider, "error", err)
		response.Fail(c, http.StatusOK, errors.ErrOSSUploadFailed, "创建OSS客户端失败: "+err.Error())
		return
	}

	// 生成对象键
	objectKey := storage.GenerateObjectKey(header.Filename)

	// 上传
	url, err := ossClient.Upload(c.Request.Context(), objectKey, file, contentType)
	if err != nil {
		logger.Error("OSS上传失败", "userID", userID, "provider", siteCfg.OssProvider, "filename", header.Filename, "key", objectKey, "error", err)
		response.Fail(c, http.StatusOK, errors.ErrOSSUploadFailed, "OSS上传失败: "+err.Error())
		return
	}

	logger.Info("OSS上传成功", "userID", userID, "provider", siteCfg.OssProvider, "filename", header.Filename, "key", objectKey, "url", url)
	response.SuccessData(c, gin.H{"url": url})
}

// isAllowedImageType 通过文件名判断是否为允许的图片类型
func isAllowedImageType(filename string) bool {
	ext := strings.ToLower(filename)
	return strings.HasSuffix(ext, ".jpg") ||
		strings.HasSuffix(ext, ".jpeg") ||
		strings.HasSuffix(ext, ".png") ||
		strings.HasSuffix(ext, ".gif") ||
		strings.HasSuffix(ext, ".webp")
}
