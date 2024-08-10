package main

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"takeout/db"
	"takeout/models"
	"takeout/service"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	db.Init()
	e.POST("/admin/employee/login", service.Login)
	//e.POST("/admin/employee/logout", service.LoginOut)
	adminGroup := e.Group("/admin")

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
	adminGroup.POST("/employee/logout", service.LoginOut)
	// 用户管理操作
	adminGroup.POST("/employee", service.CreateUser)
	adminGroup.GET("/employee/page", service.GetEmployeePage)
	e.Start(":8080")
}
