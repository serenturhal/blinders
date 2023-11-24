package auth

import (
	"blinders/packages/common"
	"context"
	"time"

	"firebase.google.com/go/v4/auth"
)

type FirebaseAuthManager struct {
	authClient *auth.Client
}

func NewFirebaseAuthManager(authClient *auth.Client) (AuthManager, error) {
	manager := &FirebaseAuthManager{
		authClient: authClient,
	}
	return manager, nil
}

func (m *FirebaseAuthManager) Generate(user *common.User) (string, error) {
	return "", nil
}

func (m *FirebaseAuthManager) Verify(token string) (*common.User, error) {
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
	return &common.User{
		ID:    authToken.UID,
		Email: user.Email,
	}, nil
}
