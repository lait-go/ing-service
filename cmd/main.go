package main

import (
	conf "tg_res/config"
	Error "tg_res/internal/err"
	LogWork "tg_res/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	conf.Cfg = conf.Config_work()
	logger := LogWork.LogInit()
	logger.Info("logger and configuration init")

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if Error.GetErr(err) {
		logger.Fatal(err)
	}




	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить резюме"),
			tgbotapi.NewKeyboardButton("Поиск"),
		),
	)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := map[string]string{
			"Добавить резюме": "Пожалуйста, отправьте ваше резюме.",
			"Поиск":           "Введите критерии для поиска.",
		}[update.Message.Text]

		if text == "" {
			text = "Выберите действие с помощью кнопок ниже."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}
}


