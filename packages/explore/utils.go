package explore

import "fmt"

const matchKey = "match:%v"

func CreateMatchKeyWithUserID(userID string) string {
	return fmt.Sprintf(matchKey, userID)
}
