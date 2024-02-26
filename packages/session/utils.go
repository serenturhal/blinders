package session

func ConstructUserKey(userID string) string {
	return "user:" + userID
}

func ConstructConnectionKey(connectionID string) string {
	return "connection:" + connectionID
}
