package bot

import (
	"context"
	tb "gopkg.in/telebot.v4"
)

var (
	menu               = &tb.ReplyMarkup{ResizeKeyboard: true}
	btnNewConversation = menu.Text("⚙ New conversation")
)

type openAiService interface {
	ChatCompletion(context.Context, int64, string) (string, error)
	NewConversation(ctx context.Context, id int64)
}

type Wrapper struct {
	bot       *tb.Bot
	config    *Config
	openaiSvc openAiService
}

func NewWrapper(config *Config, openaiService openAiService) (*Wrapper, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	settings := tb.Settings{
		Token:  config.Token,
		Poller: &tb.LongPoller{Timeout: config.Timeout},
	}

	bot, err := tb.NewBot(settings)
	if err != nil {
		return nil, err
	}

	w := &Wrapper{bot: bot, config: config, openaiSvc: openaiService}

	w.prepare()

	return w, nil
}

func (w *Wrapper) Start(_ context.Context) error {
	w.bot.Start()

	return nil
}

func (w *Wrapper) prepare() {

	menu.Reply(
		menu.Row(btnNewConversation),
	)

	w.bot.Handle(&btnNewConversation, w.newConversationHandler)
	w.bot.Handle(tb.OnText, w.chatCompletionHandler)
}

func (w *Wrapper) newConversationHandler(tbContext tb.Context) error {
	ctx := context.TODO()
	w.openaiSvc.NewConversation(ctx, tbContext.Chat().ID)
	return tbContext.Send("New conversation has been started", menu)
}

func (w *Wrapper) chatCompletionHandler(tbContext tb.Context) error {
	ctx := context.TODO()
	text := tbContext.Text()

	completion, err := w.openaiSvc.ChatCompletion(ctx, tbContext.Chat().ID, text)
	if err != nil {
		return tbContext.Send("chat completion error")
	}
	return tbContext.Send(completion, menu)
}
