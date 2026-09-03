package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"

	"wuzhispace.com/pkg/logger"
)

// TencentCOS 腾讯云 COS 实现
type TencentCOS struct {
	client *cos.Client
	bucket string
	domain string
}

// NewTencentCOS 创建腾讯云 COS 客户端实例
// Endpoint 格式: bucket-appid.cos.region.myqcloud.com
func NewTencentCOS(cfg OssConfig) (*TencentCOS, error) {
	endpoint := cfg.Endpoint
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid tencent endpoint: %w", err)
	}

	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.AccessKey,
			SecretKey: cfg.SecretKey,
			// 开启跳过证书校验时，将不安全 Transport 作为签名层的底层 Transport
			Transport: insecureTransport(cfg.InsecureSkipVerify),
		},
	})

	return &TencentCOS{
		client: client,
		bucket: cfg.Bucket,
		domain: cfg.Domain,
	}, nil
}

// insecureTransport 根据开关返回底层 Transport：开启时跳过 TLS 证书校验，否则返回 nil 沿用 SDK 默认行为
func insecureTransport(skipVerify bool) http.RoundTripper {
	if !skipVerify {
		return nil
	}
	return newInsecureHTTPClient().Transport
}

// Upload 上传文件到 COS，返回访问 URL
func (t *TencentCOS) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	}
	_, err := t.client.Object.Put(ctx, key, reader, opt)
	if err != nil {
		return "", fmt.Errorf("COS上传失败: %w", err)
	}
	return t.buildURL(key), nil
}

// Delete 从 COS 删除文件
func (t *TencentCOS) Delete(ctx context.Context, key string) error {
	_, err := t.client.Object.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("COS删除失败: %w", err)
	}
	return nil
}

// TestConnectivity 测试 COS 连通性：上传 4 字节测试文件后立即删除
func (t *TencentCOS) TestConnectivity(ctx context.Context) error {
	testKey := "blog/_test_connectivity.txt"
	_, err := t.client.Object.Put(ctx, testKey, strings.NewReader("test"), &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{ContentType: "text/plain"},
	})
	if err != nil {
		logger.Error("OSS连通测试失败", "provider", "tencent", "error", err)
		return fmt.Errorf("COS连通测试上传失败: %w", err)
	}
	if _, err = t.client.Object.Delete(ctx, testKey); err != nil {
		logger.Error("OSS连通测试失败", "provider", "tencent", "error", err)
		return fmt.Errorf("COS连通测试删除失败: %w", err)
	}
	return nil
}

// buildURL 构建文件的访问 URL
func (t *TencentCOS) buildURL(key string) string {
	if t.domain != "" {
		domain := t.domain
		if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
			domain = "https://" + domain
		}
		return strings.TrimRight(domain, "/") + "/" + key
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(t.client.BaseURL.BucketURL.String(), "/"), key)
}
