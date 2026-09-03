package model

import "encoding/json"

// SiteProfile 站长个人资料模型（单行表）
type SiteProfile struct {
	BaseModel
	Name     string          `gorm:"column:name;type:varchar(50)" json:"name"`        // 站长姓名
	Avatar   string          `gorm:"column:avatar;type:varchar(255)" json:"avatar"`   // 头像
	Title    string          `gorm:"column:title;type:varchar(100)" json:"title"`     // 头衔/一句话介绍
	Bio      string          `gorm:"column:bio;type:text" json:"bio"`                 // 个人简历/自序
	Github   string          `gorm:"column:github;type:varchar(100)" json:"github"`   // GitHub 地址
	Email    string          `gorm:"column:email;type:varchar(100)" json:"email"`     // 邮箱
	Wechat   string          `gorm:"column:wechat;type:varchar(100)" json:"wechat"`   // 微信号
	QQ       string          `gorm:"column:qq;type:varchar(100)" json:"qq"`           // QQ 号
	Address  string          `gorm:"column:address;type:varchar(200)" json:"address"` // 地址
	Skills   json.RawMessage `gorm:"column:skills;type:jsonb" json:"skills"`          // 技能列表 jsonb
	Projects json.RawMessage `gorm:"column:projects;type:jsonb" json:"projects"`      // 项目经历列表 jsonb
}

// TableName 指定表名
func (SiteProfile) TableName() string {
	return "blog_profile"
}
