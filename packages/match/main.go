package match

type UserMatch struct {
	UserID   string
	Name     string
	Gender   string
	Native   string
	Learning string
	Age      int
}

type Matcher interface {
	Match(fromID string, toID string) error
	// Suggest returns list of users that maybe match with given user
	Suggest(id string) ([]UserMatch, error)
}
