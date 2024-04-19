package database

import (
	"fmt"

	"github.com/golang-jwt-project/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBInstance() *gorm.DB {
	uri := "root:rootwdp@tcp(localhost:3306)/demo"
	client, err := gorm.Open(mysql.Open(uri), &gorm.Config{})
	if err != nil {
		fmt.Println("error in client :", err)
	}
	client.AutoMigrate(&models.User{})
	return client
}

var DB *gorm.DB = DBInstance()
