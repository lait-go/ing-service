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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	for update := range updates {
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

func backInit()tgbotapi.ReplyKeyboardMarkup{
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
}

func backCheck(update tgbotapi.Update, id int64){
		userData[id] = &User{}
		userSteps[id] = ""

		msg := tgbotapi.NewMessage(update.Message.From.ID, "Выберите действие:") 
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = keyInit()
		Bot.Send(msg)
}