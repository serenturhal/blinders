package main

import (
	"blinders/packages/authentication"
	"fmt"
)

func main() {
	token, err := authentication.GenerateTokenForUser("t7ZYtyjYCbMxOefUALu8b2P4AVO2", "minhdat15012002@gmail.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
}
