package repository

import (
	"errors"

	"gorm.io/gorm"
	"wuzhispace.com/internal/model"
)

// SiteProfileRepository 站长个人资料数据访问层
type SiteProfileRepository struct {
	DB *gorm.DB
}

// NewSiteProfileRepository 创建站长个人资料数据访问实例
func NewSiteProfileRepository(db *gorm.DB) *SiteProfileRepository {
	return &SiteProfileRepository{DB: db}
}

// Get 获取站长个人资料（取第一条，表为空时返回 nil, nil）
func (r *SiteProfileRepository) Get() (*model.SiteProfile, error) {
	var p model.SiteProfile
	err := r.DB.First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Save 保存站长个人资料（ID=0 创建，否则更新）
func (r *SiteProfileRepository) Save(p *model.SiteProfile) error {
	if p.ID == 0 {
		return r.DB.Create(p).Error
	}
	return r.DB.Save(p).Error
}
