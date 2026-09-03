package model

// User 用户模型
type User struct {
	BaseModel
	Username string `gorm:"column:username;type:varchar(50);not null;uniqueIndex" json:"username"`
	Password string `gorm:"column:password;type:varchar(100);not null" json:"-"`
	Nickname string `gorm:"column:nickname;type:varchar(50)" json:"nickname"`
	Avatar   string `gorm:"column:avatar;type:varchar(255)" json:"avatar"` // 头像URL
	Bio      string `gorm:"column:bio;type:varchar(500)" json:"bio"`       // 个人简介
	RoleID   int64  `gorm:"column:role_id;not null" json:"role_id"`
	Status   int16  `gorm:"column:status;not null;default:1" json:"status"`
}

// TableName 指定表名
func (User) TableName() string { return "sys_user" }
