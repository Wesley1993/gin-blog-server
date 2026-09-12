package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/model"
	"wuzhispace.com/internal/repository"
	"wuzhispace.com/pkg/elasticsearch"
	pkgerrors "wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/logger"
)

// ArticleService 文章服务
type ArticleService struct {
	BaseService
	Repo         *repository.ArticleRepository
	categoryRepo *repository.CategoryRepository
	esClient     *elasticsearch.ESClient // 可为 nil（ES 不可用时）
}

// NewArticleService 创建文章服务实例
func NewArticleService(repo *repository.ArticleRepository, categoryRepo *repository.CategoryRepository, db *gorm.DB, esClient *elasticsearch.ESClient) *ArticleService {
	return &ArticleService{
		Repo:         repo,
		categoryRepo: categoryRepo,
		BaseService:  BaseService{DB: db},
		esClient:     esClient,
	}
}

// Create 创建文章：封面为空且指定了分类时，自动使用分类图片作为默认封面。
// 兜底仅限创建场景：更新时封面传空保留原有“清空封面”语义，不做回填。
func (s *ArticleService) Create(userID int64, req dto.CreateArticleReq) error {
	now := time.Now()
	cover := req.Cover
	if cover == "" && req.CategoryID > 0 {
		if cat, err := s.categoryRepo.FindByID(req.CategoryID); err == nil && cat.Image != "" {
			cover = cat.Image
			logger.Info("文章未上传封面，使用分类图片兜底", "userID", userID, "categoryID", req.CategoryID, "cover", cover)
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("查询分类图片失败，跳过封面兜底", "categoryID", req.CategoryID, "error", err)
		}
	}

	article := &model.Article{
		Title:        req.Title,
		CategoryID:   req.CategoryID,
		Cover:        cover,
		Content:      req.Content,
		Tags:         req.Tags,
		Status:       req.Status,
		IsRepost:     req.IsRepost,
		RepostURL:    req.RepostURL,
		RepostAuthor: req.RepostAuthor,
		PublishedAt:  parsePublishedAt(req.PublishedAt, now),
		BaseModel: model.BaseModel{
			CreateTime: now,
			UpdateTime: now,
		},
	}

	if err := s.Repo.Create(article); err != nil {
		logger.Error("创建文章入库失败", "userID", userID, "title", req.Title, "error", err)
		return err
	}

	// ES 同步
	if s.esClient != nil {
		if err := s.esClient.SyncArticle(article); err != nil {
			logger.Error("ES同步文章失败", "error", err)
		}
	}

	logger.Info("创建文章成功", "userID", userID, "articleID", article.ID, "title", article.Title)
	return nil
}

// Update 更新文章
func (s *ArticleService) Update(req dto.UpdateArticleReq) error {
	article, err := s.Repo.FindByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("更新文章拒绝，文章不存在", "articleID", req.ID)
			return errors.New(pkgerrors.GetMsg(pkgerrors.ErrArticleNotFound))
		}
		logger.Error("更新文章查询失败", "articleID", req.ID, "error", err)
		return err
	}

	article.Title = req.Title
	article.CategoryID = req.CategoryID
	article.Cover = req.Cover
	article.Content = req.Content
	article.Tags = req.Tags
	article.Status = req.Status
	article.IsRepost = req.IsRepost
	article.RepostURL = req.RepostURL
	article.RepostAuthor = req.RepostAuthor
	// published_at 传了才更新，不传保持原值
	if req.PublishedAt != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", req.PublishedAt, time.Local); err == nil {
			article.PublishedAt = &t
		}
	}
	article.UpdateTime = time.Now()

	if err := s.Repo.Update(article); err != nil {
		logger.Error("更新文章失败", "articleID", article.ID, "error", err)
		return err
	}

	// ES 同步：更新后重新查询完整 article 再同步
	if s.esClient != nil {
		updated, err := s.Repo.FindByID(req.ID)
		if err == nil && updated != nil {
			if syncErr := s.esClient.SyncArticle(updated); syncErr != nil {
				logger.Error("ES更新后同步失败", "error", syncErr)
			}
		}
	}

	return nil
}

// Delete 逻辑删除文章
func (s *ArticleService) Delete(id int64) error {
	if err := s.Repo.Delete(id); err != nil {
		logger.Error("删除文章失败", "articleID", id, "error", err)
		return err
	}

	// ES 删除
	if s.esClient != nil {
		if err := s.esClient.DeleteArticle(id); err != nil {
			logger.Error("ES删除文章文档失败", "error", err)
		}
	}

	return nil
}

// Page 分页查询文章
func (s *ArticleService) Page(page, pageSize int, categoryID int64, status int16) ([]model.Article, int64, error) {
	return s.Repo.Page(page, pageSize, categoryID, status)
}

// PublicPage 前台分页查询，仅已发布文章（不含 content），支持分类与日期筛选
func (s *ArticleService) PublicPage(page, pageSize int, categoryID int64, date string) ([]model.Article, int64, error) {
	return s.Repo.FindPublishedPage(page, pageSize, categoryID, date)
}

// PublicDates 获取指定年月内有已发布文章的日期列表（前台日历高亮使用）
func (s *ArticleService) PublicDates(year, month int) ([]string, error) {
	return s.Repo.FindPublishedDates(year, month)
}

// PublicDetail 前台获取单篇已发布文章详情（含 content）
func (s *ArticleService) PublicDetail(id int64) (*model.Article, error) {
	return s.Repo.FindPublishedByID(id)
}

// PublicSearch 前台 ES 全文搜索，仅已发布文章，ES 不可用时降级 PG LIKE
func (s *ArticleService) PublicSearch(keyword string, page, pageSize int) ([]model.Article, int64, error) {
	if s.esClient != nil && s.esClient.Available {
		result, err := s.esClient.SearchPublished(keyword, page, pageSize)
		if err == nil && result != nil {
			articles := make([]model.Article, 0, len(result.List))
			for _, item := range result.List {
				articles = append(articles, model.Article{
					BaseModel: model.BaseModel{
						ID:         item.ID,
						CreateTime: parseTime(item.CreateTime),
					},
					Title:       item.Title,
					Tags:        item.Tags,
					CategoryID:  item.CategoryID,
					Status:      int16(item.Status),
					PublishedAt: parseTimePtr(item.PublishedAt),
				})
			}
			return articles, result.Total, nil
		}
		logger.Warn("ES前台搜索失败，降级PG", "keyword", keyword, "error", err)
	}
	return s.Repo.SearchPublishedByKeyword(keyword, page, pageSize)
}

// SearchArticle ES 全文搜索，ES 不可用时降级 PG LIKE
func (s *ArticleService) SearchArticle(keyword string, page, pageSize int) ([]model.Article, int64, error) {
	// 尝试 ES 搜索
	if s.esClient != nil && s.esClient.Available {
		result, err := s.esClient.Search(keyword, page, pageSize)
		if err == nil && result != nil {
			articles := make([]model.Article, 0, len(result.List))
			for _, item := range result.List {
				articles = append(articles, model.Article{
					BaseModel: model.BaseModel{
						ID:         item.ID,
						CreateTime: parseTime(item.CreateTime),
					},
					Title:       item.Title,
					Tags:        item.Tags,
					CategoryID:  item.CategoryID,
					Status:      int16(item.Status),
					PublishedAt: parseTimePtr(item.PublishedAt),
				})
			}
			return articles, result.Total, nil
		}
		logger.Warn("ES搜索失败，降级PG", "keyword", keyword, "error", err)
	}
	// 降级 PG LIKE
	return s.Repo.SearchByKeyword(keyword, (page-1)*pageSize, pageSize)
}

// RebuildIndex 全量重建 ES 索引
func (s *ArticleService) RebuildIndex() (success, failed int, err error) {
	if s.esClient == nil {
		return 0, 0, fmt.Errorf("ES客户端未初始化")
	}
	articles, err := s.Repo.FindAllNotDeleted()
	if err != nil {
		return 0, 0, err
	}
	return s.esClient.RebuildIndex(articles)
}

// Publish 一键发布草稿：将 status 从 0 改为 1，设置 published_at，同步 ES
func (s *ArticleService) Publish(id int64) error {
	article, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("发布文章拒绝，文章不存在", "articleID", id)
			return errors.New(pkgerrors.GetMsg(pkgerrors.ErrArticleNotFound))
		}
		logger.Error("发布文章查询失败", "articleID", id, "error", err)
		return err
	}

	if article.Status == 1 {
		return errors.New("该文章已发布，无需重复操作")
	}

	now := time.Now()
	article.Status = 1
	article.PublishedAt = &now
	article.UpdateTime = now

	if err := s.Repo.Update(article); err != nil {
		logger.Error("发布文章失败", "articleID", id, "error", err)
		return err
	}

	// ES 同步
	if s.esClient != nil {
		if syncErr := s.esClient.SyncArticle(article); syncErr != nil {
			logger.Error("发布后ES同步失败", "articleID", id, "error", syncErr)
		}
	}

	logger.Info("文章发布成功", "articleID", id, "title", article.Title)
	return nil
}

// parseTime 解析 ES 返回的时间字符串
func parseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}

// parsePublishedAt 解析请求中的发布时间字符串：为空或解析失败时回退到 fallback
func parsePublishedAt(s string, fallback time.Time) *time.Time {
	if s != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
			return &t
		}
	}
	t := fallback
	return &t
}

// parseTimePtr 解析 ES 返回的时间字符串为指针，空串或解析失败返回 nil
func parseTimePtr(s string) *time.Time {
	if s == "" {
		return nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		return &t
	}
	return nil
}
