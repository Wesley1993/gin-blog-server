package storage

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"wuzhispace.com/pkg/logger"
)

// AmazonS3 Amazon S3 实现（同时服务于 RustFS 等 S3 兼容的自建存储）
type AmazonS3 struct {
	client   *s3.Client
	bucket   string
	region   string
	domain   string
	endpoint string
	provider string
}

// NewAmazonS3 创建 Amazon S3 客户端实例
func NewAmazonS3(cfg OssConfig) (*AmazonS3, error) {
	region := cfg.Region
	if region == "" {
		region = "us-east-1"
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
		// 开启跳过证书校验时，替换 SDK 默认 HTTP 客户端的 TLS 配置
		awsconfig.WithHTTPClient(s3HTTPClient(cfg.InsecureSkipVerify)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		// 仅自建端点（如 RustFS）时覆盖 BaseEndpoint；纯 AWS S3 场景 endpoint 为空，
		// 设置空串端点会导致请求地址异常，必须保留 SDK 默认端点解析
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true
		}
	})

	provider := cfg.Provider
	if provider == "" {
		provider = "s3"
	}

	return &AmazonS3{
		client:   client,
		bucket:   cfg.Bucket,
		region:   region,
		domain:   cfg.Domain,
		endpoint: cfg.Endpoint,
		provider: provider,
	}, nil
}

// s3HTTPClient 构建 S3 用 HTTP 客户端：开启时基于 SDK BuildableClient 克隆并跳过 TLS 证书校验，
// 否则返回 SDK 默认客户端，避免影响全局默认 Transport
func s3HTTPClient(skipVerify bool) aws.HTTPClient {
	if !skipVerify {
		return awshttp.NewBuildableClient()
	}
	return awshttp.NewBuildableClient().WithTransportOptions(func(t *http.Transport) {
		t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // 用户显式开启的自签名证书/私有端点场景
	})
}

// Upload 上传文件到 S3，返回访问 URL
func (a *AmazonS3) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("S3读取上传数据失败: %w", err)
	}

	_, err = a.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(a.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(data))),
	})
	if err != nil {
		return "", fmt.Errorf("S3上传失败: %w", err)
	}

	// URL 拼接优先级：自定义域名 > 自建端点（RustFS 等） > AWS 默认虚拟主机地址
	if a.domain != "" {
		return strings.TrimRight(ensureScheme(a.domain, "https"), "/") + "/" + key, nil
	}
	if a.endpoint != "" {
		// 自建端点（RustFS 等）使用 Path Style 访问，URL 必须包含桶名：
		// {endpoint}/{bucket}/{key}，缺失桶名会导致 403/404。
		// 端点未带协议时默认 http（私有部署常见明文服务）
		return strings.TrimRight(ensureScheme(a.endpoint, "http"), "/") + "/" + a.bucket + "/" + key, nil
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", a.bucket, a.region, key), nil
}

// ensureScheme 为地址补齐协议前缀：已带 http:// 或 https:// 时原样返回，否则按 defaultScheme 补齐
func ensureScheme(addr, defaultScheme string) string {
	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		return addr
	}
	return defaultScheme + "://" + addr
}

// Delete 从 S3 删除文件
func (a *AmazonS3) Delete(ctx context.Context, key string) error {
	_, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("S3删除失败: %w", err)
	}
	return nil
}

// TestConnectivity 测试 S3 连通性
func (a *AmazonS3) TestConnectivity(ctx context.Context) error {
	_, err := a.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(a.bucket),
	})
	if err != nil {
		logger.Error("OSS连通测试失败", "provider", a.provider, "region", a.region, "error", err)
		return fmt.Errorf("S3连通测试失败: %w", err)
	}
	return nil
}
