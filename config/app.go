package config

import "os"

var LocalDBConfig = map[string]string {
	"driver": "mysql",
	"url": "root:root@/lets-go-chat",
}

var HerokuDBConfig = map[string]string {
	"driver": "mysql",
	"url": os.Getenv("DATABASE_URL"),
}

var LocalPort = "8080"
