package service

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"takeout/db"
	"takeout/models"
	"takeout/models/dto"
	"takeout/utils"
	"time"
)

func AddCategory(e echo.Context) error {
	// 新增分类
	var categoryDTO dto.CategoryDTO
	if err := e.Bind(&categoryDTO); err != nil {
		result := new(models.Result[interface{}])
		result.Error("参数绑定失败")
		return e.JSON(400, result)
	}

	category := new(models.Category)
	category.Name = categoryDTO.Name
	category.Type = categoryDTO.Type
	category.Sort = categoryDTO.Sort
	category.CreateUser = e.Get("empId").(int64)
	category.UpdateUser = e.Get("empId").(int64)

	answer := db.DB.Create(&category)
	if answer.Error != nil {
		result := new(models.Result[interface{}])
		result.Error(answer.Error.Error())
		return e.JSON(http.StatusInternalServerError, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	result.Msg = "新增分类成功"
	return e.JSON(http.StatusOK, result)
}

func GetCategoryPage(e echo.Context) error {
	name := e.QueryParam("name")
	page, _ := strconv.Atoi(e.QueryParam("page"))
	pageSize, _ := strconv.Atoi(e.QueryParam("pageSize"))
	nowType, _ := strconv.Atoi(e.QueryParam("type"))

	var categories []models.Category
	var pageResult models.PageResult
	condition := db.DB
	if nowType != 0 {
		condition = condition.Where(&models.Category{Type: nowType})
	}
	if name != "" {
		condition = condition.Where("name like ?", "%"+name+"%")
	}
	pageResult.Total = condition.Find(&categories).RowsAffected
	findResult := condition.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("sort asc , create_time desc").Find(&categories)
	if findResult.Error != nil {
		result := new(models.Result[interface{}])
		result.Error(findResult.Error.Error())
		return e.JSON(http.StatusInternalServerError, result)
	}
	pageResult.Records = models.ToInterfaceSlice(categories)
	result := new(models.Result[models.PageResult])
	result.SuccessWithObject(pageResult)

	return e.JSON(http.StatusOK, result)
}

func UpdateCategory(e echo.Context) error {
	var categoryDTO dto.CategoryDTO
	if err := e.Bind(&categoryDTO); err != nil {
		result := new(models.Result[interface{}])
		result.Error("参数绑定失败")
		return e.JSON(400, result)
	}

	category := new(models.Category)
	category.Id = categoryDTO.Id
	category.Name = categoryDTO.Name
	category.Type = categoryDTO.Type
	category.Sort = categoryDTO.Sort
	category.UpdateUser = e.Get("empId").(int64)

	answer := db.DB.Model(&models.Category{Id: category.Id}).Updates(map[string]interface{}{
		"id":          category.Id,
		"name":        category.Name,
		"type":        category.Type,
		"sort":        category.Sort,
		"update_user": category.UpdateUser,
		"update_time": utils.CustomTime{Time: time.Now()},
	})

	if answer.Error != nil {
		result := new(models.Result[interface{}])
		result.Error(answer.Error.Error())
		return e.JSON(http.StatusInternalServerError, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	result.Msg = "修改分类成功"
	return e.JSON(http.StatusOK, result)
}

func StartOrStopCategory(e echo.Context) error {
	id, _ := strconv.Atoi(e.QueryParam("id"))
	status, _ := strconv.Atoi(e.Param("status"))
	category := new(models.Category)
	category.Id = int64(id)
	e.Logger().Infof("id: %d, status: %d", id, category.Status)
	//GORM的零值处理: GORM 默认情况下会忽略零值（例如 0、""、false）的字段。这意味着，如果 status 被设置为零值（0），
	//GORM 会认为这个字段未被修改，从而不会将其包含在生成的更新语句中。
	//为了避免这种情况，你可以使用 Select 方法指定要更新的字段，即使它们是零值。
	//例如，下面的代码将更新所有字段，即使它们是零值：
	//db.Model(&user).select("name", "age").Updates(models.Category{
	//    Id:          category.Id,
	//    Status:      int8(status),
	//    UpdateTime:  utils.CustomTime{Time: time.Now()},
	//    UpdateUser:  category.UpdateUser,
	//}))
	//或者
	//db.Model(&user).Updates(map[string]interface{}{"name": "jinni", "age": 18})
	answer := db.DB.Debug().Model(&category).Select("Status", "UpdateTime", "UpdateUser").Updates(models.Category{Id: category.Id,
		Status:     status,
		UpdateTime: utils.CustomTime{Time: time.Now()},
		UpdateUser: category.UpdateUser})
	if answer.Error != nil {
		result := new(models.Result[interface{}])
		result.Error(answer.Error.Error())
		return e.JSON(http.StatusInternalServerError, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	result.Msg = "修改分类状态成功"
	return e.JSON(http.StatusOK, result)
}

func DeleteCategoryById(e echo.Context) error {
	id, _ := strconv.Atoi(e.QueryParam("id"))

	category := new(models.Category)
	category.Id = int64(id)

	answer := db.DB.Delete(&category)
	if answer.Error != nil {
		result := new(models.Result[interface{}])
		result.Error(answer.Error.Error())
		return e.JSON(http.StatusInternalServerError, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	result.Msg = "删除分类成功"
	return e.JSON(http.StatusOK, result)
}

func GetCategoryByType(e echo.Context) error {
	nowType, _ := strconv.Atoi(e.QueryParam("type"))

	var categories []models.Category
	condition := db.DB.Where(&models.Category{Status: 1})
	if nowType != 0 {
		condition = condition.Where(&models.Category{Type: nowType})
	}
	answer := condition.Order("sort asc,create_time desc").Find(&categories)
	if answer.Error != nil {
		result := new(models.Result[interface{}])
		result.Error(answer.Error.Error())
		return e.JSON(http.StatusInternalServerError, result)
	}
	result := new(models.Result[interface{}])
	result.SuccessWithObject(models.ToInterfaceSlice(categories))
	return e.JSON(http.StatusOK, result)
}
