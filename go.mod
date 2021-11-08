module github.com/oleksandr-chornovol/lets-go-chat

go 1.16

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/tools/godep v0.0.0-20180126220526-ce0bfadeb516 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20211107104306-e0b2ad06fe42 // indirect
	golang.org/x/tools v0.1.7 // indirect
	pkg/hasher v1.0.0
)

replace pkg/hasher => ./pkg/hasher
