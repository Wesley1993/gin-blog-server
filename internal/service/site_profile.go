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
)

const siteProfileCacheKey = "site:profile"
const siteProfileCacheTTL = 1 * time.Hour

// SiteProfileService 站长个人资料服务（区别于登录用户个人信息 ProfileService）
type SiteProfileService struct {
	BaseService
	Repo *repository.SiteProfileRepository
}

// NewSiteProfileService 创建站长个人资料服务实例
func NewSiteProfileService(repo *repository.SiteProfileRepository, db *gorm.DB, redis *pkgredis.RedisClient) *SiteProfileService {
	return &SiteProfileService{
		Repo:        repo,
		BaseService: BaseService{DB: db, Redis: redis},
	}
}

// GetProfile 获取站长个人资料（优先读 Redis 缓存；表为空时返回空默认值，不报错）
func (s *SiteProfileService) GetProfile() (*dto.SiteProfileResp, error) {
	// 尝试从 Redis 读取缓存
	if s.Redis != nil {
		cached, err := s.Redis.Get(siteProfileCacheKey)
		if err == nil && cached != "" {
			var resp dto.SiteProfileResp
			if jsonErr := json.Unmarshal([]byte(cached), &resp); jsonErr == nil {
				return &resp, nil
			}
		}
	}

	// 从 DB 读取
	p, err := s.Repo.Get()
	if err != nil {
		return nil, err
	}

	resp := toSiteProfileResp(p)

	// 写入 Redis 缓存
	if s.Redis != nil {
		data, jsonErr := json.Marshal(resp)
		if jsonErr == nil {
			if cacheErr := s.Redis.SetWithTTL(siteProfileCacheKey, string(data), siteProfileCacheTTL); cacheErr != nil {
				logger.Error("SiteProfile写入Redis缓存失败", "error", cacheErr)
			}
		}
	}

	return resp, nil
}

// SaveProfile 整体覆盖保存站长个人资料，保存后删除 Redis 缓存
func (s *SiteProfileService) SaveProfile(req *dto.SaveSiteProfileReq) error {
	if req.Skills == nil {
		req.Skills = make([]*dto.SiteProfileSkill, 0)
	}
	if req.Projects == nil {
		req.Projects = make([]*dto.SiteProfileProject, 0)
	}

	skills, err := json.Marshal(req.Skills)
	if err != nil {
		logger.Error("SiteProfile序列化技能失败", "error", err)
		return err
	}
	projects, err := json.Marshal(req.Projects)
	if err != nil {
		logger.Error("SiteProfile序列化项目失败", "error", err)
		return err
	}

	p := &model.SiteProfile{
		Name:     req.Name,
		Avatar:   req.Avatar,
		Title:    req.Title,
		Bio:      req.Bio,
		Github:   req.Contacts.Github,
		Email:    req.Contacts.Email,
		Wechat:   req.Contacts.Wechat,
		QQ:       req.Contacts.QQ,
		Address:  req.Contacts.Address,
		Skills:   skills,
		Projects: projects,
	}

	// 先获取已有记录（用于获取 ID，实现覆盖更新）
	existing, err := s.Repo.Get()
	if err != nil {
		logger.Error("SiteProfile查询已有记录失败", "error", err)
		return err
	}
	if existing != nil {
		p.ID = existing.ID
	}

	p.UpdateTime = time.Now()
	if p.ID == 0 {
		p.CreateTime = time.Now()
	}

	if err := s.Repo.Save(p); err != nil {
		logger.Error("SiteProfile保存失败", "error", err)
		return err
	}

	// 删除 Redis 缓存
	if s.Redis != nil {
		if _, err := s.Redis.Del(siteProfileCacheKey); err != nil {
			logger.Error("SiteProfile删除Redis缓存失败", "error", err)
		}
	}

	logger.Info("SiteProfile保存成功", "name", req.Name)
	return nil
}

// toSiteProfileResp 将模型转换为响应结构（nil 时返回空默认值，反序列化 jsonb 字段）
func toSiteProfileResp(p *model.SiteProfile) *dto.SiteProfileResp {
	resp := &dto.SiteProfileResp{
		Skills:   make([]*dto.SiteProfileSkill, 0),
		Projects: make([]*dto.SiteProfileProject, 0),
	}
	if p == nil {
		return resp
	}

	resp.Name = p.Name
	resp.Avatar = p.Avatar
	resp.Title = p.Title
	resp.Bio = p.Bio
	resp.Contacts = dto.SiteProfileContact{
		Github:  p.Github,
		Email:   p.Email,
		Wechat:  p.Wechat,
		QQ:      p.QQ,
		Address: p.Address,
	}

	if len(p.Skills) > 0 {
		_ = json.Unmarshal(p.Skills, &resp.Skills)
	}
	if len(p.Projects) > 0 {
		_ = json.Unmarshal(p.Projects, &resp.Projects)
	}
	if resp.Skills == nil {
		resp.Skills = make([]*dto.SiteProfileSkill, 0)
	}
	if resp.Projects == nil {
		resp.Projects = make([]*dto.SiteProfileProject, 0)
	}

	return resp
}
