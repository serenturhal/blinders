package authentication

import (
	"blinders/packages/authentication/models"
	"blinders/packages/authentication/token"
	"os"
	"time"
)

var Manager token.Maker

func init() {
	if Manager == nil {
		secretKey := os.Getenv("JWT_SECRET")
		durationString := os.Getenv("JWT_DURATION")
		tokenDuration, err := time.ParseDuration(durationString)
		if err != nil {
			Manager, err = token.NewJWTManager(token.JwtOptions{
				SecretKey: secretKey,
			})
			if err != nil {
				panic(err)
			}
			return
		}
		Manager, err = token.NewJWTManager(token.JwtOptions{
			SecretKey:     secretKey,
			TokenDuration: tokenDuration,
		})
		if err != nil {
			panic(err)
		}
	}
}

// Generate jwtToken from userID and userEmail
func GenerateTokenForUser(userID string, userEmail string) (string, error) {
	user := &models.User{
		ID:    userID,
		Email: userEmail,
	}
	return Manager.Generate(user)
}

// Verify jwtToken, return userID, userEmail, error respectively
func VerifyToken(token string) (string, string, error) {
	user, err := Manager.Verify(token)
	if err != nil {
		return "", "", err
	}
	return user.ID, user.Email, nil
}
