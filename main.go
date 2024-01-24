package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("Could not connect to DB")
	}

	//migrate data and term[late
	_ = db.AutoMigrate(&Schema{})

}
