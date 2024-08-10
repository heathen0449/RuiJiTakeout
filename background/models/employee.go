package models

import "time"

type Employee struct {
	Id         int64     `gorm:"primaryKey;autoIncrement"`
	UserName   string    `gorm:"unique; not null; uniqueIndex; size:32; column:username"`
	Name       string    `gorm:"size:32"`
	Password   string    `gorm:"size:64"`
	Phone      string    `gorm:"size:11"`
	Sex        string    `gorm:"size:2"`
	IdNumber   string    `gorm:"size:18"`
	Status     int8      `gorm:"default:1"`
	CreateTime time.Time `gorm:"autoCreateTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime"`
	CreateUser int64
	UpdateUser int64
}
