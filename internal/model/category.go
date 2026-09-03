package model

// Category 分类模型
type Category struct {
	BaseModel
	Name     string     `gorm:"column:name;type:varchar(60);not null" json:"name"`
	ParentID int64      `gorm:"column:parent_id;not null;default:0" json:"parent_id"`
	Sort     int        `gorm:"column:sort;not null;default:0" json:"sort"`
	Status   int16      `gorm:"column:status;not null;default:1" json:"status"`
	Image    string     `gorm:"column:image;type:varchar(500);not null;default:''" json:"image"`
	Children []Category `gorm:"-" json:"children,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "blog_category"
}
