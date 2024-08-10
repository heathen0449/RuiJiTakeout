package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"takeout/db"
	"takeout/models"
	"takeout/models/dto"
	"takeout/models/vo"
	"time"
)

func LoginOut(e echo.Context) error {
	m := new(models.Result[interface{}])
	m.Success()
	return e.JSON(http.StatusOK, m)
}

func Login(e echo.Context) error {
	var employeeDTO dto.EmployeeLoginDTO
	if err := e.Bind(&employeeDTO); err != nil {
		result := new(models.Result[interface{}])
		result.Error("参数绑定失败")
		return e.JSON(http.StatusBadRequest, result)
	}
	var employees models.Employee
	db.DB.Where("username = ?", employeeDTO.UserName).First(&employees)
	if employees == (models.Employee{}) {
		result := new(models.Result[interface{}])
		result.Error("用户不存在")
		return e.JSON(http.StatusBadRequest, result)
	}
	hash := md5.New()
	hash.Write([]byte(employeeDTO.Password))
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	if employees.Password != hashString {
		result := new(models.Result[interface{}])
		result.Error("密码错误")
		return e.JSON(http.StatusBadRequest, result)
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["empId"] = employees.Id
	claims["role"] = employees.Name
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	t, err := token.SignedString([]byte("itcast"))

	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("token生成失败")
		return e.JSON(http.StatusBadRequest, result)
	}
	result := new(models.Result[vo.EmployeeLoginVO])

	result.SuccessWithObject(vo.EmployeeLoginVO{
		Id:       int(employees.Id),
		Name:     employees.Name,
		UserName: employees.UserName,
		Token:    t,
	})
	result.Msg = "登录成功"
	return e.JSON(http.StatusOK, result)
}
