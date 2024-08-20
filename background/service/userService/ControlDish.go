package userService

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"takeout/db"
	"takeout/models"
)

func GetDishByCategory(e echo.Context) error {
	categoryId, _ := strconv.ParseInt(e.QueryParam("categoryId"), 10, 64)
	var dishList []models.Dish
	err := db.DB.Model(&models.Dish{}).Where("category_id = ?", categoryId).Find(&dishList).Error
	if err != nil {
		answer := new(models.Result[interface{}])
		answer.Error("获取菜品失败")
		return e.JSON(http.StatusInternalServerError, answer)
	}
	for key, dish := range dishList {
		var dishFavorites []models.DishFlavor
		err := db.DB.Model(&models.DishFlavor{}).Where("dish_id = ?", dish.Id).Find(&dishFavorites).Error
		if err != nil {
			answer := new(models.Result[interface{}])
			answer.Error("获取菜品口味失败")
			return e.JSON(http.StatusInternalServerError, answer)
		}
		dishList[key].Flavors = dishFavorites
	}
	answer := new(models.Result[[]models.Dish])
	answer.SuccessWithObject(dishList)
	return e.JSON(http.StatusOK, answer)
}
