package service

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/model"
	"wuzhispace.com/internal/repository"
	"wuzhispace.com/pkg/logger"
)

// 分类业务错误
var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryInUse    = errors.New("category in use")
)

// CategoryService 分类服务
type CategoryService struct {
	Repo *repository.CategoryRepository
	DB   *gorm.DB
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(repo *repository.CategoryRepository, db *gorm.DB) *CategoryService {
	return &CategoryService{
		Repo: repo,
		DB:   db,
	}
}

// Tree 获取分类树
func (s *CategoryService) Tree() ([]model.Category, error) {
	cats, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}
	return s.Repo.BuildTree(cats), nil
}

// Create 新增分类
func (s *CategoryService) Create(req dto.CreateCategoryReq) error {
	now := time.Now()
	cat := &model.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		Status:   req.Status,
		Image:    req.Image,
		BaseModel: model.BaseModel{
			CreateTime: now,
			UpdateTime: now,
		},
	}
	if err := s.Repo.Create(cat); err != nil {
		logger.Error("创建分类失败", "name", req.Name, "error", err)
		return err
	}
	logger.Info("创建分类成功", "categoryID", cat.ID, "name", cat.Name)
	return nil
}

// Update 更新分类
func (s *CategoryService) Update(req dto.UpdateCategoryReq) error {
	cat, err := s.Repo.FindByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("更新分类拒绝，分类不存在", "categoryID", req.ID)
			return ErrCategoryNotFound
		}
		logger.Error("更新分类查询失败", "categoryID", req.ID, "error", err)
		return err
	}

	cat.Name = req.Name
	cat.ParentID = req.ParentID
	cat.Sort = req.Sort
	cat.Status = req.Status
	cat.Image = req.Image
	cat.UpdateTime = time.Now()

	if err := s.Repo.Update(cat); err != nil {
		logger.Error("更新分类失败", "categoryID", cat.ID, "error", err)
		return err
	}
	logger.Info("更新分类成功", "categoryID", cat.ID, "name", cat.Name)
	return nil
}

// Delete 删除分类（先检查是否被文章引用）
func (s *CategoryService) Delete(id int64) error {
	_, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("删除分类拒绝，分类不存在", "categoryID", id)
			return ErrCategoryNotFound
		}
		logger.Error("删除分类查询失败", "categoryID", id, "error", err)
		return err
	}

	count, err := s.Repo.CountArticlesByCategoryID(id)
	if err != nil {
		logger.Error("删除分类统计文章引用失败", "categoryID", id, "error", err)
		return err
	}
	if count > 0 {
		logger.Warn("删除分类拒绝，分类被文章引用", "categoryID", id, "articleCount", count)
		return ErrCategoryInUse
	}

	if err := s.Repo.Delete(id); err != nil {
		logger.Error("删除分类失败", "categoryID", id, "error", err)
		return err
	}
	logger.Info("删除分类成功", "categoryID", id)
	return nil
}
