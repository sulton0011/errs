package errs

import (
	botV5 "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BroadcastBot struct handles sending messages to multiple Telegram chats
type BroadcastBot struct {
	bot     *botV5.BotAPI
	chatIDs []int64
}

// NewBroadcastBot creates a new instance of BroadcastBot
func NewBroadcastBot(token string, chatIDs []int64) error {
	if token == "" || len(chatIDs) == 0 {
		return New("Failed to create Telegram bot. Invalid token or chat ID.")
	}

	b, err := botV5.NewBotAPI(token)
	if err != nil {
		return Wrap(err, "failed to create telegram bot")
	}

	bot = &BroadcastBot{
		bot:     b,
		chatIDs: chatIDs,
	}

	return nil
}

// SendMessage sends a message to all configured chat IDs
func (bb *BroadcastBot) sendMessage(msg string) error {
	var errs error
	for _, chatID := range bb.chatIDs {
		if err := bb.sendToChat(chatID, msg); err != nil {
			errs = Join(" && ", errs, err)
		}
	}
	return errs
}

// sendToChat sends a message to a specific chat ID
func (bb *BroadcastBot) sendToChat(chatID int64, msg string) error {
	m := botV5.NewMessage(chatID, msg)
	m.ParseMode = "Markdown"
	_, err := bb.bot.Send(m)
	if err != nil {
		return WrapF(err, "failed to send message to chat %d", chatID)
	}
	return nil
}
