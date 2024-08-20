package models

import (
	"takeout/utils"
)

type User struct {
	Id         int64            `gorm:"primaryKey;autoIncrement;"`
	OpenId     string           `gorm:"type:varchar(45); column:openid"`
	Name       string           `gorm:"type:varchar(32)"`
	Phone      string           `gorm:"type:varchar(11)"`
	sex        string           `gorm:"type:varchar(2)"`
	IdNumber   string           `gorm:"type:varchar(18)"`
	Avatar     string           `gorm:"type:varchar(500)"`
	CreateTime utils.CustomTime `gorm:"autoCreateTime, dataType:datetime"`
}
