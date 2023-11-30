package suggestion

import (
	"blinders/packages/common"
	"errors"
	"fmt"
	"strings"
)

type MessageSuggestionPrompt struct {
	embed    string // embed string that could be use to embed user's context and messages to make complete prompt
	ctx      common.UserContext
	messages []common.Message
}

func (p MessageSuggestionPrompt) Build() (string, error) {
	msgs := []string{}
	for _, msg := range p.messages {
		switch msg.FromID {
		case p.ctx.UserID:
			msgs = append(msgs, fmt.Sprintf("\tsender: %s", msg.Content))
		default:
			msgs = append(msgs, fmt.Sprintf("\treceiver: %s", msg.Content))
		}
	}
	return fmt.Sprintf(p.embed, p.ctx.Learning.Lang, p.ctx.Learning.Level, strings.Join(msgs, "\n")), nil
}

func (p *MessageSuggestionPrompt) Update(objs ...any) error {
	for _, obj := range objs {
		switch doc := obj.(type) {
		case common.UserContext:
			p.ctx = doc
		case *common.UserContext:
			p.ctx = *doc
		case []common.Message:
			p.messages = doc
		default:
			return errors.New("messageSuggestionPrompter: expected(*common.UserContext, []common.Message) got unknown")
		}
	}
	return nil
}

func NewMessageSuggestionPrompt() *MessageSuggestionPrompt {
	embed, _ := randomEmbed()

	return &MessageSuggestionPrompt{
		embed: embed,
	}
}
