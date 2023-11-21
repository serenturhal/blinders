package authentication

//
// import (
// 	"os"
// 	"time"
// )
//
// // Public function to outter scope of package
//
// var Manager Maker
//
// func init() {
// 	if Manager == nil {
// 		secretKey := os.Getenv("JWT_SECRET")
// 		durationString := os.Getenv("JWT_DURATION")
// 		tokenDuration, err := time.ParseDuration(durationString)
// 		if err != nil {
// 			Manager, err = NewJWTManager(JwtOptions{
// 				SecretKey: secretKey,
// 			})
// 			if err != nil {
// 				panic(err)
// 			}
// 			return
// 		}
// 		Manager, err = NewJWTManager(JwtOptions{
// 			SecretKey:     secretKey,
// 			TokenDuration: tokenDuration,
// 		})
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }
//
// // Generate jwtToken from userID and userEmail
// func GenerateTokenForUser(userID string, userEmail string) (string, error) {
// 	user := &User{
// 		ID:    userID,
// 		Email: userEmail,
// 	}
// 	return Manager.Generate(user)
// }
//
// // Verify jwtToken, return userID, userEmail, error respectively
// func VerifyToken(token string) (string, string, error) {
// 	user, err := Manager.Verify(token)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	return user.ID, user.Email, nil
// }
