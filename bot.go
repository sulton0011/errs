package errs

import (
	botV5 "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BroadcastBot struct handles sending messages to multiple Telegram chats
type broadcastBot struct {
	bot     *botV5.BotAPI
	chatIDs []int64
	trimSpace bool
}

type BroadcastBotParams struct {
	ServiceName string
	Token       string
	ChatIDs     []int64
	TrimSpace   bool
}

// NewBroadcastBot creates a new instance of BroadcastBot
func NewBroadcastBot(params BroadcastBotParams) error {
	if params.Token == "" || len(params.ChatIDs) == 0 {
		return New("Failed to create Telegram bot. Invalid token or chat ID.")
	}

	sTitle = params.ServiceName

	b, err := botV5.NewBotAPI(params.Token)
	if err != nil {
		return Wrap(err, "failed to create telegram bot")
	}

	bot = &broadcastBot{
		bot:     b,
		chatIDs: params.ChatIDs,
		trimSpace: params.TrimSpace,
	}

	return nil
}

// SendMessage sends a message to all configured chat IDs
func (bb *broadcastBot) sendMessage(msg string) error {
	var errs error
	for _, chatID := range bb.chatIDs {
		if err := bb.sendToChat(chatID, msg); err != nil {
			errs = Join(" && ", errs, err)
		}
	}
	return errs
}

// sendToChat sends a message to a specific chat ID
func (bb *broadcastBot) sendToChat(chatID int64, msg string) error {
	m := botV5.NewMessage(chatID, msg)
	m.ParseMode = "Markdown"
	_, err := bb.bot.Send(m)
	if err != nil {
		return WrapF(err, "failed to send message to chat %d", chatID)
	}
	return nil
}
