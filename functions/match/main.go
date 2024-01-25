package main

import (
	"context"
	"fmt"

	"blinders/packages/match"

	"github.com/aws/aws-lambda-go/lambda"
)

var matcher match.Matcher

func init() {
	matcher = &match.MockMatcher{}
}

type MatchEvent struct {
	FromID string `json:"fromId"`
	ToID   string `json:"toId"`
}

func HandleRequest(_ context.Context, event *MatchEvent) error {
	if event == nil {
		return fmt.Errorf("received nil event")
	}

	if err := matcher.Match(event.FromID, event.ToID); err != nil {
		return fmt.Errorf("cannot match players, err: %v", err)
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
