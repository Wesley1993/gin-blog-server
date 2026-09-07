package repository

import (
	"gorm.io/gorm"
	"wuzhispace.com/internal/model"
)

// ArticleRepository 文章数据访问层
type ArticleRepository struct {
	DB *gorm.DB
}

// NewArticleRepository 创建文章数据访问实例
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{DB: db}
}

// Create 创建文章
func (r *ArticleRepository) Create(article *model.Article) error {
	return r.DB.Create(article).Error
}

// FindByID 根据 ID 查询文章（未删除）
func (r *ArticleRepository) FindByID(id int64) (*model.Article, error) {
	var article model.Article
	err := r.DB.Where("id = ? AND is_deleted = 0", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// Update 更新文章
func (r *ArticleRepository) Update(article *model.Article) error {
	return r.DB.Save(article).Error
}

// Delete 逻辑删除文章
func (r *ArticleRepository) Delete(id int64) error {
	return r.DB.Model(&model.Article{}).Where("id = ? AND is_deleted = 0", id).Update("is_deleted", 1).Error
}

// Page 分页查询，支持按 category_id 和 status 筛选
func (r *ArticleRepository) Page(page, pageSize int, categoryID int64, status int16) ([]model.Article, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := r.DB.Model(&model.Article{}).Where("is_deleted = 0")
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []model.Article
	offset := (page - 1) * pageSize
	if err := query.Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// FindPublishedPage 分页查询已发布文章（不含 content 字段，节省带宽），支持分类与日期筛选
func (r *ArticleRepository) FindPublishedPage(page, pageSize int, categoryID int64, date string) ([]model.Article, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := r.DB.Model(&model.Article{}).Where("is_deleted = 0 AND status = 1")
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if date != "" {
		query = query.Where("TO_CHAR(published_at, 'YYYY-MM-DD') = ?", date)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []model.Article
	offset := (page - 1) * pageSize
	err := query.Select("id, title, category_id, cover, tags, status, is_repost, repost_url, repost_author, published_at, create_time, update_time").
		Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// FindPublishedByID 根据 ID 查询单篇已发布文章（含 content）
func (r *ArticleRepository) FindPublishedByID(id int64) (*model.Article, error) {
	var article model.Article
	err := r.DB.Where("id = ? AND is_deleted = 0 AND status = 1", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// SearchPublishedByKeyword PG LIKE 降级搜索（仅已发布，不含 content 字段）
func (r *ArticleRepository) SearchPublishedByKeyword(keyword string, page, pageSize int) ([]model.Article, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	like := "%" + keyword + "%"
	query := r.DB.Model(&model.Article{}).
		Where("is_deleted = 0 AND status = 1 AND (title LIKE ? OR content LIKE ? OR tags LIKE ?)", like, like, like)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []model.Article
	offset := (page - 1) * pageSize
	err := query.Select("id, title, category_id, cover, tags, status, is_repost, repost_url, repost_author, published_at, create_time, update_time").
		Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// FindPublishedDates 查询指定年月内有已发布文章的日期列表（格式 YYYY-MM-DD，前台日历高亮使用）
func (r *ArticleRepository) FindPublishedDates(year, month int) ([]string, error) {
	var dates []string
	sql := `SELECT DISTINCT TO_CHAR(published_at, 'YYYY-MM-DD')
		FROM blog_article
		WHERE is_deleted = 0 AND status = 1
		  AND EXTRACT(YEAR FROM published_at) = ?
		  AND EXTRACT(MONTH FROM published_at) = ?
		ORDER BY 1`
	if err := r.DB.Raw(sql, year, month).Scan(&dates).Error; err != nil {
		return nil, err
	}
	return dates, nil
}

// FindAllNotDeleted 查所有未删除文章（用于 ES 重建索引）
func (r *ArticleRepository) FindAllNotDeleted() ([]model.Article, error) {
	var articles []model.Article
	err := r.DB.Where("is_deleted = 0").Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// SearchByKeyword PG LIKE 降级搜索
func (r *ArticleRepository) SearchByKeyword(keyword string, page, pageSize int) ([]model.Article, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	like := "%" + keyword + "%"
	query := r.DB.Where("is_deleted = 0 AND (title LIKE ? OR content LIKE ? OR tags LIKE ?)", like, like, like)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []model.Article
	offset := (page - 1) * pageSize
	if err := query.Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
