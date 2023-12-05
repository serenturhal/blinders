package context

import "blinders/packages/common"

func GetUserContext(userID string) (common.UserContext, error) {
	return common.UserContext{
		UserID: userID,
		Native: common.LanguageContext{
			Lang:  common.LangVi,
			Level: common.Intermediate,
		},
		Learning: common.LanguageContext{
			Lang:  common.LangEn,
			Level: common.Beginner,
		},
	}, nil
}
