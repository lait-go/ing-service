package bot_tg

import (
	conf "tg_res/config"
	Error "tg_res/internal/err"
	LogWork "tg_res/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func TgInit() {
	logger := LogWork.LogInit()

	var err error
	Bot, err = tgbotapi.NewBotAPI(conf.Token)
	if Error.GetErr(err) {
		logger.Fatal("ошибка токена")
	}

	keyboard := keyInit()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		id := update.Message.Chat.ID

		if _, exists := userSteps[id]; !exists {
			msg := tgbotapi.NewMessage(id, "")
			msg.ReplyMarkup = keyboard
			Bot.Send(msg)
		}

		HandleUpdate(update)
		// distrebute(update)
	}
}


func distrebute(update tgbotapi.Update){
	switch update.Message.Text{
		case "Добавить резюме":
			HandleUpdate(update)
	}
} 

func keyInit()tgbotapi.ReplyKeyboardMarkup{
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить резюме"),
			tgbotapi.NewKeyboardButton("Поиск"),
		),
	)
}