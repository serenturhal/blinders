package utils

import (
	"context"
	"os"
	"time"

	"blinders/packages/auth"
	"blinders/packages/common"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var (
	authenticator   auth.Manager
	credentialsPath = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
)

func init() {
	if authenticator == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsPath))
		if err != nil {
			panic(err)
		}
		authenticator, err = auth.NewFirebaseAuthManager(app)
		if err != nil {
			panic(err)
		}
	}
}

func VerifyFireStoreToken(token string) (*common.User, error) {
	return authenticator.Verify(token)
}
