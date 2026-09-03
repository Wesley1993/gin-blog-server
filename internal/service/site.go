package service

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/model"
	"wuzhispace.com/internal/repository"
	"wuzhispace.com/pkg/logger"
	pkgredis "wuzhispace.com/pkg/redis"
	"wuzhispace.com/pkg/storage"
)

const siteConfigCacheKey = "site:config"
const siteConfigCacheTTL = 1 * time.Hour

// SiteConfigService 站点配置服务
type SiteConfigService struct {
	BaseService
	Repo *repository.SiteConfigRepository
}

// NewSiteConfigService 创建站点配置服务实例
func NewSiteConfigService(repo *repository.SiteConfigRepository, db *gorm.DB, redis *pkgredis.RedisClient) *SiteConfigService {
	return &SiteConfigService{
		Repo:        repo,
		BaseService: BaseService{DB: db, Redis: redis},
	}
}

// GetConfig 获取站点配置（优先读 Redis 缓存）
func (s *SiteConfigService) GetConfig() (*model.SiteConfig, error) {
	// 尝试从 Redis 读取缓存
	if s.Redis != nil {
		cached, err := s.Redis.Get(siteConfigCacheKey)
		if err == nil && cached != "" {
			var cfg model.SiteConfig
			if jsonErr := json.Unmarshal([]byte(cached), &cfg); jsonErr == nil {
				// 兼容旧结构缓存：DB 中 oss_provider 为 NOT NULL DEFAULT 'aliyun'，
				// provider 为空说明缓存是无该字段的旧版本写入，作废回源 DB，
				// 避免厂商路由回落默认值
				if cfg.OssProvider != "" {
					return &cfg, nil
				}
				logger.Warn("SiteConfig缓存为旧结构（缺provider字段），作废回源DB")
				if _, delErr := s.Redis.Del(siteConfigCacheKey); delErr != nil {
					logger.Error("SiteConfig删除旧缓存失败", "error", delErr)
				}
			}
		}
	}

	// 从 DB 读取
	cfg, err := s.Repo.Get()
	if err != nil {
		return nil, err
	}

	// 写入 Redis 缓存
	if s.Redis != nil {
		data, jsonErr := json.Marshal(cfg)
		if jsonErr == nil {
			if cacheErr := s.Redis.SetWithTTL(siteConfigCacheKey, string(data), siteConfigCacheTTL); cacheErr != nil {
				logger.Error("SiteConfig写入Redis缓存失败", "error", cacheErr)
			}
		}
	}

	return cfg, nil
}

// SaveConfig 保存站点配置，保存后删除 Redis 缓存
func (s *SiteConfigService) SaveConfig(cfg *model.SiteConfig) error {
	cfg.UpdateTime = time.Now()

	if err := s.Repo.Save(cfg); err != nil {
		logger.Error("SiteConfig保存失败", "error", err)
		return err
	}

	// 删除 Redis 缓存
	if s.Redis != nil {
		if _, err := s.Redis.Del(siteConfigCacheKey); err != nil {
			logger.Error("SiteConfig删除Redis缓存失败", "error", err)
		}
	}

	return nil
}

// GetOssConfigSource 直接从 DB 读取站点配置（不走缓存，用于配置保存链路取已有记录，
// 避免旧结构缓存中缺失 provider 字段导致存储实例缓存清理用错 key）
func (s *SiteConfigService) GetOssConfigSource() (*model.SiteConfig, error) {
	return s.Repo.Get()
}

// GetOssConfig 从 DB 读取 OSS 配置并转换为 storage.OssConfig
func (s *SiteConfigService) GetOssConfig() (*storage.OssConfig, error) {
	cfg, err := s.Repo.Get()
	if err != nil {
		logger.Error("读取OSS配置失败", "error", err)
		return nil, err
	}

	return &storage.OssConfig{
		AccessKey:          cfg.OssAccessKey,
		SecretKey:          cfg.OssSecretKey,
		Bucket:             cfg.OssBucket,
		Endpoint:           cfg.OssEndpoint,
		Domain:             cfg.OssDomain,
		Provider:           cfg.OssProvider,
		Region:             cfg.OssRegion,
		InsecureSkipVerify: cfg.OssInsecure != 0,
	}, nil
}

// GetStats 获取站点运行统计（建站日期、运行天数、文章/分类/用户数）
func (s *SiteConfigService) GetStats() (*dto.SiteStatsResp, error) {
	resp := &dto.SiteStatsResp{}

	// 建站日期与运行天数（直接查 DB，不走缓存，保证 founded_at 实时）
	cfg, err := s.Repo.Get()
	if err == nil && cfg != nil {
		resp.FoundedAt = cfg.FoundedAt
		if founded, parseErr := time.ParseInLocation("2006-01-02", cfg.FoundedAt, time.Local); parseErr == nil {
			days := int64(time.Since(founded).Hours() / 24)
			if days < 0 {
				days = 0
			}
			resp.RunningDays = days
		}
	}

	// 文章总数（未删除）
	if err := s.DB.Model(&model.Article{}).Where("is_deleted = 0").Count(&resp.ArticleCount).Error; err != nil {
		return nil, err
	}

	// 分类数
	if err := s.DB.Model(&model.Category{}).Count(&resp.CategoryCount).Error; err != nil {
		return nil, err
	}

	// 用户数
	if err := s.DB.Model(&model.User{}).Count(&resp.UserCount).Error; err != nil {
		return nil, err
	}

	return resp, nil
}
