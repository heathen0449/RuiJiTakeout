package service

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"takeout/db"
	"takeout/models"
	"takeout/models/dto"
	"takeout/models/vo"
	"takeout/utils"
	"time"
)

func ChangeDish(e echo.Context) error {
	var dishDTO dto.DishDTO
	var dish models.Dish
	if err := e.Bind(&dishDTO); err != nil {
		result := new(models.Result[interface{}])
		result.Error("参数绑定失败")
		fmt.Println("参数绑定失败")
		return e.JSON(400, result)
	}
	dish.Id = dishDTO.Id
	dish.Name = dishDTO.Name
	dish.CategoryId = dishDTO.CategoryId
	dish.Price, _ = strconv.ParseFloat(dishDTO.Price, 64)
	dish.Image = dishDTO.Image
	dish.Description = dishDTO.Description
	dish.Status = dishDTO.Status
	dish.UpdateTime = utils.CustomTime{Time: time.Now()}
	dish.Flavors = dishDTO.Flavors
	dish.UpdateUser = e.Get("empId").(int64)
	// FullSaveAssociations: true: 这是一个 GORM 会话选项，用来确保在保存主记录时，关联的数据也会被保存或更新。
	//对于一对多关系，这个选项确保 DishFlavor 数据被正确更新。
	answer := db.DB.Model(&dish).Session(&gorm.Session{FullSaveAssociations: true}).Save(
		map[string]interface{}{
			"name":        dish.Name,
			"category_id": dish.CategoryId,
			"price":       dish.Price,
			"image":       dish.Image,
			"description": dish.Description,
			"status":      dish.Status,
			"update_user": dish.UpdateUser,
			"update_time": dish.UpdateTime,
		})
	// 1. 查找当前已有的 Flavors
	var existingFlavors []models.DishFlavor
	db.DB.Where("dish_id = ?", dish.Id).Find(&existingFlavors)

	// 将现有的 Flavors 转换为 map，便于查找
	existingFlavorMap := make(map[int64]models.DishFlavor)
	for _, flavor := range existingFlavors {
		existingFlavorMap[flavor.Id] = flavor
	}

	// 2. 处理新传入的 Flavors
	newFlavorIDs := make(map[int64]bool)
	for _, flavor := range dish.Flavors {
		if flavor.Id == 0 {
			// 新增的 Flavor
			flavor.DishId = dish.Id
			db.DB.Create(&flavor)
		} else {
			// 更新已有的 Flavor
			db.DB.Model(&flavor).Where("id = ?", flavor.Id).Updates(flavor)
			newFlavorIDs[flavor.Id] = true
		}
	}

	// 3. 删除不在新数据中的旧 Flavors
	for _, existingFlavor := range existingFlavors {
		if !newFlavorIDs[existingFlavor.Id] {
			db.DB.Delete(&existingFlavor)
		}
	}
	if answer.Error != nil {
		result := new(models.Result[interface{}])
		result.Error("更新失败")
		fmt.Println(answer.Error)
		return e.JSON(400, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	result.Msg = "更新成功"
	return e.JSON(200, result)
}

func DeleteDishes(e echo.Context) error {
	var idList []int
	for _, v := range strings.Split(e.QueryParam("ids"), ",") {
		ans, _ := strconv.Atoi(v)
		idList = append(idList, ans)
	}
	answer := db.DB.Delete(&models.Dish{}, idList)
	if answer.Error != nil {
		result := new(models.Result[interface{}])
		result.Error("删除失败")
		return e.JSON(400, result)
	}
	answer = db.DB.Where("dish_id in (?)", idList).Delete(&models.DishFlavor{})
	if answer.Error != nil {
		result := new(models.Result[interface{}])
		result.Error("删除失败")
		return e.JSON(400, result)
	}
	result := new(models.Result[interface{}])
	result.Msg = "删除成功"
	result.Success()
	return e.JSON(200, result)
}

func AddDish(e echo.Context) error {
	var dish models.Dish
	if err := e.Bind(&dish); err != nil {
		result := new(models.Result[interface{}])
		result.Error("参数绑定失败")
		return e.JSON(400, result)
	}
	dish.CreateUser = e.Get("empId").(int64)
	dish.UpdateUser = e.Get("empId").(int64)

	err := db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Select("Status").Updates(&dish).Error
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("添加失败")
		return e.JSON(400, result)
	}
	result := new(models.Result[interface{}])
	result.Msg = "添加成功"
	return e.JSON(200, result)
}

func GetDishById(e echo.Context) error {
	id, _ := strconv.Atoi(e.Param("id"))
	var dish models.Dish
	err := db.DB.Preload("Flavors").First(&dish, id).Error
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("查询失败")
		return e.JSON(400, result)
	}

	result := new(models.Result[models.Dish])
	result.SuccessWithObject(dish)
	return e.JSON(200, result)
}

func GetDishByCategoryId(e echo.Context) error {
	categoryId, _ := strconv.ParseInt(e.QueryParam("categoryId"), 10, 64)
	var dish []models.Dish
	err := db.DB.Where("category_id = ?", categoryId).Find(&dish).Error
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("查询失败")
		return e.JSON(400, result)
	}
	var name models.Category
	err = db.DB.Model(&models.Category{Id: categoryId}).Select("Name").First(&name).Error
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("查询失败")
		return e.JSON(400, result)
	}
	var answer []vo.DishVOWithCategoryName
	for _, v := range dish {
		answer = append(answer, vo.DishVOWithCategoryName{
			Id:           v.Id,
			Name:         v.Name,
			Price:        v.Price,
			Image:        v.Image,
			Description:  v.Description,
			Status:       v.Status,
			UpdateTime:   v.UpdateTime,
			CategoryId:   v.CategoryId,
			CategoryName: name.Name,
		})
	}
	result := new(models.Result[[]vo.DishVOWithCategoryName])
	result.SuccessWithObject(answer)
	return e.JSON(200, result)
}

func GetDishPage(e echo.Context) error {
	categoryId, _ := strconv.ParseInt(e.QueryParam("categoryId"), 10, 64)
	name := e.QueryParam("name")
	page, _ := strconv.Atoi(e.QueryParam("page"))
	pageSize, _ := strconv.Atoi(e.QueryParam("pageSize"))
	status, _ := strconv.Atoi(e.QueryParam("status"))
	var finalPage models.PageResult
	var dish []models.Dish
	condition := db.DB
	if categoryId != 0 {
		condition = condition.Where("category_id = ?", categoryId)
	}
	if name != "" {
		condition = condition.Where("name like ?", "%"+name+"%")
	}
	if e.QueryParam("status") != "" {
		condition = condition.Where("status = ?", status)
	}
	//condition := db.DB.Where("category_id = ? and name like ? and status = ?", categoryId, "%"+name+"%", status)
	ans := condition.Find(&dish)
	if ans.Error != nil {
		result := new(models.Result[interface{}])
		result.Error("查询失败")
		return e.JSON(400, result)
	}
	finalPage.Total = ans.RowsAffected
	err := condition.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("update_time desc, id").Find(&dish).Error
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("查询失败")
		return e.JSON(400, result)
	}
	var finalResult []vo.DishVOWithCategoryName
	for _, v := range dish {
		var now vo.DishVOWithCategoryName
		var name models.Category
		err = db.DB.Select("Name").First(&name, v.CategoryId).Error
		if err != nil {
			result := new(models.Result[interface{}])
			result.Error("查询失败")
			return e.JSON(400, result)
		}
		now.CategoryId = v.CategoryId
		now.CategoryName = name.Name
		now.Id = v.Id
		now.Name = v.Name
		now.Price = v.Price
		now.Image = v.Image
		now.Description = v.Description
		now.Status = v.Status
		now.UpdateTime = v.UpdateTime
		finalResult = append(finalResult, now)
		//fmt.Println("finalResult", finalResult)
	}

	finalPage.Records = models.ToInterfaceSlice(finalResult)
	result := new(models.Result[models.PageResult])
	result.SuccessWithObject(finalPage)
	return e.JSON(200, result)
}

func StartOrStopDish(e echo.Context) error {
	status, _ := strconv.Atoi(e.Param("status"))
	id, _ := strconv.Atoi(e.QueryParam("id"))
	err := db.DB.Model(&models.Dish{Id: int64(id)}).
		Updates(map[string]interface{}{"status": int32(status),
			"update_user": e.Get("empId").(int64),
			"update_time": utils.CustomTime{Time: time.Now()}},
		).Error
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error("修改失败")
		return e.JSON(400, result)
	}
	result := new(models.Result[interface{}])
	result.Success()
	result.Msg = "修改成功"
	return e.JSON(200, result)
}
