package main

import (
	"blinders/packages/authentication"
	"fmt"
	"time"

	"os"
)

var (
	defaultTokenDuration = time.Minute * 15
)

func main() {
	secretKey := os.Getenv("JWT_SECRET")
	durationString := os.Getenv("JWT_DURATION")
	tokenDuration, err := time.ParseDuration(durationString)
	if err != nil {
		tokenDuration = defaultTokenDuration
	}

	maker, err := authentication.NewJWTManager(authentication.JwtOptions{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration,
	})
	if err != nil {
		panic(err)
	}

	token, err := maker.Generate(&authentication.User{
		ID:    "t7ZYtyjYCbMxOefUALu8b2P4AVO2",
		Email: "minhdat15012002@gmail.com",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(token)
}
