module github.com/oleksandr-chornovol/lets-go-chat

go 1.16

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	pkg/hasher v1.0.0
)

replace pkg/hasher => ./pkg/hasher
