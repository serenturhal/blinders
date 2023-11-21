package authentication

import (
	"context"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type FirebaseAuthManager struct {
	authClient *auth.Client
}

func NewFirebaseAuthManager() (Maker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// init firebase admin app with GOOGLE_APPLICATION_CREDENTIALS variable
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	manager := &FirebaseAuthManager{
		authClient: authClient,
	}
	return manager, nil
}

func (m *FirebaseAuthManager) Generate(user *User) (string, error) {
	// Generate token are generated with firebase backend server
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// defer cancel()
	// token, err := m.authClient.GetUserByProviderID(ctx, user.ID)
	//
	// if err != nil {
	// 	return "", err
	// }
	// return token, nil
	return "", nil
}

func (m *FirebaseAuthManager) Verify(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	authToken, err := m.authClient.VerifyIDTokenAndCheckRevoked(ctx, token)
	if err != nil {
		return nil, err
	}
	user, err := m.authClient.GetUser(ctx, authToken.UID)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:    authToken.UID,
		Email: user.Email,
	}, nil
}
