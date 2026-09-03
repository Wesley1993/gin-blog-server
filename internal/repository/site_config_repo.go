package repository

import (
	"gorm.io/gorm"
	"wuzhispace.com/internal/model"
)

// SiteConfigRepository 站点配置数据访问层
type SiteConfigRepository struct {
	DB *gorm.DB
}

// NewSiteConfigRepository 创建站点配置数据访问实例
func NewSiteConfigRepository(db *gorm.DB) *SiteConfigRepository {
	return &SiteConfigRepository{DB: db}
}

// Get 获取站点配置（取第一条）
func (r *SiteConfigRepository) Get() (*model.SiteConfig, error) {
	var cfg model.SiteConfig
	err := r.DB.First(&cfg).Error
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save 保存站点配置（ID=0 创建，否则更新）
func (r *SiteConfigRepository) Save(cfg *model.SiteConfig) error {
	if cfg.ID == 0 {
		return r.DB.Create(cfg).Error
	}
	return r.DB.Save(cfg).Error
}
