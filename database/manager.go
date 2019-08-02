package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

var db *gorm.DB

var dbpast *gorm.DB

func GetDB() *gorm.DB {
	return db
}
func GetPastDB() *gorm.DB {
	return dbpast
}
func Connect() {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	host := os.Getenv("db_host")
	port := os.Getenv("db_port")
	database := os.Getenv("db_name")
	database1 := os.Getenv("db_name1")
	dialect := os.Getenv("db_type")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC", username, password, host, port, database)

	log.Println("Connecting to database.")
	log.Println(fmt.Sprintf("Database credentials: \n"+
		"	Host: %s\n"+
		"	Port: %s\n"+
		"	Database name: %s\n"+
		"	Username: %s\n"+
		"	Password: %s", host, port, database, username, password))
	var err error
	db, err = gorm.Open(dialect, dataSource)
	if err == nil {
		log.Println("Connected successfully1.")
	} else {
		log.Fatal("Connection failed.", err)
	}

	dataSource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC", username, password, host, port, database1)
	dbpast, err = gorm.Open(dialect, dataSource)
	if err == nil {
		log.Println("Connected successfully2.")
	} else {
		log.Fatal("Connection failed.", err)
	}
//	db.LogMode(true)
//	drop()
	db.LogMode(true)
	migrate()
	foreignkey()
	index()
//	insertSome()
}
