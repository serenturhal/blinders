package auth

import "blinders/packages/common"

type Manager interface {
	Generate(user *common.User) (string, error)
	Verify(token string) (*common.User, error)
}
