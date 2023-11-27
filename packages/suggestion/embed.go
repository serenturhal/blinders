package suggestion

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	randSource             *rand.Rand = nil
	messageSuggestionEmbed            = []string{
		`
Sender context:
    Language: %s
    Level: %s
Recent messages: [
%s
]
Prompt:
    Given the recent messages in the conversation, the user's context,
    provide a message that the user could send next.
    Ensure that the suggestion is contextually relevant, considerate, and aligned with the user's language proficiency.
    Aim for a response that flows seamlessly within the ongoing conversation.
    Just return the completed message, not any guides.`,
	}
)

func init() {
	randSource = rand.New(rand.NewSource(time.Now().Unix()))
}

func randomEmbed(i ...int) (string, error) {
	if len(i) == 1 {
		if 0 > i[0] || i[0] >= len(messageSuggestionEmbed) {
			return "", fmt.Errorf("embed: index out of length, got (%d)", i[0])
		}
		return messageSuggestionEmbed[i[0]], nil
	}
	return messageSuggestionEmbed[randSource.Intn(len(messageSuggestionEmbed))], nil
}
