package model

import "time"

// BaseModel 公共基础模型
type BaseModel struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp without time zone;not null" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp without time zone;not null" json:"update_time"`
}
