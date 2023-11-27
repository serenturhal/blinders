package suggestion

import (
	"blinders/packages/common"
	"context"
	"errors"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

var DefaultSuggesterOptions = GPTSuggesterOptions{
	ChatModel:         openai.GPT3Dot5TurboInstruct,
	TextModel:         openai.GPT3Dot5TurboInstruct,
	NChat:             3,
	NText:             1,
	ModelTemperateure: 0.6,
}

type GPTSuggesterOptions struct {
	ChatModel         string
	TextModel         string
	NChat             int
	NText             int
	ModelTemperateure float32
}

type GPTSuggester struct {
	prompter Prompter
	client   *openai.Client
	GPTSuggesterOptions
}

func (s *GPTSuggester) ChatCompletion(ctx context.Context, userContext common.UserContext, msgs []common.Message) ([]string, error) {
	suggestions := []string{}
	err := s.prompter.Update(userContext, msgs)
	if err != nil {
		return suggestions, err
	}

	prompt, err := s.prompter.Build()
	if err != nil {
		return suggestions, err
	}

	req := openai.CompletionRequest{
		Model:       s.ChatModel,
		Prompt:      prompt,
		N:           s.NChat,
		Temperature: 0.6,
	}
	rsp, err := s.client.CreateCompletion(ctx, req)
	if err != nil {
		return suggestions, err
	}

	if len(rsp.Choices) == 0 {
		return suggestions, errors.New("gptSuggester: got empty reply from server")
	}

	for _, choice := range rsp.Choices {
		fmt.Printf("choice: %v\n", choice)
		suggestions = append(suggestions, choice.Text)
	}
	return suggestions, nil
}

func (s *GPTSuggester) TextCompletion(ctx context.Context, prompt string) (string, error) {
	req := openai.CompletionRequest{
		Model:       s.TextModel,
		Prompt:      prompt,
		N:           s.NText,
		Temperature: s.ModelTemperateure,
	}
	rsp, err := s.client.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(rsp.Choices) == 0 {
		return "", errors.New("gptSuggester: got empty reply from server")
	}
	return rsp.Choices[0].Text, nil
}

func (s *GPTSuggester) _mustImplementSuggester() {
	var _ Suggester = s
}

func NewGPTSuggestor(client *openai.Client, prompter Prompter, opts ...Option) (*GPTSuggester, error) {
	gptSuggester := &GPTSuggester{
		client:              client,
		prompter:            prompter,
		GPTSuggesterOptions: DefaultSuggesterOptions,
	}
	for _, opt := range opts {
		opt(gptSuggester)
	}
	return gptSuggester, nil
}

func optionAdapter(closer func(s *GPTSuggester)) Option {
	return func(i any) {
		switch s := i.(type) {
		case *GPTSuggester:
			closer(s)
		}
	}
}

func WithTemperature(temperature float32) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.ModelTemperateure = temperature
	})
}

func WithNText(N int) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.NText = N
	})
}

func WithNChat(N int) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.NChat = N
	})
}

func WithTextModel(model string) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.TextModel = model
	})
}

func WithChatModel(model string) Option {
	return optionAdapter(func(s *GPTSuggester) {
		s.ChatModel = model
	})
}
