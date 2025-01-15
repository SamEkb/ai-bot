package openai

import (
	"context"
	"errors"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Service struct {
	client *openai.Client
	chats  Chats
}

func NewService(cfg *Config) (*Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	client := openai.NewClient(option.WithAPIKey(cfg.Token))

	return &Service{client: client, chats: []*Chat{}}, nil
}

func (s *Service) NewConversation(_ context.Context, id int64) {
	s.chats.Reset(id)
}

func (s *Service) ChatCompletion(ctx context.Context, id int64, prompt string) (string, error) {
	chat, err := s.chats.Find(id)
	if err != nil && !errors.Is(err, ErrNoChat) {
		return "", err
	}

	newMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatCompletionMessageRole(openai.ChatCompletionMessageParamRoleUser),
		Content: prompt,
	}

	if errors.Is(err, ErrNoChat) {
		s.chats, chat = s.chats.Create(id, newMessage)
	} else {
		chat.AddMessage(newMessage)
	}

	var messages []openai.ChatCompletionMessageParamUnion
	for _, m := range chat.Messages {
		switch m.Role {
		case "system":
			messages = append(messages, openai.SystemMessage(m.Content))
		case "assistant":
			messages = append(messages, openai.AssistantMessage(m.Content))
		case "user":
			messages = append(messages, openai.UserMessage(m.Content))
		default:
			messages = append(messages, openai.UserMessage(m.Content))
		}
	}

	chatCompletion, err := s.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F(messages),
			Model:    openai.F(openai.ChatModelGPT3_5Turbo),
		},
	)
	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", errors.New("no choices returned from chat completion")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}
