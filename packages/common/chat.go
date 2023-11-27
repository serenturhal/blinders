package common

const (
	Individual ChatType = "individual"
	Group      ChatType = "group"
)

func (t ChatType) String() string {
	switch t {
	case Individual:
		return "individual"
	case Group:
		return "individual"
	default:
		return "unknown"
	}
}

type (
	ChatType string
	ChatRoom struct {
		Type    ChatType `json:"type"`
		UserIDs []string `json:"uids"`
	}
)
