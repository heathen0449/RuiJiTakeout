package dto

import "takeout/models"

type DishDTO struct {
	Id          int64               `json:"id"`
	Name        string              `json:"name"`
	CategoryId  int64               `json:"categoryId"`
	Price       string              `json:"price"`
	Image       string              `json:"image"`
	Description string              `json:"description"`
	Status      int32               `json:"status"`
	Flavors     []models.DishFlavor `json:"flavors"`
}
