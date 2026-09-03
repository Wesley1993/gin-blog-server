package dto

// CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	MenuName string `json:"menu_name" binding:"required"`
	ParentID int64  `json:"parent_id"`
	MenuType int    `json:"menu_type" binding:"required,oneof=1 2 3"` // 1目录 2页面 3按钮
	Path     string `json:"path"`
	Perms    string `json:"perms"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
}

// UpdateMenuReq 更新菜单请求
type UpdateMenuReq struct {
	ID       int64  `json:"id" binding:"required"`
	MenuName string `json:"menu_name"`
	ParentID int64  `json:"parent_id"`
	MenuType int    `json:"menu_type"`
	Path     string `json:"path"`
	Perms    string `json:"perms"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
}
