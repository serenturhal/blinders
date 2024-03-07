package auth

import (
	"context"
	"fmt"

	"blinders/packages/db/repo"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseManager struct {
	App      *firebase.App
	Client   *auth.Client
	UserRepo *repo.UsersRepo
}

func (m FirebaseManager) Verify(jwt string) (*UserAuth, error) {
	authToken, err := m.Client.VerifyIDToken(context.Background(), jwt)
	if err != nil {
		return nil, err
	}

	firebaseUID := authToken.UID
	email := authToken.Claims["email"].(string)
	name := authToken.Claims["name"].(string)

	user, err := m.UserRepo.GetUserByFirebaseUID(firebaseUID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	userAuth := UserAuth{
		Email:  email,
		Name:   name,
		AuthID: firebaseUID,
		ID:     user.ID.Hex(),
	}

	return &userAuth, nil
}

func NewFirebaseManager(userRepo *repo.UsersRepo, adminConfig []byte) (*FirebaseManager, error) {
	manager := FirebaseManager{}
	opt := option.WithCredentialsJSON(adminConfig)
	newApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	manager.App = newApp

	newClient, err := manager.App.Auth(context.Background())
	if err != nil {
		return nil, err
	}
	manager.Client = newClient

	manager.UserRepo = userRepo

	return &manager, nil
}
