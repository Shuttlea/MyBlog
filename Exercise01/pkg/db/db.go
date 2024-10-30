package db

import (
	"bufio"
	"fmt"
	"log"
	"myblog/structs"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadCfg(n int) (string, string) {
	file, err := os.Open("admin_credentials.txt")
	if err != nil {
		log.Println("Can't open admin_credetials.txt : ", err)
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for ; n != 0; n-- {
		scanner.Scan()
	}
	scanner.Scan()
	first := scanner.Text()
	scanner.Scan()
	second := scanner.Text()
	return first, second
}

func ConnectDB() *gorm.DB {
	database, user := LoadCfg(2)
	dsn := fmt.Sprintf("host=localhost user=%s database=%s port=5432 sslmode=disable", user, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Can't open database: ", err.Error())
	}
	err = db.AutoMigrate(&structs.Article{})
	if err != nil {
		log.Println("Can't megrate database: ", err.Error())
		panic(err)
	}
	return db
}

func CloseDB(db *gorm.DB) {
	dbInstance, _ := db.DB()
	_ = dbInstance.Close()
}

func ReadDB(page int, db *gorm.DB) structs.Data1 {
	var total int64
	db.Table("articles").Count(&total)
	if page < 1 || page*3-int(total) > 2 {
		page = 1
	}
	limit := 3
	totalPages := int(total / 3)
	offset := (totalPages - page) * 3
	if total%3 != 0 {
		totalPages++
		offset += int(total % 3)
		if page == totalPages {
			limit = int(total % 3)
			offset = 0
		}
	}
	var res structs.Data1
	result := db.Offset(offset).Limit(limit).Find(&res.Nodes)
	if result.Error != nil {
		log.Println("error to find in database: ", result.Error)
	}
	res.Page = page
	res.TotalPages = totalPages

	return res
}

func ReadArticle(a int, db *gorm.DB) structs.Article {
	if a == 0 {
		a = 1
	}
	var res structs.Article
	db.Where("id = ?", strconv.Itoa(a)).First(&res)

	return res
}

func WriteDB(a structs.Article, db *gorm.DB) {
	var id int64
	db.Table("articles").Count(&id)
	a.ID = int(id) + 1

	result := db.Create(&a)
	if result.Error != nil {
		log.Println("Can't create row: ", result.Error)
	}
}
