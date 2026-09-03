package model

import "encoding/json"

// Role 角色模型
type Role struct {
	BaseModel
	RoleName    string          `gorm:"column:role_name;type:varchar(50);not null" json:"role_name"`
	MenuIDs     json.RawMessage `gorm:"column:menu_ids;type:jsonb" json:"menu_ids"`
	ButtonPerms json.RawMessage `gorm:"column:button_perms;type:jsonb" json:"button_perms"`
	IsSuper     int16           `gorm:"column:is_super;not null;default:0" json:"is_super"`
}

// TableName 指定表名
func (Role) TableName() string { return "sys_role" }
