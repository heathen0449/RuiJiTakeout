package userService

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"takeout/db"
	"takeout/models"
)

func GetShopStatus(e echo.Context) error {
	result, err := db.Rdb.Get(db.Context, "SHOP_STATUS").Result()
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("获取失败")
		return e.JSON(400, result)
	}
	answer := new(models.Result[int])
	answer.Msg = "获取成功"
	data, _ := strconv.Atoi(result)
	answer.SuccessWithObject(data)
	return e.JSON(200, answer)
}
