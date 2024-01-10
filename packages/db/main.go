package db

import (
	"fmt"
	"strings"

	"blinders/packages/common"
)

type UserData struct {
	UserID   string          `json:"userID"`
	Native   common.Language `json:"nativeLanguage"`
	Learning common.Language `json:"learningLanguage"`
}

func (d UserData) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("Native language: %s\n", d.Native))
	str.WriteString(fmt.Sprintf("Learning language: %s\n", d.Native))
	return str.String()
}

func GetUserData(userID string) (UserData, error) {
	return UserData{
		UserID: userID,
		Native: common.Language{
			Lang:  common.LangVi,
			Level: common.Intermediate,
		},
		Learning: common.Language{
			Lang:  common.LangEn,
			Level: common.Beginner,
		},
	}, nil
}
