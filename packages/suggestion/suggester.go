package suggestion

import (
	"context"

	"blinders/packages/common"
)

type Suggester interface {
	ChatCompletion(context.Context, common.UserData, []common.Message, ...Prompter) ([]string, error)
	TextCompletion(context.Context, common.UserData, string) ([]string, error)
}
