package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	res, err := bcrypt.GenerateFromPassword([]byte("asdqwezxc1488"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
}
