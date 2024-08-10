package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"takeout/db"
	"takeout/models"
	"takeout/models/dto"
)

func CreateUser(e echo.Context) error {
	var employeeDTO dto.EmployeeDTO
	if err := e.Bind(&employeeDTO); err != nil {
		result := new(models.Result[interface{}])
		result.Error("参数绑定失败")
		return e.JSON(http.StatusBadRequest, result)
	}

	employee := new(models.Employee)
	// 将 employeeDTO 转换为 JSON 字符串
	employeeJSON, err := json.Marshal(employeeDTO)
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("JSON 编码失败")
		return e.JSON(http.StatusInternalServerError, result)
	}
	if err := json.Unmarshal(employeeJSON, employee); err != nil {
		result := new(models.Result[interface{}])
		result.Error("JSON 解码失败")
		return e.JSON(http.StatusInternalServerError, result)
	}
	// 从jwt获取
	empId := e.Get("empId")
	employee.CreateUser = empId.(int64)
	employee.UpdateUser = empId.(int64)
	// 密码加密
	employee.Password = "123456"
	hash := md5.New()
	hash.Write([]byte(employee.Password))
	hashInBytes := hash.Sum(nil)
	employee.Password = hex.EncodeToString(hashInBytes)
	// 保存到数据库
	result := db.DB.Create(&employee)
	if result.Error != nil {
		errorMess := result.Error
		result := new(models.Result[interface{}])
		findError := fmt.Sprintf("%v", errorMess)
		if strings.Contains(findError, "Duplicate entry") {
			result.Error("创建用户失败 用户名" + employeeDTO.UserName + "已存在")
		} else {
			result.Error("创建用户失败 ---" + findError)
		}
		return e.JSON(http.StatusInternalServerError, result)
	}
	response := new(models.Result[interface{}])
	response.Success()
	response.Msg = "创建用户成功"
	return e.JSON(http.StatusOK, response)
}

func GetEmployeePage(e echo.Context) error {
	name := e.QueryParam("name")
	page, _ := strconv.Atoi(e.QueryParam("page"))

	pageSize, _ := strconv.Atoi(e.QueryParam("pageSize"))
	var employees []models.Employee
	// 这里是一个指针 *PageResult
	pageResult := new(models.PageResult)
	var result *gorm.DB
	if name != "" {
		result = db.DB.Where("name like ?", "%"+name+"%").Offset((page - 1) * pageSize).Limit(pageSize).
			Order("create_time desc").Find(&employees)
	} else {
		result = db.DB.Offset((page - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&employees)
	}
	// []Interface{} 不能转换成 []Employee
	pageResult.Records = models.ToInterfaceSlice(employees)
	pageResult.Total = result.RowsAffected
	response := new(models.Result[models.PageResult])
	response.SuccessWithObject(*pageResult)
	return e.JSON(http.StatusOK, response)

}
