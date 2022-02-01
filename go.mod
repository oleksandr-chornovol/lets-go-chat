module github.com/oleksandr-chornovol/lets-go-chat

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deepmap/oapi-codegen v1.9.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/google/wire v0.5.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871 // indirect
	pkg/hasher v1.0.0
)

replace pkg/hasher => ./pkg/hasher
