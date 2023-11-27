package suggestion

import (
	"blinders/packages/common"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildPrompt(t *testing.T) {
	_, _, _, prompter := initPrompt(t)
	prompt, err := prompter.Build()
	assert.Nil(t, err)
	assert.NotEmpty(t, prompt)
}

func TestUpdatePrompt(t *testing.T) {
	sender, receiver, userContext, prompter := initPrompt(t)
	assert.Equal(t, userContext, prompter.ctx)

	newContext := newUserContext(
		&receiver,
		common.LanguageContext{},
		common.LanguageContext{},
	)
	assert.Nil(t, prompter.Update(newContext))
	assert.Equal(t, newContext, prompter.ctx)

	lenMessages := 5
	msgs := []common.Message{}
	for i := 0; i < lenMessages; i++ {
		switch i % 2 {
		case 0:
			msgs = append(msgs, *common.NewMessage(sender.ID, receiver.ID, fmt.Sprintf("test_%d", i)))
		default:
			msgs = append(msgs, *common.NewMessage(receiver.ID, sender.ID, fmt.Sprintf("test_%d", i)))
		}
	}
	assert.Nil(t, prompter.Update(msgs))
	assert.Equal(t, msgs, prompter.messages)

	assert.NotNil(t, prompter.Update(sender))
}

func initPrompt(t *testing.T) (sender common.User, receiver common.User, senderContext common.UserContext, prompter MessageSuggestionPrompt) {
	sender = newUser("sender", "sender@email")
	receiver = newUser("receiver", "receiver@email")
	senderContext = newUserContext(
		&sender,
		common.LanguageContext{
			Lang:  common.LangVi,
			Level: common.Beginner,
		}, common.LanguageContext{
			Lang:  common.LangEn,
			Level: common.Beginner,
		})
	lenMessages := 10
	msgs := []common.Message{}
	for i := 0; i < lenMessages; i++ {
		switch i % 2 {
		case 0:
			msgs = append(msgs, *common.NewMessage(sender.ID, receiver.ID, fmt.Sprintf("msg_%d", i)))
		default:
			msgs = append(msgs, *common.NewMessage(receiver.ID, sender.ID, fmt.Sprintf("msg_%d", i)))
		}
	}

	prompter = NewMessageSuggestionPrompt()
	assert.NotNil(t, prompter)
	assert.Nil(t, prompter.Update(senderContext, msgs))
	return
}

func newUserContext(user *common.User, native common.LanguageContext, language common.LanguageContext) common.UserContext {
	return common.UserContext{
		UserID:   user.ID,
		Native:   native,
		Learning: language,
	}
}

func newUser(id string, email string) common.User {
	return common.User{
		ID:    id,
		Email: email,
	}
}
