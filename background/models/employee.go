package models

import (
	"takeout/utils"
)

type Employee struct {
	Id         int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName   string           `gorm:"unique; not null; uniqueIndex; size:32; column:username" json:"username"`
	Name       string           `gorm:"size:32" json:"name"`
	Password   string           `gorm:"size:64" json:"password"`
	Phone      string           `gorm:"size:11" json:"phone"`
	Sex        string           `gorm:"size:2" json:"sex"`
	IdNumber   string           `gorm:"size:18" json:"idNumber"`
	Status     int8             `gorm:"default:1" json:"status"`                         // 1正常 0锁定
	CreateTime utils.CustomTime `json:"createTime" gorm:"autoCreateTime; type:datetime"` // 因为时间格式的问题，所以需要自定义时间类型
	UpdateTime utils.CustomTime `json:"updateTime" gorm:"autoUpdateTime; type:datetime"` // 同时还要标记数据库中也是时间类型
	CreateUser int64            `json:"createUser"`
	UpdateUser int64            `json:"updateUser"`
}
