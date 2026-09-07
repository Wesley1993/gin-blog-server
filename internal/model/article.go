package model

import "time"

// Article 文章模型
type Article struct {
	BaseModel
	Title        string     `gorm:"column:title;type:varchar(120);not null" json:"title"`
	CategoryID   int64      `gorm:"column:category_id;not null;default:0" json:"category_id"`
	Cover        string     `gorm:"column:cover;type:varchar(255)" json:"cover"`
	Content      string     `gorm:"column:content;type:text;not null" json:"content"`
	Tags         string     `gorm:"column:tags;type:varchar(200)" json:"tags"`
	Status       int16      `gorm:"column:status;not null;default:0" json:"status"`              // 0草稿 1已发布
	IsDeleted    int16      `gorm:"column:is_deleted;not null;default:0" json:"-"`               // 0未删 1已删
	IsRepost     int16      `gorm:"column:is_repost;not null;default:0" json:"is_repost"`        // 0原创 1转载
	RepostURL    string     `gorm:"column:repost_url;type:varchar(500)" json:"repost_url"`       // 原文链接（转载时必填）
	RepostAuthor string     `gorm:"column:repost_author;type:varchar(100)" json:"repost_author"` // 原作者
	PublishedAt  *time.Time `gorm:"column:published_at" json:"published_at"`                     // 发布时间（草稿可为空）
}

// TableName 指定表名
func (Article) TableName() string {
	return "blog_article"
}
