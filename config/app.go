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
		//"db_url": os.Getenv("CLEARDB_DATABASE_URL"),
		"db_url": "mysql://b733c72e090730:5245eefd@(https://lets-go-chat-chornovol.herokuapp.com)/heroku_187154ebea373f1",
	},
}

func Get(parameter string) string {
	if os.Getenv("PORT") == "" {
		return app["local"][parameter]
	} else {
		return app["heroku"][parameter]
	}
}
