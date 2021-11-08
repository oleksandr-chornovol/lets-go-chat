module github.com/oleksandr-chornovol/lets-go-chat

go 1.16

require (
	github.com/google/uuid v1.3.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	pkg/hasher v1.0.0
)

replace pkg/hasher => ./pkg/hasher
