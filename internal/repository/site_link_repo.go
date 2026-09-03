package repository

import (
	"gorm.io/gorm"
	"wuzhispace.com/internal/model"
)

// SiteLinkRepository 常用网站数据访问层
type SiteLinkRepository struct {
	DB *gorm.DB
}

// NewSiteLinkRepository 创建常用网站数据访问实例
func NewSiteLinkRepository(db *gorm.DB) *SiteLinkRepository {
	return &SiteLinkRepository{DB: db}
}

// FindAll 查询全部记录（含停用），按 sort 升序、id 升序
func (r *SiteLinkRepository) FindAll() ([]model.SiteLink, error) {
	var links []model.SiteLink
	err := r.DB.Order("sort ASC, id ASC").Find(&links).Error
	if err != nil {
		return nil, err
	}
	return links, nil
}

// FindPublished 查询启用中的记录，按 sort 升序、id 升序
func (r *SiteLinkRepository) FindPublished() ([]model.SiteLink, error) {
	var links []model.SiteLink
	err := r.DB.Where("status = 1").Order("sort ASC, id ASC").Find(&links).Error
	if err != nil {
		return nil, err
	}
	return links, nil
}

// FindByID 根据 ID 查询
func (r *SiteLinkRepository) FindByID(id int64) (*model.SiteLink, error) {
	var link model.SiteLink
	err := r.DB.Where("id = ?", id).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// Create 创建
func (r *SiteLinkRepository) Create(link *model.SiteLink) error {
	return r.DB.Create(link).Error
}

// Update 更新
func (r *SiteLinkRepository) Update(link *model.SiteLink) error {
	return r.DB.Save(link).Error
}

// Delete 物理删除
func (r *SiteLinkRepository) Delete(id int64) error {
	return r.DB.Where("id = ?", id).Delete(&model.SiteLink{}).Error
}
