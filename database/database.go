package database

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	userModel "github.com/oleksandr-chornovol/lets-go-chat/app/models"
	"github.com/oleksandr-chornovol/lets-go-chat/config"
	"github.com/oleksandr-chornovol/lets-go-chat/database/drivers"
	"log"
	"time"
)

var db *sql.DB

var migrations = []string {
	`CREATE TABLE IF NOT EXISTS users(id varchar(100) primary key, name varchar(100), password varchar(100), UNIQUE(name))`,
}

func Init() {
	//dbConfig := config.LocalDBConfig
	dbConfig := config.HerokuDBConfig
	database, err := sql.Open(dbConfig["driver"], dbConfig["url"])
	if err != nil {
		log.Println(err)
	}

	db = database

	var driver drivers.DBDriverInterface
	switch dbConfig["driver"] {
	case "mysql":
		driver = drivers.MySqlDriver{DB: database}
	}

	userModel.DBDriver = driver
}

func Migrate() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	for _, query := range migrations {
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			log.Println(err)
		}
	}
}
