package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/masonyc/go-gin-example/pkg/settings"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func init() {
	var (
		err            error
		dbType, dbPath string
	)
	sec, err := settings.Cfg.GetSection("database")

	if err != nil {
		log.Fatalf("Fail to get section 'database': %v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbPath = sec.Key("PATH").String()
	db, err = gorm.Open(dbType, dbPath)

	if err != nil {
		log.Println(err)
	}

	db.SingularTable(true)
	db.LogMode(true)
}

func CloseDB() {
	defer db.Close()
}
