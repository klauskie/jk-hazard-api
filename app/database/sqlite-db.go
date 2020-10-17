package database

import (
	"fmt"
	"klaus.com/jkapi/app/entity"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB : database variable
var DB *gorm.DB

const (
	DbNAME = "jk-hazard.db"
	DbNameFULL = "/go/src/app/tmp/" + DbNAME
)

func init() {
	var err error
	DB, err = gorm.Open("sqlite3", DbNAME)
	if err != nil {
		panic(err)
	}
	if err = DB.DB().Ping(); err != nil {
		panic(err)
	}

	fmt.Println("You are connected to database jk-hazard...")

	DB.AutoMigrate(&entity.User{})
}