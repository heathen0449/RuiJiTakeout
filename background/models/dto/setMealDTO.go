package dto

import (
	"takeout/models"
	"takeout/utils"
)

type SetMealDTO struct {
	Id            int64                `json:"id"`
	CategoryId    int64                `json:"categoryId" `
	Name          string               `json:"name" `
	Price         utils.Price          `json:"price" `
	Status        int32                `json:"status" `
	Description   string               `json:"description" `
	Image         string               `json:"image"`
	CreateTime    utils.CustomTime     `json:"createTime" `
	UpdateTime    utils.CustomTime     `json:"updateTime"`
	CreateUser    int64                `json:"createUser"`
	UpdateUser    int64                `json:"updateUser"`
	SetmealDishes []models.SetmealDish `json:"setmealDishes"`
}
