package authentication

type Maker interface {
	Generate(user *User) (string, error)
	Verify(token string) (*User, error)
}
