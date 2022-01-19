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
	`CREATE TABLE IF NOT EXISTS users(id varchar(100) primary key, name varchar(100), password varchar(100), last_session_end datetime default current_timestamp, unique(name))`,
	`CREATE TABLE IF NOT EXISTS tokens(id varchar(100) primary key, user_id varchar(100), expires_at datetime)`,
	`CREATE TABLE IF NOT EXISTS messages(id int auto_increment primary key, user_id varchar(100), text varchar(1000), created_at datetime)`,
}

func Init() {
	database, err := sql.Open(config.Get("db_driver"), config.Get("db_url"))
	if err != nil {
		log.Println(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db = database

	log.Println("!!!DB_URL!!!", config.Get("db_url"))
	switch config.Get("db_driver") {
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
