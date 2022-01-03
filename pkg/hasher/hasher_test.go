package hasher

import (
	"fmt"
	"testing"
)

func ExampleCheckPasswordHash() {
	fmt.Println(CheckPasswordHash("password", "$2a$10$Kt0YB3SgXJuUpek5anTDguHyKXUEbE4EIyzQXrfzYzsNB9ExZflSe"))
	//Output:true
}

func BenchmarkHashPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HashPassword("password")
	}
}
