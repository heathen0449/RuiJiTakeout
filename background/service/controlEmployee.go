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
	"takeout/utils"
	"time"
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
	//employee.CreateTime = utils.CustomTime{Time: time.Now()}
	//employee.UpdateTime = utils.CustomTime{Time: time.Now()}
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
		pageResult.Total = result.RowsAffected
	} else {
		pageResult.Total = db.DB.Find(&employees).RowsAffected
		result = db.DB.Offset((page - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&employees)
	}
	// []Interface{} 不能转换成 []Employee
	pageResult.Records = models.ToInterfaceSlice(employees)

	response := new(models.Result[models.PageResult])
	response.SuccessWithObject(*pageResult)
	return e.JSON(http.StatusOK, response)

}

func StartOrStop(e echo.Context) error {
	id, err := strconv.Atoi(e.QueryParam("id"))
	status, err := strconv.Atoi(e.Param("status"))
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("参数绑定失败")
		return e.JSON(http.StatusBadRequest, result)
	}
	// 这里操作的是echo本身的日志
	e.Logger().Infof("id: %d, status: %d", id, status)
	answer := db.DB.Model(&models.Employee{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
		"update_user": e.Get("empId"),
	})
	if answer.Error != nil {
		result := new(models.Result[interface{}])
		result.Error("操作失败")
		return e.JSON(http.StatusInternalServerError, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	return e.JSON(http.StatusOK, result)
}

func GetEmployeeById(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("参数绑定失败")
		return e.JSON(http.StatusBadRequest, result)
	}
	var employee models.Employee
	answer := db.DB.First(&employee, id)
	if answer.Error != nil {
		result := new(models.Result[interface{}])
		result.Error("查询失败")
		return e.JSON(http.StatusInternalServerError, result)
	}
	result := new(models.Result[models.Employee])
	result.SuccessWithObject(employee)
	return e.JSON(http.StatusOK, result)
}
func UpdateEmployee(e echo.Context) error {
	var employeeDTO dto.EmployeeDTO
	if err := e.Bind(&employeeDTO); err != nil {
		result := new(models.Result[interface{}])
		result.Error("参数绑定失败")
		return e.JSON(http.StatusBadRequest, result)
	}
	employee := new(models.Employee)
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
	employee.UpdateUser = e.Get("empId").(int64)
	employee.UpdateTime = utils.CustomTime{Time: time.Now()}
	e.Logger().Infof("employee: %v", employee.UpdateTime)
	answer := db.DB.Model(&employee).Updates(employee)
	//answer := db.DB.Model(&models.Employee{}).Where("id = ?", employee.Id).Updates(employee)

	if answer.Error != nil {
		e.Logger().Errorf("Failed to update employee: %v", answer.Error)
		result := new(models.Result[interface{}])
		result.Error("更新失败")
		return e.JSON(http.StatusInternalServerError, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	result.Msg = "更新成功"
	return e.JSON(http.StatusOK, result)
}
