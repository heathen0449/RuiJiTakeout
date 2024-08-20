package vo

import "takeout/utils"

type SetmealPageVO struct {
	Id           int64            `json:"id"`
	CategoryId   int64            `json:"categoryId"`
	Name         string           `json:"name"`
	Price        float64          `json:"price"`
	Status       int32            `json:"status"`
	Description  string           `json:"description"`
	Image        string           `json:"image"`
	UpdateTime   utils.CustomTime `json:"updateTime"`
	CategoryName string           `json:"categoryName"`
}
