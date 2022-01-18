package config

import "os"

var app = map[string]map[string]string{
	"local": {
		"port": "8080",
		"db_driver": "mysql",
		"db_url": "root:root@/lets-go-chat",
	},
	"heroku": {
		"port": os.Getenv("PORT"),
		"db_driver": "mysql",
		//"db_url": os.Getenv("DATABASE_URL"),
		"db_url": os.Getenv("CLEARDB_DATABASE_URL"),
	},
}

func Get(parameter string) string {
	if os.Getenv("PORT") == "" {
		return app["local"][parameter]
	} else {
		return app["heroku"][parameter]
	}
}
