package service

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"strconv"
	"takeout/db"
	"takeout/models"
)

func SetService(e echo.Context) error {
	status := e.Param("status")
	if err := db.Rdb.Set(db.Context, "SHOP_STATUS", status, 0).Err(); err != nil {
		result := new(models.Result[interface{}])
		result.Error("设置失败")
		return e.JSON(400, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	result.Msg = "设置成功"
	return e.JSON(200, result)
}

func GetService(e echo.Context) error {
	status, err := db.Rdb.Get(db.Context, "SHOP_STATUS").Result()
	if err != nil {
		result := new(models.Result[interface{}])
		fmt.Println("获取失败" + err.Error())
		result.Error("获取失败")
		return e.JSON(400, result)
	}
	result := new(models.Result[int])
	result.Msg = "获取成功"
	data, _ := strconv.Atoi(status)
	result.SuccessWithObject(data)
	return e.JSON(200, result)
}
