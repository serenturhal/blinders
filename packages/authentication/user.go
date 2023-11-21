package authentication

// User type in Authenticate scope, maybe use User type in common package
type User struct {
	ID    string `json:"user_id"`    // firebase id of user
	Email string `json:"user_email"` // gmail address of user
}
