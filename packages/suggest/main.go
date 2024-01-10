package suggest

import (
	"context"

	"blinders/packages/db"
)

type Prompter interface {
	Build() (string, error)
	Update(...interface{}) error
}

type Suggester interface {
	ChatCompletion(context.Context, db.UserData, []Message, ...Prompter) ([]string, error)
	TextCompletion(context.Context, db.UserData, string) ([]string, error)
}
