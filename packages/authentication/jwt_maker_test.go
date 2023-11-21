package authentication

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func randomUser() *User {
	return &User{
		Email: strconv.FormatInt(int64(rand.Intn(1000)), 10),
		ID:    strconv.FormatInt(int64(rand.Intn(1000)), 10),
	}
}

func TestJWT(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		Name          string        // name of testcase
		Secret        string        // secret that passed in construct
		Duration      time.Duration // token duration
		Sleep         time.Duration // time sleep before verify token
		User          *User         // expect user claim
		VerifyToken   string        // token using to test verify, if empty, use the generated token
		ExpectedError error         // error expected for testcase
	}{
		{
			Name:          "Success",
			Secret:        "sample",
			Duration:      time.Second,
			User:          randomUser(),
			ExpectedError: nil,
		},
		{
			Name:          "empty secret",
			Duration:      time.Second,
			User:          randomUser(),
			ExpectedError: ErrEmptySecret,
		},
		{
			Name:          "invalid token",
			Secret:        "sample",
			Duration:      time.Second,
			User:          randomUser(),
			VerifyToken:   "invalid token",
			ExpectedError: ErrInvalidToken,
		},
	}

	for index := range testcases {
		tc := testcases[index]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			opts := JwtOptions{
				SecretKey:     tc.Secret,
				TokenDuration: tc.Duration,
			}
			manager, err := NewJWTManager(opts)
			// If there is error, then error must equal to expected error
			if err != nil {
				assert.Equal(t, tc.ExpectedError, err)
				return
			}
			assert.NotNil(t, manager)

			token, err := manager.Generate(tc.User)
			if err != nil {
				assert.Equal(t, tc.ExpectedError, err)
				return
			}
			// token must not be empty
			assert.NotEqual(t, "", token)
			time.Sleep(tc.Sleep)

			if tc.VerifyToken != "" {
				token = tc.VerifyToken
			}

			user, err := manager.Verify(token)
			if err != nil {
				assert.Equal(t, tc.ExpectedError, err)
				return
			}
			assert.Equal(t, tc.User, user)
		})
	}
}

func TestExpiredToken(t *testing.T) {
	t.Parallel()
	opts := JwtOptions{
		SecretKey:     "sample",
		TokenDuration: time.Second,
	}
	manager, err := NewJWTManager(opts)
	assert.Nil(t, err)
	assert.NotNil(t, manager)
	user := randomUser()
	jwtToken, err := manager.Generate(user)
	assert.Nil(t, err)
	assert.NotEmpty(t, jwtToken)
	time.Sleep(time.Second * 2)

	usr, err := manager.Verify(jwtToken)
	assert.NotNil(t, err)
	assert.Nil(t, usr)
}
