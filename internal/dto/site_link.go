package dto

// SiteLinkItem 常用网站展示项（公开接口契约字段）
type SiteLinkItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

// CreateSiteLinkReq 新增常用网站请求
type CreateSiteLinkReq struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      int16  `json:"status"`
}

// UpdateSiteLinkReq 更新常用网站请求
type UpdateSiteLinkReq struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      int16  `json:"status"`
}
