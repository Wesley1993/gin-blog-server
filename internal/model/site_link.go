package model

// SiteLink 常用网站（友情链接）模型
type SiteLink struct {
	BaseModel
	Name        string `gorm:"column:name;type:varchar(50);not null" json:"name"`             // 网站名称
	URL         string `gorm:"column:url;type:varchar(255);not null" json:"url"`              // 链接地址
	Icon        string `gorm:"column:icon;type:varchar(500);not null;default:''" json:"icon"` // 图标图片 URL
	Description string `gorm:"column:description;type:varchar(200)" json:"description"`       // 可选描述
	Sort        int    `gorm:"column:sort;not null;default:0" json:"sort"`                    // 排序，小靠前
	Status      int16  `gorm:"column:status;not null;default:1" json:"status"`                // 1 启用 0 停用
}

// TableName 指定表名
func (SiteLink) TableName() string {
	return "blog_link"
}
