package auth

type UserAuth struct {
	Email  string
	Name   string
	AuthID string
}

type Manager interface {
	Verify(jwt string) (*UserAuth, error)
}
