package database

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"

	"github.com/oleksandr-chornovol/lets-go-chat/config"
	"github.com/oleksandr-chornovol/lets-go-chat/database/drivers"
)

var db *sql.DB

var Driver drivers.DBDriverInterface

var migrations = []string {
	`CREATE TABLE IF NOT EXISTS users(id varchar(100) primary key, name varchar(100), password varchar(100), UNIQUE(name))`,
	`CREATE TABLE IF NOT EXISTS tokens(id varchar(100) primary key, user_id varchar(100), expires_at datetime)`,
}

func Init() {
	dbConfig := config.LocalDBConfig
	database, err := sql.Open(dbConfig["driver"], dbConfig["url"])
	if err != nil {
		log.Println(err)
	}

	db = database

	switch dbConfig["driver"] {
	case "mysql":
		SetDriver(drivers.MySqlDriver{DB: db})
	}
}

func SetDriver(driver drivers.DBDriverInterface) {
	Driver = driver
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
