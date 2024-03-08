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

	if m.UserRepo == nil {
		return nil, fmt.Errorf("firebaseMangager: userRepo must not be nil")
	}

	userAuth := UserAuth{
		Email:  email,
		Name:   name,
		AuthID: firebaseUID,
	}

	return &userAuth, nil
}

func NewFirebaseManager(adminConfig []byte) (*FirebaseManager, error) {
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

	return &manager, nil
}
