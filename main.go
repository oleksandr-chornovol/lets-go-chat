package main

import (
	"fmt"
	"pkg/hasher"
)

func main() {
	password := "password"
	hash, _ := hasher.HashPassword(password)
	fmt.Println(hash)
	fmt.Println(hasher.CheckPasswordHash(password, hash))
}
