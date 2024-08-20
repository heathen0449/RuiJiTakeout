package models

import "takeout/utils"

type Category struct {
	Id         int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	Type       int              `json:"type" gorm:"type:int"` //1-菜品分类 2-套餐分类
	Name       string           `gorm:"size:32 uniqueIndex" json:"name"`
	Sort       int              `gorm:"default:0; type:int" json:"sort"`   //用于分类数据的排序
	Status     int              `json:"status" gorm:"default:1; type:int"` // 1启用 0禁用
	CreateTime utils.CustomTime `json:"createTime" gorm:"autoCreateTime, type:datetime"`
	UpdateTime utils.CustomTime `json:"updateTime" gorm:"autoUpdateTime, type:datetime"`
	CreateUser int64            `json:"createUser"`
	UpdateUser int64            `json:"updateUser"`
}
