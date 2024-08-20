package userService

import (
	"github.com/labstack/echo/v4"
	"takeout/db"
	"takeout/models"
)

func GetCategoryByType(e echo.Context) error {
	typeName := e.QueryParam("type")
	condition := db.DB.Model(&models.Category{})
	if typeName != "" {
		condition = condition.Where("type = ?", typeName)
	}
	var categoryList []models.Category
	err := condition.Find(&categoryList).Error

	if err != nil {
		result := new(models.Result[interface{}])
		result.Error(err.Error())
		return e.JSON(500, result)
	}

	result := new(models.Result[[]models.Category])
	result.SuccessWithObject(categoryList)
	return e.JSON(200, result)
}
