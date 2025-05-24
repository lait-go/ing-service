package main

import (
	conf "tg_res/config"
	"tg_res/db"
	bot_tg "tg_res/internal/bot"

	// Error "tg_res/internal/err"
	LogWork "tg_res/log"
)

func main() {
	conf.Cfg = conf.Config_work()
	logger := LogWork.LogInit()
	logger.Info("логер и конфигурации созданы")

	

	db.ProcessingDB()

	defer db.Db.Close()

	bot_tg.TgInit()

	// keyboard := tgbotapi.NewReplyKeyboard(
	// 	tgbotapi.NewKeyboardButtonRow(
	// 		tgbotapi.NewKeyboardButton("Добавить резюме"),
	// 		tgbotapi.NewKeyboardButton("Поиск"),
	// 	),
	// )

	// for update := range updates {
	// 	if update.Message == nil {
	// 		continue
	// 	}

	// 	text := map[string]string{
	// 		"Добавить резюме": bot_tg.RegistationUser(),
	// 		"Поиск":           "Введите критерии для поиска.",
	// 	}[update.Message.Text]

	// 	if text == "" {
	// 		text = "Выберите действие с помощью кнопок ниже."
	// 	}

	// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	// 	msg.ReplyMarkup = keyboard
	// 	bot.Send(msg)
	// }
}


