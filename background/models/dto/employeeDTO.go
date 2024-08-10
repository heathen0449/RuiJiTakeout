package dto

// 前端提交与数据库差别太大，所以用这个结构体来接收前端提交的数据

type EmployeeDTO struct {
	Id       int64  `json:"id,omitempy"`
	IdNumber string `json:"idNumber"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Sex      string `json:"sex"`
	UserName string `json:"username"`
}
