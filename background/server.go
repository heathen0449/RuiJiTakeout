package main

import "C"
import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"os"
	"takeout/db"
	"takeout/models"
	"takeout/service"
	"takeout/service/userService"
	"takeout/utils"
)

func main() {
	e := echo.New()
	// 配置 Echo 的日志记录器
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetOutput(os.Stdout)
	// 这是一个中间件，用于记录请求日志
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
		Output: os.Stdout, // Ensure logs are output to the console
	}))
	db.Context = context.Background()
	e.Use(middleware.Recover())
	db.Init()
	defer db.Close()
	db.RedisInit()
	defer db.RedisClose()

	e.POST("/admin/employee/login", service.Login)
	e.POST("/user/user/login", userService.Login)
	e.GET("/user/shop/status", userService.GetShopStatus)

	adminGroup := e.Group("/admin")
	userGroup := e.Group("/user")
	err := utils.BucketInit()
	if err != nil {
		fmt.Printf("oss client error: %v\n", err)
	}

	adminGroup.Use(echojwt.WithConfig(echojwt.Config{
		TokenLookup: "header:token",
		SigningKey:  []byte("itcast"),
		ErrorHandler: func(c echo.Context, err error) error {
			result := new(models.Result[interface{}])
			result.Error(err.Error())
			return c.JSON(401, result)
		},
		SuccessHandler: func(c echo.Context) {
			// 名称特定是user 通过 ContextKey string 设置，默认为user
			token, ok := c.Get("user").(*jwt.Token)
			if !ok || token == nil {
				c.Error(errors.New("JWT token missing or invalid"))
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.Error(errors.New("failed to cast claims as jwt.MapClaims"))
				return
			}
			// 将 empId 从 float64 转换为 int64
			// 在 Go 中，JWT 中的数值类型通常会被解析为 float64，尤其是在使用 jwt.MapClaims 时。
			//如果你直接将这些数值转换为 int64 而没有进行类型检查或转换，就会导致类型不匹配的错误。
			if empIdFloat, ok := claims["empId"].(float64); ok {
				empId := int64(empIdFloat)
				c.Set("empId", empId)
			} else {
				c.Error(errors.New("empId not found or is of incorrect type"))
				return
			}

			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			} else {
				c.Error(errors.New("role not found or is of incorrect type"))
				return
			}
		},
	}))

	userGroup.Use(echojwt.WithConfig(echojwt.Config{
		TokenLookup: "header:authentication",
		SigningKey:  []byte("itheima"),
	}))
	adminGroup.POST("/employee/logout", service.LoginOut)
	// 用户管理操作
	adminGroup.POST("/employee", service.CreateUser)
	//单页查询 用户查询
	adminGroup.GET("/employee/page", service.GetEmployeePage)
	// 修改用户启用，禁用
	adminGroup.POST("/employee/status/:status", service.StartOrStop)
	// 修改用户信息
	adminGroup.PUT("/employee", service.UpdateEmployee)
	// 根据用户id查询用户信息
	adminGroup.GET("/employee/:id", service.GetEmployeeById)

	adminGroup.POST("/category", service.AddCategory)
	adminGroup.GET("/category/page", service.GetCategoryPage)
	adminGroup.PUT("/category", service.UpdateCategory)
	adminGroup.POST("/category/status/:status", service.StartOrStopCategory)
	adminGroup.DELETE("/category", service.DeleteCategoryById)
	adminGroup.GET("/category/list", service.GetCategoryByType)

	adminGroup.PUT("/dish", service.ChangeDish)
	adminGroup.DELETE("/dish", service.DeleteDishes)
	adminGroup.POST("/dish", service.AddDish)
	adminGroup.GET("/dish/:id", service.GetDishById)
	adminGroup.GET("/dish/list", service.GetDishByCategoryId)
	adminGroup.GET("/dish/page", service.GetDishPage)
	adminGroup.POST("/dish/status/:status", service.StartOrStopDish)
	adminGroup.GET("/shop/status", service.GetService)
	adminGroup.PUT("/shop/:status", service.SetService)

	adminGroup.PUT("/setmeal", service.ChangeSetmeal)
	adminGroup.GET("/setmeal/:id", service.GetSetmealById)
	adminGroup.GET("/setmeal/page", service.GetSetmealPage)
	adminGroup.POST("/setmeal/status/:status", service.StartOrStopSetmeal)
	adminGroup.DELETE("/setmeal", service.DeleteSetmeal)
	adminGroup.POST("/setmeal", service.CreateSetmeal)

	adminGroup.POST("/common/upload", func(e echo.Context) error {
		file, err := e.FormFile("file")
		if err != nil {
			result := new(models.Result[interface{}])
			result.Error(err.Error())
			result.Msg = "参数绑定失败"
			e.Logger().Errorf("参数绑定失败: %v", err)
			return e.JSON(400, result)
		}
		src, err := file.Open()
		if err != nil {
			result := new(models.Result[interface{}])
			result.Error(err.Error())
			result.Msg = "文件打开失败"
			e.Logger().Errorf("文件打开失败: %v", err)
			return e.JSON(400, result)
		}
		defer src.Close()
		url, err := utils.UploadImage(src, file.Filename)
		if err != nil {
			result := new(models.Result[interface{}])
			result.Error(err.Error())
			result.Msg = "文件上传失败"
			e.Logger().Errorf("文件上传失败: %v", err)
			return e.JSON(400, result)
		}
		result := new(models.Result[string])
		result.SuccessWithObject(url)
		return e.JSON(200, result)
	})

	userGroup.GET("/category/list", userService.GetCategoryByType)
	userGroup.GET("/setmeal/list", userService.SetmealByCategoryId)
	userGroup.GET("/setmeal/dish/:id", userService.GetMealBySetmealId)
	userGroup.GET("/dish/list", userService.GetDishByCategory)
	e.Start(":8080")
}
