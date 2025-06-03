package bot_tg

import (
	"fmt"
	"reflect"
	"tg_res/db"
	Error "tg_res/internal/err"

	// Error "tg_res/internal/err"
	// LogWork "tg_res/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct{
	Id   int     `db:"id"`
	Prof string  `db:"profession"`
	Name string  `db:"name"`
	Num  string  `db:"phone"`
	Tg   string  `db:"tg"`
}

var userSteps = make(map[int64]string)       // шаг регистрации по chatID
var userData = make(map[int64]*User)         // данные пользователя по chatID

func HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	id := update.Message.Chat.ID
	text := update.Message.Text

	if text == "/start" {
		msg := tgbotapi.NewMessage(id, "Добро пожаловать! Выберите действие:")
		msg.ReplyMarkup = keyInit() 
		Bot.Send(msg)
	}

	if text == "Назад" {
		backCheck(update, id)
	}

	switch userSteps[id] {
	case "":
		switch text{
		case "Добавить резюме":
			userData[id] = &User{}
			userSteps[id] = "name"
			msg := tgbotapi.NewMessage(id, "Введите ваше имя:")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			msg.ReplyMarkup = backInit()
			Bot.Send(msg)
			
		case "Поиск":
			userSteps[id] = "search"
			msg := tgbotapi.NewMessage(id, "Кого ищем?")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			msg.ReplyMarkup = backInit()
			Bot.Send(msg)
		default:
		}
			
		
	case "name":
			userData[id].Name = text
			userSteps[id] = "phone"
			Bot.Send(tgbotapi.NewMessage(id, "Введите ваш номер телефона:"))

	case "phone":
			userData[id].Num = text
			userSteps[id] = "prof"
			Bot.Send(tgbotapi.NewMessage(id, "Какие услуги вы предоставляете?"))

	case "prof":
		user := userData[id]
		user.Prof = text
		user.Tg = update.Message.From.UserName
	
		UserAdded(id, user)

	case "search":
		user := searchUser(text)

		val := reflect.ValueOf(user)
		typ := reflect.TypeOf(user)
	
		for i := 1; i < val.NumField(); i++ {
			field := typ.Field(i)
			value := val.Field(i)
			Bot.Send(tgbotapi.NewMessage(id, fmt.Sprintf("%s = %v", field.Name, value.Interface())))
		}
		
		delete(userSteps, id)
	}

}

func UserAdded(id int64, user *User) {
	query, err := db.Used_sql_with_parms("../db/migrations/add_person.sql")
	if err != nil {
		Bot.Send(tgbotapi.NewMessage(id, "Ошибка при подготовке SQL запроса"))
	} else {
		_, err = db.Db.Exec(query, user.Prof, user.Name, user.Num, user.Tg)
		if err != nil {
			Bot.Send(tgbotapi.NewMessage(id, "Ошибка при записи в базу данных"))
		} else {
			Bot.Send(tgbotapi.NewMessage(id, "Вы успешно зарегистрированы!"))
		}

		delete(userSteps, id)
		delete(userData, id)
	}
}

func searchUser(proff string) User{
	var user User

	date, err := db.Used_sql_with_parms("../db/migrations/serch_user.sql")
	Error.GetErr(err)

	rows, err := db.Db.Query(date, proff)
	Error.GetErr(err)

	for rows.Next(){
		rows.Scan(&user.Id, &user.Prof, &user.Name, &user.Num, &user.Tg)
	}

	return user
}

