package main

import (
	"blinders/packages/authentication"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	tokenMaker      authentication.Maker
	defaultDuration = time.Minute * 15
)

type User struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
}

func init() {
	secretKey := os.Getenv("JWT_SECRET")
	durationString := os.Getenv("JWT_DURATION")
	tokenDuration, err := time.ParseDuration(durationString)
	if err != nil {
		tokenDuration = defaultDuration
	}
	tokenMaker, err = authentication.NewJWTManager(authentication.JwtOptions{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration,
	})
	if err != nil {
		panic(fmt.Sprintf("Cannot init jwt maker: %s", err.Error()))
	}
}

func handleRequest(ctx context.Context, user *User) (*string, error) {
	if user == nil {
		return nil, errors.New("Nil user")
	}

	token, err := tokenMaker.Generate(&authentication.User{
		ID:    user.UserID,
		Email: user.UserEmail,
	})
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func main() {
	lambda.Start(handleRequest)
}
