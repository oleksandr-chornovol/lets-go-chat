package main

import (
	"github.com/oleksandr-chornovol/lets-go-chat/app/http"
	"github.com/oleksandr-chornovol/lets-go-chat/database"
)

//var database *sql.DB
//
//type User struct{
//	Id int
//	Email string
//	Created_at string
//	Is_active string
//}

func main() {
	database.Init()
	database.Migrate()
	http.InitRoutes()
	http.StartServer()

	//db, err := sql.Open("mysql", "root:root@/lets-go-chat")
	//
	//if err != nil {
	//	//log.Println(err)
	//}
	//database = db
	//defer db.Close()
	//
	////results, err := database.Query("select * from blast_bs.open_account_request")
	//res, err := database.Exec("insert into users (id, name, password) values (?, ?, ?)",
	//	"id", "name", "pass")
	//
	////for results.Next() {
	////	var user User
	////	// for each row, scan the result into our tag composite object
	////	err = results.Scan(&user.Id, &user.Email, &user.Created_at, &user.Is_active)
	//if err != nil {
	//	panic(err.Error()) // proper error handling instead of panic in your app
	//} else {
	//	log.Println(res)
	//}
	//
	////	// and then print out the tag's Name attribute
	////	log.Printf(user.Email)
	////}
}
