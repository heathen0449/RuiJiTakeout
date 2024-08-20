package userService

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"takeout/db"
	"takeout/models"
	"takeout/models/dto"
	"takeout/models/vo"
)

func SetmealByCategoryId(e echo.Context) error {
	categoryId, _ := strconv.ParseInt(e.QueryParam("categoryId"), 10, 64)
	var setmeals []dto.SetmealUser
	if err := db.DB.Where("category_id = ?", categoryId).Find(&setmeals).Error; err != nil {
		result := new(models.Result[interface{}])
		result.Error(err.Error())
		return e.JSON(500, result)
	}
	result := new(models.Result[[]dto.SetmealUser])
	result.SuccessWithObject(setmeals)
	return e.JSON(200, result)
}

func GetMealBySetmealId(e echo.Context) error {
	setmealId, _ := strconv.ParseInt(e.Param("id"), 10, 64)
	var setmealDishes []models.SetmealDish
	if err := db.DB.Where("setmeal_id = ?", setmealId).Find(&setmealDishes).Error; err != nil {
		result := new(models.Result[interface{}])
		result.Error(err.Error())
		return e.JSON(500, result)
	}
	var dishes []vo.SetMealDishesUserVO
	for _, setmealDish := range setmealDishes {
		dishes = append(dishes, vo.SetMealDishesUserVO{
			Copies: setmealDish.Copies,
			Name:   setmealDish.Name,
		})
		var dish models.Dish
		if err := db.DB.Where("id = ?", setmealDish.DishId).Select("Description", "Image").
			Take(&dish).Error; err != nil {
			result := new(models.Result[interface{}])
			result.Error(err.Error())
			return e.JSON(500, result)
		}
		dishes[len(dishes)-1].Description = dish.Description
		dishes[len(dishes)-1].Image = dish.Image
	}
	result := new(models.Result[[]vo.SetMealDishesUserVO])
	result.SuccessWithObject(dishes)
	return e.JSON(200, result)
}
