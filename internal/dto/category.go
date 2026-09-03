package dto

// CreateCategoryReq 创建分类请求
type CreateCategoryReq struct {
	Name     string `json:"name" binding:"required"`
	ParentID int64  `json:"parent_id"`
	Sort     int    `json:"sort"`
	Status   int16  `json:"status"`
	Image    string `json:"image"`
}

// UpdateCategoryReq 更新分类请求
type UpdateCategoryReq struct {
	ID       int64  `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	ParentID int64  `json:"parent_id"`
	Sort     int    `json:"sort"`
	Status   int16  `json:"status"`
	Image    string `json:"image"`
}
