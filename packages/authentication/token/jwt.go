package token

import (
	"blinders/packages/authentication/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	DefaultTokenDuration = time.Minute * 15
)

type JwtOptions struct {
	SecretKey     string
	TokenDuration time.Duration
}

type JWTManager struct {
	JwtOptions
}

var (
	ErrUnexpectedSigningMethod = errors.New("Unexpected signin method")
	ErrEmptySecret             = errors.New("Secret key must not be empty")
	ErrInvalidToken            = errors.New("Invalid token")
	ErrInvalidClaims           = errors.New("Token have invalid claims")
)

type Payload struct {
	UserID string `json:"user_id"` // Firebase uid of user
	Email  string `json:"email"`   // Gmail address of user
	jwt.RegisteredClaims
}

func NewJWTManager(opts JwtOptions) (Maker, error) {
	if len(opts.SecretKey) == 0 {
		return nil, ErrEmptySecret
	}
	// if tokenDuration not specify
	if opts.TokenDuration == 0 {
		opts.TokenDuration = DefaultTokenDuration
	}
	return &JWTManager{
		JwtOptions: opts,
	}, nil
}

func (m *JWTManager) Generate(user *models.User) (string, error) {
	claims := &Payload{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.TokenDuration)), // unix timestamp
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(m.SecretKey))
	return ss, err
}

func (m *JWTManager) Verify(token string) (*models.User, error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&Payload{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return []byte(m.SecretKey), nil
		},
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// convert claims to userclaims
	claims, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidClaims
	}
	return &models.User{
		ID:    claims.UserID,
		Email: claims.Email,
	}, nil
}
