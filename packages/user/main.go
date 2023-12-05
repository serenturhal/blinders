package user

import "blinders/packages/common"

func GetUserData(userID string) (common.UserData, error) {
	return common.UserData{
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
