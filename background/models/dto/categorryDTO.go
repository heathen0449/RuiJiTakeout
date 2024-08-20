package dto

type CategoryDTO struct {
	Id   int64  `json:"id"`
	Type int    `json:"type"` //类型 1 菜品分类 2 套餐分类
	Name string `json:"name"`
	Sort int    `json:"sort"`
}
