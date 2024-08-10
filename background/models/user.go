package models

import "time"

type User struct {
	Id         int64     `gorm:"primaryKey;autoIncrement"`
	OpenId     string    `gorm:"size:45; column: openid"`
	Name       string    `gorm:"size:32"`
	Phone      string    `gorm:"size:11"`
	sex        string    `gorm:"size:2"`
	IdNumber   string    `gorm:"size:18"`
	Avatar     string    `gorm:"size:500"`
	CreateTime time.Time `gorm:"autoCreateTime"`
}
