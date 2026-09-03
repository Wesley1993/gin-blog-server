package service

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/model"
	"wuzhispace.com/internal/repository"
)

// 常用网站业务错误
var ErrLinkNotFound = errors.New("link not found")

// SiteLinkService 常用网站服务
type SiteLinkService struct {
	Repo *repository.SiteLinkRepository
}

// NewSiteLinkService 创建常用网站服务实例
func NewSiteLinkService(repo *repository.SiteLinkRepository) *SiteLinkService {
	return &SiteLinkService{Repo: repo}
}

// List 获取全部记录（含停用），按 sort 升序、id 升序
func (s *SiteLinkService) List() ([]model.SiteLink, error) {
	links, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}
	if links == nil {
		links = []model.SiteLink{}
	}
	return links, nil
}

// PublicList 获取启用中的常用网站（公开接口契约结构）
func (s *SiteLinkService) PublicList() ([]*dto.SiteLinkItem, error) {
	links, err := s.Repo.FindPublished()
	if err != nil {
		return nil, err
	}

	items := make([]*dto.SiteLinkItem, 0, len(links))
	for i := range links {
		items = append(items, &dto.SiteLinkItem{
			ID:          links[i].ID,
			Name:        links[i].Name,
			URL:         links[i].URL,
			Icon:        links[i].Icon,
			Description: links[i].Description,
			Sort:        links[i].Sort,
		})
	}
	return items, nil
}

// Create 新增常用网站
func (s *SiteLinkService) Create(req dto.CreateSiteLinkReq) error {
	now := time.Now()
	link := &model.SiteLink{
		Name:        req.Name,
		URL:         req.URL,
		Icon:        req.Icon,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
		BaseModel: model.BaseModel{
			CreateTime: now,
			UpdateTime: now,
		},
	}
	return s.Repo.Create(link)
}

// Update 更新常用网站
func (s *SiteLinkService) Update(id int64, req dto.UpdateSiteLinkReq) error {
	link, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLinkNotFound
		}
		return err
	}

	link.Name = req.Name
	link.URL = req.URL
	link.Icon = req.Icon
	link.Description = req.Description
	link.Sort = req.Sort
	link.Status = req.Status
	link.UpdateTime = time.Now()

	return s.Repo.Update(link)
}

// Delete 删除常用网站
func (s *SiteLinkService) Delete(id int64) error {
	if _, err := s.Repo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLinkNotFound
		}
		return err
	}
	return s.Repo.Delete(id)
}
