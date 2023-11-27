package common

import (
	"fmt"
	"strings"
)

type LanguageContext struct {
	Lang  Lang  `json:"language"`
	Level Level `json:"languageLevel"`
}

func (c LanguageContext) String() string {
	str := strings.Builder{}
	str.WriteString("[")
	str.WriteString(fmt.Sprintf("language: %s, ", c.Lang))
	str.WriteString(fmt.Sprintf("level: %s", c.Level))
	str.WriteString("]")
	return str.String()
}

type UserContext struct {
	UserID   string          `json:"userID"`
	Native   LanguageContext `json:"nativeLanguage"`
	Learning LanguageContext `json:"learningLanguage"`
}

func (c UserContext) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("User: %s\n", c.UserID))
	str.WriteString(fmt.Sprintf("Native language: %s\n", c.Native))
	str.WriteString(fmt.Sprintf("Learning language: %s\n", c.Native))
	return str.String()
}
