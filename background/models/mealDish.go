package models

import "takeout/utils"

type Setmeal struct {
	Id            int64            `json:"id" gorm:"autoIncrement;primaryKey type:bigint; not null"`
	CategoryId    int64            `json:"categoryId" gorm:"type:bigint; not null"`
	Name          string           `json:"name" gorm:"type:varchar(32); uniqueIndex; not null"`
	Price         float64          `json:"price" gorm:"type:decimal(10,2)"`
	Status        int32            `json:"status" gorm:"type:int; default:1"`
	Description   string           `json:"description" gorm:"type:varchar(255)"`
	Image         string           `json:"image" gorm:"type:varchar(255)"`
	CreateTime    utils.CustomTime `json:"createTime" gorm:"type:datetime; autoCreateTime"`
	UpdateTime    utils.CustomTime `json:"updateTime" gorm:"type:datetime; autoUpdateTime"`
	CreateUser    int64            `json:"createUser" gorm:"type:bigint"`
	UpdateUser    int64            `json:"updateUser" gorm:"type:bigint"`
	SetmealDishes []SetmealDish    `json:"setmealDishes" gorm:"foreignKey:SetmealId; references:Id"`
}

type SetmealDish struct {
	Id        int64   `json:"id" gorm:"autoIncrement;primaryKey type:bigint; not null"`
	SetmealId int64   `json:"setmealId" gorm:"type:bigint; not null"`
	DishId    int64   `json:"dishId" gorm:"type:bigint; not null"`
	Name      string  `json:"name" gorm:"type:varchar(32)"`
	Price     float64 `json:"price" gorm:"type:decimal(10,2)"`
	Copies    int32   `json:"copies" gorm:"type:int"`
}
