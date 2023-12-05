package common

import (
	"fmt"
	"strings"
)

type UserData struct {
	UserID   string   `json:"userID"`
	Native   Language `json:"nativeLanguage"`
	Learning Language `json:"learningLanguage"`
}

func (d UserData) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("Native language: %s\n", d.Native))
	str.WriteString(fmt.Sprintf("Learning language: %s\n", d.Native))
	return str.String()
}
