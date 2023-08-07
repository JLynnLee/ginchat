package main

import (
	"fmt"
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/ginchat"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	err = db.AutoMigrate(&models.Message{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&models.Group{})
	if err != nil {
		return
	}
	//db.AutoMigrate(&models.UserBasic{})
	return
	// Create
	//db.Create(&Product{Code: "D42", Price: 100})
	user := &models.UserBasic{}
	//user.Name = "张三"
	//user.PassWord = "12134"
	//db.Create(user)
	// Read
	fmt.Println(db.First(user, 1)) // 根据整型主键查找
	//db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Model(user).Update("PassWord", 123433)
	// Update - 更新多个字段
	//db.Model(user).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	//db.Model(user).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(user, 7)
}
