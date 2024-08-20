package service

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
	"takeout/db"
	"takeout/models"
	"takeout/models/dto"
	"takeout/models/vo"
	"takeout/utils"
	"time"
)

func ChangeSetmeal(e echo.Context) error {
	var setMealDTO dto.SetMealDTO
	//var setMeal models.Setmeal
	if err := e.Bind(&setMealDTO); err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("参数绑定失败: %v", err)
		result.Error("参数绑定失败")
		return e.JSON(400, result)
	}
	setMealDTO.UpdateUser = e.Get("empId").(int64)
	setMealDTO.UpdateTime = utils.CustomTime{Time: time.Now()}
	var price float64
	switch setMealDTO.Price.(type) {
	case float64:
		price = setMealDTO.Price.(float64)
	case string:
		price, _ = strconv.ParseFloat(setMealDTO.Price.(string), 64)
	}

	err := db.DB.Model(&models.Setmeal{}).Where("id = ?", setMealDTO.Id).Updates(
		map[string]interface{}{
			"CategoryId":  setMealDTO.CategoryId,
			"Name":        setMealDTO.Name,
			"Price":       price,
			"Status":      setMealDTO.Status,
			"Description": setMealDTO.Description,
			"Image":       setMealDTO.Image,
			"UpdateUser":  setMealDTO.UpdateUser,
			"UpdateTime":  setMealDTO.UpdateTime}).Error
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("操作失败")
		e.Logger().Errorf("Failed to update setmeal: %v", err)
		return e.JSON(500, result)
	}
	// 1. 查询现有的 dishes
	var setMealDishes []models.SetmealDish
	err = db.DB.Where("setmeal_id = ?", setMealDTO.Id).Find(&setMealDishes).Error
	if err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("Failed to update setmeal: %v", err)
		result.Error("操作失败")
		return e.JSON(500, result)
	}
	nowIds := make(map[string]interface{})
	for _, dish := range setMealDishes {
		nowIds[dish.Name] = nil
	}

	for _, dish := range setMealDTO.SetmealDishes {
		// 这里处理的是存在的 dish
		if _, ok := nowIds[dish.Name]; ok {
			err = db.DB.Model(&dish).Where("name", dish.Name).Updates(map[string]interface{}{
				"price":  dish.Price,
				"copies": dish.Copies,
			}).Error
			if err != nil {
				result := new(models.Result[interface{}])
				e.Logger().Errorf("Failed to update setmeal: %v", err)
				result.Error("操作失败")
				return e.JSON(500, result)
			}
			delete(nowIds, dish.Name)
		} else {
			// 这里处理的是不存在的 dish
			dish.SetmealId = setMealDTO.Id
			err = db.DB.Create(&dish).Error
			if err != nil {
				result := new(models.Result[interface{}])
				e.Logger().Errorf("Failed to update setmeal: %v", err)
				result.Error("操作失败")
				return e.JSON(500, result)
			}
		}
	}
	for _, dish := range setMealDishes {
		if _, ok := nowIds[dish.Name]; ok {
			err = db.DB.Delete(&dish).Error
			if err != nil {
				result := new(models.Result[interface{}])
				e.Logger().Errorf("Failed to update setmeal: %v", err)
				result.Error("操作失败")
				return e.JSON(500, result)
			}
		}
	}
	result := new(models.Result[interface{}])
	result.Success()
	return e.JSON(200, result)
}

func GetSetmealPage(e echo.Context) error {
	page, _ := strconv.Atoi(e.QueryParam("page"))
	pageSize, _ := strconv.Atoi(e.QueryParam("pageSize"))
	status, _ := strconv.ParseInt(e.QueryParam("status"), 10, 32)
	categoryId := e.QueryParam("categoryId")
	name := e.QueryParam("name")
	condition := db.DB
	if name != "" {
		condition = condition.Where("name like ?", "%"+name+"%")
	}
	if categoryId != "" {
		condition = condition.Where("category_id = ?")
	}
	if status != 0 {
		condition = condition.Where("status = ?", status)
	}
	pageAnswer := new(models.PageResult)
	var total int64
	err := condition.Model(&models.Setmeal{}).Count(&total).Error
	if err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("Failed to query setmeal: %v", err)
		result.Error("操作失败")
		return e.JSON(500, result)
	}
	pageAnswer.Total = total
	setMeals := make([]models.Setmeal, 0)
	err = condition.Offset((page - 1) * pageSize).Limit(pageSize).Order("").Find(&setMeals).Error
	if err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("Failed to query setmeal: %v", err)
		result.Error("操作失败")
		return e.JSON(500, result)
	}
	var FinalAnswer []vo.SetmealPageVO
	for _, v := range setMeals {
		var setmealPageVO vo.SetmealPageVO
		setmealPageVO.Name = v.Name
		setmealPageVO.Id = v.Id
		setmealPageVO.CategoryId = v.CategoryId
		setmealPageVO.Description = v.Description
		setmealPageVO.Image = v.Image
		setmealPageVO.Price = v.Price
		setmealPageVO.Status = v.Status
		setmealPageVO.UpdateTime = v.UpdateTime
		var name string
		err = db.DB.Model(&models.Category{}).Where("id = ?", v.CategoryId).Select("Name").Scan(&name).Error
		if err != nil {
			result := new(models.Result[interface{}])
			e.Logger().Errorf("Failed to query setmeal: %v", err)
			result.Error("操作失败")
			return e.JSON(500, result)
		}
		setmealPageVO.CategoryName = name
		FinalAnswer = append(FinalAnswer, setmealPageVO)
	}
	pageAnswer.Records = models.ToInterfaceSlice(FinalAnswer)
	answer := new(models.Result[models.PageResult])
	answer.SuccessWithObject(*pageAnswer)
	return e.JSON(200, answer)
}

func StartOrStopSetmeal(e echo.Context) error {
	id, _ := strconv.ParseInt(e.QueryParam("id"), 10, 64)
	status, _ := strconv.ParseInt(e.Param("status"), 10, 32)
	err := db.DB.Model(&models.Setmeal{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"update_user": e.Get("empId").(int64),
		"update_time": utils.CustomTime{Time: time.Now()},
	}).Error
	if err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("Failed to update setmeal: %v", err)
		result.Error("操作失败")
		return e.JSON(500, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	return e.JSON(200, result)
}

func DeleteSetmeal(e echo.Context) error {
	idsList := strings.Split(e.Param("ids"), ",")
	ids := make([]int64, 0)
	for _, v := range idsList {
		id, _ := strconv.ParseInt(v, 10, 64)
		ids = append(ids, id)
	}
	err := db.DB.Where("id in (?)", ids).Delete(&models.Setmeal{}).Error
	if err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("Failed to delete setmeal: %v", err)
		result.Error("操作失败")
		return e.JSON(500, result)
	}
	err = db.DB.Where("setmeal_id in (?)", ids).Delete(&models.SetmealDish{}).Error
	if err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("Failed to delete setmeal: %v", err)
		result.Error("操作失败")
		return e.JSON(500, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	return e.JSON(200, result)
}

func CreateSetmeal(e echo.Context) error {
	var setmeal models.Setmeal
	var setmealDTO dto.SetMealDTO
	if err := e.Bind(&setmealDTO); err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("Failed to bind setmeal: %v", err)
		result.Error("参数绑定失败")
		return e.JSON(400, result)
	}
	setmeal.CategoryId = setmealDTO.CategoryId
	setmeal.Name = setmealDTO.Name
	switch setmealDTO.Price.(type) {
	case float64:
		setmeal.Price = setmealDTO.Price.(float64)
	case string:
		setmeal.Price, _ = strconv.ParseFloat(setmealDTO.Price.(string), 64)
	}
	setmeal.Status = setmealDTO.Status
	setmeal.Description = setmealDTO.Description
	setmeal.Image = setmealDTO.Image
	setmeal.SetmealDishes = setmealDTO.SetmealDishes
	setmeal.CreateUser = e.Get("empId").(int64)
	setmeal.UpdateUser = e.Get("empId").(int64)
	setmeal.CreateTime = utils.CustomTime{Time: time.Now()}
	setmeal.UpdateTime = utils.CustomTime{Time: time.Now()}
	err := db.DB.Create(&setmeal).Error
	if err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("Failed to create setmeal: %v", err)
		result.Error("操作失败")
		return e.JSON(500, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	return e.JSON(200, result)
}

func GetSetmealById(e echo.Context) error {
	id := e.Param("id")
	var setmeal models.Setmeal
	err := db.DB.Where("id = ?", id).First(&setmeal).Error
	if err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("Failed to query setmeal: %v", err)
		result.Error("操作失败")
		return e.JSON(500, result)
	}
	var setmealDishes []models.SetmealDish
	err = db.DB.Where("setmeal_id = ?", id).Find(&setmealDishes).Error
	if err != nil {
		result := new(models.Result[interface{}])
		e.Logger().Errorf("Failed to query setmeal: %v", err)
		result.Error("操作失败")
		return e.JSON(500, result)
	}
	setmeal.SetmealDishes = setmealDishes
	result := new(models.Result[interface{}])
	result.SuccessWithObject(setmeal)
	return e.JSON(200, result)
}
