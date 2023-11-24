package common

import (
	"fmt"
	"strings"
)

type User struct {
	ID    string `json:"userID"`
	Email string `json:"userEmail"`
}

func (u User) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("UserID: %s", u.ID))
	str.WriteString(fmt.Sprintf("Email: %s", u.Email))
	return str.String()
}
