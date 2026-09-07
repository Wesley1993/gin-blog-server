package dto

import (
	"errors"
	"strings"
)

// CreateArticleReq 创建文章请求
type CreateArticleReq struct {
	Title        string `json:"title" binding:"required"`
	CategoryID   int64  `json:"category_id" binding:"required"`
	Cover        string `json:"cover"`
	Content      string `json:"content" binding:"required"`
	Tags         string `json:"tags"`
	Status       int16  `json:"status"`
	IsRepost     int16  `json:"is_repost"`     // 0原创 1转载
	RepostURL    string `json:"repost_url"`    // 原文链接（转载时必填）
	RepostAuthor string `json:"repost_author"` // 原作者
	PublishedAt  string `json:"published_at"`  // 可选，格式 "2006-01-02 15:04:05"，未传则默认当前时间
}

// UpdateArticleReq 更新文章请求
type UpdateArticleReq struct {
	ID           int64  `json:"id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	CategoryID   int64  `json:"category_id" binding:"required"`
	Cover        string `json:"cover"`
	Content      string `json:"content" binding:"required"`
	Tags         string `json:"tags"`
	Status       int16  `json:"status"`
	IsRepost     int16  `json:"is_repost"`     // 0原创 1转载
	RepostURL    string `json:"repost_url"`    // 原文链接（转载时必填）
	RepostAuthor string `json:"repost_author"` // 原作者
	PublishedAt  string `json:"published_at"`  // 可选，格式 "2006-01-02 15:04:05"，传了才更新
}

// ValidateRepost 校验转载字段：转载时原文链接必填，返回错误表示参数不合法
func ValidateRepost(isRepost int16, repostURL string) error {
	if isRepost == 1 && strings.TrimSpace(repostURL) == "" {
		return errors.New("转载文章必须填写原文链接")
	}
	return nil
}
