package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	qiniuclient "github.com/qiniu/go-sdk/v7/client"
	qiniustorage "github.com/qiniu/go-sdk/v7/storage"

	"wuzhispace.com/pkg/logger"
)

// QiniuOSS 七牛云 Kodo 实现
type QiniuOSS struct {
	bucket    string
	domain    string
	mac       *qbox.Mac
	putPolicy qiniustorage.PutPolicy
	// client 自定义 HTTP 客户端（跳过 TLS 证书校验时非 nil，否则传 nil 沿用 SDK 默认客户端）
	httpClient *qiniuclient.Client
}

// NewQiniuOSS 创建七牛云存储客户端实例
func NewQiniuOSS(cfg OssConfig) (*QiniuOSS, error) {
	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)

	// 七牛 SDK 的 storage.Config 无 Transport 注入点，通过 NewFormUploaderEx/NewBucketManagerEx
	// 传入自定义 client.Client（包装跳过证书校验的 http.Transport）实现，不修改全局 DefaultClient
	var httpClient *qiniuclient.Client
	if cfg.InsecureSkipVerify {
		httpClient = &qiniuclient.Client{Client: newInsecureHTTPClient()}
	}

	return &QiniuOSS{
		bucket: cfg.Bucket,
		domain: cfg.Domain,
		mac:    mac,
		putPolicy: qiniustorage.PutPolicy{
			Scope: cfg.Bucket,
		},
		httpClient: httpClient,
	}, nil
}

// Upload 上传文件到七牛云，返回访问 URL
func (q *QiniuOSS) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("七牛云读取上传数据失败: %w", err)
	}

	upToken := q.putPolicy.UploadToken(q.mac)
	cfg := qiniustorage.Config{UseHTTPS: true}
	formUploader := qiniustorage.NewFormUploaderEx(&cfg, q.httpClient)

	var ret qiniustorage.PutRet
	putExtra := qiniustorage.PutExtra{
		MimeType: contentType,
	}

	err = formUploader.Put(ctx, &ret, upToken, key, bytes.NewReader(data), int64(len(data)), &putExtra)
	if err != nil {
		return "", fmt.Errorf("七牛云上传失败: %w", err)
	}

	return fmt.Sprintf("%s/%s", strings.TrimRight(q.accessDomain(), "/"), key), nil
}

// Delete 从七牛云删除文件
func (q *QiniuOSS) Delete(_ context.Context, key string) error {
	cfg := qiniustorage.Config{UseHTTPS: true}
	bucketManager := qiniustorage.NewBucketManagerEx(q.mac, &cfg, q.httpClient)
	if err := bucketManager.Delete(q.bucket, key); err != nil {
		return fmt.Errorf("七牛云删除失败: %w", err)
	}
	return nil
}

// TestConnectivity 测试七牛云连通性：上传 4 字节测试文件后立即删除
func (q *QiniuOSS) TestConnectivity(ctx context.Context) error {
	testKey := "blog/_test_connectivity.txt"
	if _, err := q.Upload(ctx, testKey, strings.NewReader("test"), "text/plain"); err != nil {
		logger.Error("OSS连通测试失败", "provider", "qiniu", "error", err)
		return fmt.Errorf("七牛云连通测试上传失败: %w", err)
	}
	if err := q.Delete(ctx, testKey); err != nil {
		logger.Error("OSS连通测试失败", "provider", "qiniu", "error", err)
		return fmt.Errorf("七牛云连通测试删除失败: %w", err)
	}
	return nil
}

// accessDomain 获取访问域名（七牛云必须配置测试域名或自定义域名）
func (q *QiniuOSS) accessDomain() string {
	if q.domain != "" {
		domain := q.domain
		if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
			domain = "https://" + domain
		}
		return domain
	}
	return ""
}
