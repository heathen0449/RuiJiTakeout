package vo

import (
	"takeout/utils"
)

type DishVOWithCategoryName struct {
	Id           int64            `json:"id"`
	Name         string           `json:"name"`
	CategoryId   int64            `json:"categoryId" `
	Price        float64          `json:"price" `
	Image        string           `json:"image"`
	Description  string           `json:"description"`
	Status       int32            `json:"status"`
	UpdateTime   utils.CustomTime `json:"updateTime"`
	CategoryName string           `json:"categoryName"`
}
