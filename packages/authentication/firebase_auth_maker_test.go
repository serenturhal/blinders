package authentication

import (
	"context"
	"io"
	"testing"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/stretchr/testify/assert"
)

func TestFirebaseSuccess(t *testing.T) {
	app, err := firebase.NewApp(context.Background(), nil)
	assert.Nil(t, err)
	assert.NotNil(t, app)

	authclient, err := app.Auth(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, authclient)

	maker, err := NewFirebaseAuthManager()
	assert.Nil(t, err)
	assert.NotNil(t, maker)
	var (
		validUser = &User{
			ID:    "t7ZYtyjYCbMxOefUALu8b2P4AVO2",
			Email: "minhdat15012002@gmail.com",
		}
	)
	validToken, err := GetUserIDToken(validUser)
	// validToken, err := maker.Generate(validUser)
	assert.Nil(t, err)
	assert.NotEmpty(t, validToken)
	u, err := maker.Verify(validToken)
	assert.Nil(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, u, validUser)
}
func TestFirebaseWithInvalidUser(t *testing.T) {
	app, err := firebase.NewApp(context.Background(), nil)
	assert.Nil(t, err)
	assert.NotNil(t, app)

	authclient, err := app.Auth(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, authclient)

	maker, err := NewFirebaseAuthManager()
	assert.Nil(t, err)
	assert.NotNil(t, maker)
	var (
		invalidToken = "invalidTOken"
	)
	u, err := maker.Verify(invalidToken)

	// validToken, err := maker.Generate(validUser)
	assert.NotNil(t, err)
	assert.Nil(t, u)
}

func GetUserIDToken(user *User) (string, error) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return "", err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return "", err
	}
	tokenString, err := client.CustomToken(ctx, user.ID)
	if err != nil {
		return "", err
	}

	idToken, err := signInWithCustomToken(tokenString)
	if err != nil {
		return "", err
	}
	return idToken, nil
}

const (
	verifyCustomTokenURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s"
)

var (
	apiKey = os.Getenv("API_KEY")
)

// see https://github.com/firebase/firebase-admin-go/blob/1d2a52c3c8195451b5ad2e0a173906bd6eb9529d/integration/auth/auth_test.go#L199
func signInWithCustomToken(token string) (string, error) {
	req, err := json.Marshal(map[string]interface{}{
		"token":             token,
		"returnSecureToken": true,
	})
	if err != nil {
		return "", err
	}

	resp, err := postRequest(fmt.Sprintf(verifyCustomTokenURL, apiKey), req)
	if err != nil {
		return "", err
	}
	var respBody struct {
		IDToken string `json:"idToken"`
	}
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		return "", err
	}
	return respBody.IDToken, nil
}

func postRequest(url string, req []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
