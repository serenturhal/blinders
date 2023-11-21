package token

import "blinders/packages/authentication/models"

type Maker interface {
	Generate(user *models.User) (string, error)
	Verify(token string) (*models.User, error)
}
