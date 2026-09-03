package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/service"
	pkgerrors "wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/response"
)

// ArticleHandler 文章处理器
type ArticleHandler struct {
	articleService *service.ArticleService
}

// NewArticleHandler 创建文章处理器实例
func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

// normalizeRepost 规范化转载字段：非转载时清空原文链接与原作者，转载时去除链接首尾空白
func normalizeRepost(isRepost *int16, repostURL, repostAuthor *string) {
	if *isRepost != 1 {
		*isRepost = 0
		*repostURL = ""
		*repostAuthor = ""
		return
	}
	*repostURL = strings.TrimSpace(*repostURL)
	*repostAuthor = strings.TrimSpace(*repostAuthor)
}

// Page 分页列表
// @Summary 文章分页列表
// @Tags 文章管理
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(10)
// @Param category_id query int false "分类ID"
// @Param status query int false "状态 0草稿 1已发布"
// @Success 200 {object} response.Response
// @Router /api/article/page [get]
func (h *ArticleHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.ParseInt(c.DefaultQuery("category_id", "0"), 10, 64)
	status, _ := strconv.ParseInt(c.DefaultQuery("status", "-1"), 10, 16)

	articles, total, err := h.articleService.Page(page, pageSize, categoryID, int16(status))
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "获取文章列表失败")
		return
	}

	response.SuccessPage(c, articles, total, page, pageSize)
}

// Create 创建文章
// @Summary 创建文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateArticleReq true "文章信息"
// @Success 200 {object} response.Response
// @Router /api/article/create [post]
func (h *ArticleHandler) Create(c *gin.Context) {
	var req dto.CreateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	// 转载校验：转载时原文链接必填；非转载时清空转载字段
	normalizeRepost(&req.IsRepost, &req.RepostURL, &req.RepostAuthor)
	if err := dto.ValidateRepost(req.IsRepost, req.RepostURL); err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, err.Error())
		return
	}

	// 从 JWT context 获取 userID（Phase 1 中间件注入）
	userID, _ := c.Get("userID")
	var uid int64
	if id, ok := userID.(int64); ok {
		uid = id
	}

	if err := h.articleService.Create(uid, req); err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "创建文章失败")
		return
	}

	response.Success(c)
}

// Update 更新文章
// @Summary 更新文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateArticleReq true "文章信息"
// @Success 200 {object} response.Response
// @Router /api/article/update [put]
func (h *ArticleHandler) Update(c *gin.Context) {
	var req dto.UpdateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	// 转载校验：转载时原文链接必填；非转载时清空转载字段
	normalizeRepost(&req.IsRepost, &req.RepostURL, &req.RepostAuthor)
	if err := dto.ValidateRepost(req.IsRepost, req.RepostURL); err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, err.Error())
		return
	}

	if err := h.articleService.Update(req); err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "更新文章失败")
		return
	}

	response.Success(c)
}

// Delete 逻辑删除文章
// @Summary 逻辑删除文章
// @Tags 文章管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/article/{id} [delete]
func (h *ArticleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	if err := h.articleService.Delete(id); err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "删除文章失败")
		return
	}

	response.Success(c)
}

// SearchArticle ES 全文搜索（降级 PG LIKE）
// @Summary ES全文搜索文章
// @Tags 文章管理
// @Produce json
// @Security BearerAuth
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(10)
// @Success 200 {object} response.Response
// @Router /api/article/es/search [get]
func (h *ArticleHandler) SearchArticle(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "请输入搜索关键词")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	articles, total, err := h.articleService.SearchArticle(keyword, page, pageSize)
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "搜索失败")
		return
	}

	response.SuccessPage(c, articles, total, page, pageSize)
}

// RebuildIndex 手动重建 ES 全量索引
// @Summary 手动重建ES全量索引
// @Description 仅超级管理员可用，全量从PG同步到ES
// @Tags 文章管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/article/es/rebuild [post]
func (h *ArticleHandler) RebuildIndex(c *gin.Context) {
	success, failed, err := h.articleService.RebuildIndex()
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "重建索引失败: "+err.Error())
		return
	}

	response.SuccessData(c, gin.H{
		"success": success,
		"failed":  failed,
	})
}

// PublicPage 前台文章分页列表（仅已发布，不含 content）
// @Summary 前台文章分页列表（仅已发布）
// @Description 博客前台公开接口，返回列表不包含 content 字段，节省带宽
// @Tags 博客公开接口
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(10)
// @Param category_id query int false "分类ID"
// @Param date query string false "发布日期（YYYY-MM-DD）"
// @Success 200 {object} response.Response
// @Router /api/public/article/page [get]
func (h *ArticleHandler) PublicPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.ParseInt(c.DefaultQuery("category_id", "0"), 10, 64)
	date := c.Query("date")

	articles, total, err := h.articleService.PublicPage(page, pageSize, categoryID, date)
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "获取文章列表失败")
		return
	}

	response.SuccessPage(c, articles, total, page, pageSize)
}

// PublicDates 前台文章发布日历（指定年月内有文章发布的日期列表）
// @Summary 前台文章发布日历（指定年月）
// @Description 博客前台公开接口，返回指定年月内有已发布文章的日期列表（YYYY-MM-DD）
// @Tags 博客公开接口
// @Produce json
// @Param year query int true "年份"
// @Param month query int true "月份 1-12"
// @Success 200 {object} response.Response
// @Router /api/public/article/dates [get]
func (h *ArticleHandler) PublicDates(c *gin.Context) {
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "参数错误")
		return
	}
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 1 || month > 12 {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	dates, err := h.articleService.PublicDates(year, month)
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "获取发布日历失败")
		return
	}

	response.SuccessData(c, dates)
}

// PublicDetail 前台文章详情（仅已发布，含 content）
// @Summary 前台文章详情（仅已发布）
// @Description 博客前台公开接口，根据 ID 获取已发布文章详情（含 content）
// @Tags 博客公开接口
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/public/article/{id} [get]
func (h *ArticleHandler) PublicDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "参数错误")
		return
	}

	article, err := h.articleService.PublicDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusOK, pkgerrors.ErrArticleNotFound, pkgerrors.GetMsg(pkgerrors.ErrArticleNotFound))
			return
		}
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "获取文章详情失败")
		return
	}

	response.SuccessData(c, article)
}

// PublicSearch 前台 ES 全文搜索（仅已发布，降级 PG）
// @Summary 前台ES全文搜索（仅已发布）
// @Description 博客前台公开接口，仅搜索已发布文章，ES 不可用时降级 PG LIKE
// @Tags 博客公开接口
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(10)
// @Success 200 {object} response.Response
// @Router /api/public/article/es/search [get]
func (h *ArticleHandler) PublicSearch(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.Fail(c, http.StatusOK, pkgerrors.CodeBadRequest, "请输入搜索关键词")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	articles, total, err := h.articleService.PublicSearch(keyword, page, pageSize)
	if err != nil {
		response.Fail(c, http.StatusOK, pkgerrors.CodeServerError, "搜索失败")
		return
	}

	response.SuccessPage(c, articles, total, page, pageSize)
}
