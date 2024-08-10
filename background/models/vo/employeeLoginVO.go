package vo

type EmployeeLoginVO struct {
	Id       int    `json:"id"`       // 员工id 数据库主键？
	Name     string `json:"name"`     // 员工姓名
	Token    string `json:"token"`    // 登录token
	UserName string `json:"userName"` // 登录用户名
}
