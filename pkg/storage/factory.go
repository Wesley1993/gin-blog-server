package storage

import (
	"fmt"
	"sync"

	"wuzhispace.com/pkg/logger"
)

var (
	instancePool sync.Map
)

// instanceCacheKey 生成实例缓存 key（含 region 与 insecure，避免 S3 等按 region 路由的厂商变更区域、
// 或切换跳过证书校验开关后命中旧实例）
func instanceCacheKey(provider, endpoint, bucket, region string, insecure bool) string {
	return fmt.Sprintf("%s:%s:%s:%s:%t", provider, endpoint, bucket, region, insecure)
}

// NewStorage 根据 provider 创建对应的存储实例（带实例缓存）
func NewStorage(provider string, cfg OssConfig) (ObjectStorage, error) {
	// provider 为空时回落阿里云（兼容历史配置），必须告警暴露路由异常
	if provider == "" {
		logger.Warn("存储厂商未配置，回落阿里云", "bucket", cfg.Bucket)
		provider = "aliyun"
	}

	// 生成缓存 key（包含 region 与 insecure，防止区域或证书校验开关变更后复用旧客户端）
	cacheKey := instanceCacheKey(provider, cfg.Endpoint, cfg.Bucket, cfg.Region, cfg.InsecureSkipVerify)

	if instance, ok := instancePool.Load(cacheKey); ok {
		return instance.(ObjectStorage), nil
	}

	cfg.Provider = provider
	var storage ObjectStorage
	var err error

	switch provider {
	case "aliyun":
		storage, err = NewAliyunOSS(cfg)
	case "tencent":
		storage, err = NewTencentCOS(cfg)
	case "qiniu":
		storage, err = NewQiniuOSS(cfg)
	case "s3":
		storage, err = NewAmazonS3(cfg)
	case "rustfs":
		// RustFS 为 S3 兼容协议的自建对象存储，复用 S3 客户端（自建端点走 path style）
		storage, err = NewAmazonS3(cfg)
	default:
		logger.Error("存储厂商不支持", "provider", provider)
		return nil, fmt.Errorf("unsupported storage provider: %s", provider)
	}

	if err != nil {
		logger.Error("存储客户端创建失败", "provider", provider, "error", err)
		return nil, err
	}

	instancePool.Store(cacheKey, storage)
	return storage, nil
}

// InvalidateInstance 清除指定配置的缓存实例（配置变更时调用）
func InvalidateInstance(provider, endpoint, bucket, region string, insecure bool) {
	if provider == "" {
		// 兼容旧调用：历史配置无 provider 时同时清理阿里云回落实例的 key
		instancePool.Delete(instanceCacheKey("aliyun", endpoint, bucket, region, insecure))
		return
	}
	instancePool.Delete(instanceCacheKey(provider, endpoint, bucket, region, insecure))
}
