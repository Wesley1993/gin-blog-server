package storage

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"

	"wuzhispace.com/pkg/logger"
)

// ObjectStorage 对象存储接口
type ObjectStorage interface {
	// Upload 上传文件，返回访问 URL
	Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error)
	// Delete 删除文件
	Delete(ctx context.Context, key string) error
	// TestConnectivity 测试连通性（上传测试文件后删除）
	TestConnectivity(ctx context.Context) error
}

// OssConfig OSS 配置
type OssConfig struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
	Domain    string
	Provider  string
	Region    string
	// InsecureSkipVerify 跳过 HTTPS 证书校验（自签名证书/私有端点场景）
	InsecureSkipVerify bool
}

// newInsecureHTTPClient 创建跳过 TLS 证书校验的 HTTP 客户端（仅供各厂商实现内部使用）
func newInsecureHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // 用户显式开启的自签名证书/私有端点场景
		},
	}
}

// AliyunOSS 阿里云 OSS 实现
type AliyunOSS struct {
	client *oss.Client
	bucket *oss.Bucket
	config OssConfig
}

// NewAliyunOSS 创建阿里云 OSS 客户端实例
func NewAliyunOSS(cfg OssConfig) (*AliyunOSS, error) {
	var clientOpts []oss.ClientOption
	if cfg.InsecureSkipVerify {
		clientOpts = append(clientOpts, oss.HTTPClient(newInsecureHTTPClient()))
	}
	client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取OSS Bucket失败: %w", err)
	}

	return &AliyunOSS{
		client: client,
		bucket: bucket,
		config: cfg,
	}, nil
}

// Upload 上传文件到 OSS，返回访问 URL
func (o *AliyunOSS) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	options := []oss.Option{
		oss.ContentType(contentType),
	}

	err := o.bucket.PutObject(key, reader, options...)
	if err != nil {
		return "", fmt.Errorf("OSS上传失败: %w", err)
	}

	url := o.buildURL(key)
	return url, nil
}

// Delete 从 OSS 删除文件
func (o *AliyunOSS) Delete(ctx context.Context, key string) error {
	err := o.bucket.DeleteObject(key)
	if err != nil {
		return fmt.Errorf("OSS删除失败: %w", err)
	}
	return nil
}

// TestConnectivity 测试 OSS 连通性：上传 4 字节测试文件后立即删除
func (o *AliyunOSS) TestConnectivity(ctx context.Context) error {
	testKey := "blog/_test_connectivity.txt"
	testContent := strings.NewReader("test")

	err := o.bucket.PutObject(testKey, testContent, oss.ContentType("text/plain"))
	if err != nil {
		logger.Error("OSS连通测试失败", "provider", "aliyun", "error", err)
		return fmt.Errorf("OSS连通测试上传失败: %w", err)
	}

	err = o.bucket.DeleteObject(testKey)
	if err != nil {
		logger.Error("OSS连通测试失败", "provider", "aliyun", "error", err)
		return fmt.Errorf("OSS连通测试删除失败: %w", err)
	}

	return nil
}

// buildURL 构建文件的访问 URL
func (o *AliyunOSS) buildURL(key string) string {
	if o.config.Domain != "" {
		domain := o.config.Domain
		if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
			domain = "https://" + domain
		}
		return strings.TrimRight(domain, "/") + "/" + key
	}
	// 使用默认 endpoint 构建 URL
	endpoint := o.config.Endpoint
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}
	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(endpoint, "/"), o.config.Bucket, key)
}

// GenerateObjectKey 生成对象键：blog/{year}/{month}/{uuid}.{ext}
func GenerateObjectKey(filename string) string {
	now := time.Now()
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg"
	}
	uid := uuid.New().String()
	return fmt.Sprintf("blog/%d/%02d/%s%s", now.Year(), now.Month(), uid, ext)
}
