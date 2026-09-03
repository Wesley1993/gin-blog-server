package dto

// ProfileDTO 个人信息DTO
type ProfileDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	RoleID   int64  `json:"role_id"`
}

// UpdateProfileReq 更新个人信息请求
type UpdateProfileReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

// SiteStatsResp 站点运行统计响应
type SiteStatsResp struct {
	FoundedAt     string `json:"founded_at"`     // 建站日期
	RunningDays   int64  `json:"running_days"`   // 运行天数
	ArticleCount  int64  `json:"article_count"`  // 文章总数（未删除）
	CategoryCount int64  `json:"category_count"` // 分类数
	UserCount     int64  `json:"user_count"`     // 用户数
}
