package match

type UserMatch struct {
	UserID string
}

type Matcher interface {
	Match(fromID string, toID string) error
}
