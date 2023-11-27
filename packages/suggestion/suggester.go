package suggestion

import (
	"blinders/packages/common"
	"context"
)

type Suggester interface {
	ChatCompletion(ctx context.Context, userContext common.UserContext, recentMessages []common.Message) ([]string, error)
	TextCompletion(ctx context.Context, prompt string) (string, error)
}
