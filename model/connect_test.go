package model

import (
	"testing"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
)

func TestGetProductById(t *testing.T) {
	db, err := gorm.Open("mysql", "michong:michong@tcp(192.168.1.34)/ck_dev?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&ProductViewedDetail{})
	fmt.Println(db.HasTable(&Product{}))

	db.Create(&Product{
		Title: "test",
		Amount: 200,
	})

	var product Product
	db.First(&product, 1) // find product with id 1
	db.First(&product, "Title = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	fmt.Println(product)

	// Delete - delete product
	//db.Delete(&product)
}