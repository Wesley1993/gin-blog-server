package model

// Menu 菜单权限模型
type Menu struct {
	BaseModel
	ParentID int64  `gorm:"column:parent_id;not null;default:0" json:"parent_id"`
	MenuName string `gorm:"column:menu_name;type:varchar(50);not null" json:"menu_name"`
	MenuType int16  `gorm:"column:menu_type;not null" json:"menu_type"` // 1目录 2页面 3按钮
	Path     string `gorm:"column:path;type:varchar(100)" json:"path"`
	Perms    string `gorm:"column:perms;type:varchar(100)" json:"perms"`
	Sort     int    `gorm:"column:sort;not null;default:0" json:"sort"`
	Status   int16  `gorm:"column:status;not null;default:1" json:"status"`
	Children []Menu `gorm:"-" json:"children,omitempty"`
}

// TableName 指定表名
func (Menu) TableName() string { return "sys_menu" }
