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
		"db_url": "iuy6iqcs4m7sbfz1:k0kyjat4v8hesrn2@d3y0lbg7abxmbuoi.chr7pe7iynqr.eu-west-1.rds.amazonaws.com:3306/ojrqgsiogkv0qqlz",
	},
}

func Get(parameter string) string {
	if os.Getenv("PORT") == "" {
		return app["local"][parameter]
	} else {
		return app["heroku"][parameter]
	}
}
