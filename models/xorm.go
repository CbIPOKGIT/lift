package models

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

const DATABASE = "nano"

func CreateConnection() (*xorm.Engine, error) {
	return xorm.NewEngine("sqlite3", fmt.Sprintf("%s.db", DATABASE))
}
