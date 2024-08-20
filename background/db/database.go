package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() {
	// 数据库连接配置
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	dsn := "root:ljh200050@tcp(127.0.0.1:3306)/project?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 强制使用单数表名
		},
	})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}
	fmt.Println("Database connection established")
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("Failed to get database connection:", err)
	}
	if err := sqlDB.Close(); err != nil {
		fmt.Println("Failed to close database connection:", err)
	}
	fmt.Println("Database connection closed")
}
