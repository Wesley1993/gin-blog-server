package dto

import "encoding/json"

// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	RoleName    string          `json:"role_name" binding:"required"`
	MenuIDs     json.RawMessage `json:"menu_ids" swaggertype:"array,integer"`          // 菜单ID数组
	ButtonPerms json.RawMessage `json:"button_perms" swaggertype:"array,string"`       // 按钮权限标识数组
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	ID          int64           `json:"id" binding:"required"`
	RoleName    string          `json:"role_name" binding:"required"`
	MenuIDs     json.RawMessage `json:"menu_ids" swaggertype:"array,integer"`          // 菜单ID数组
	ButtonPerms json.RawMessage `json:"button_perms" swaggertype:"array,string"`       // 按钮权限标识数组
}
