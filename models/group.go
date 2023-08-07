package models

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Desc     string
	OwnerId  uint // 关系信息
	TargetId uint // 对应谁
	Type     int  // 对应类型 0 1 3
}

func (table *Group) TableName() string {
	return "group"
}
