package auth

type UserAuth struct {
	Email  string
	Name   string
	AuthID string
}

// TODO: add interface for role-based auth
type Manager interface {
	Verify(jwt string) (*UserAuth, error)
}
