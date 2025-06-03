package bot_tg

import (
	"fmt"
	"os/exec"
	"reflect"
	"tg_res/db"
	Error "tg_res/internal/err"
	"tg_res/internal/utils"

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
		
		backCheck(update, id)

	case "search":
		user := searchUser(text)

		for _, row := range user{
			val := reflect.ValueOf(row)
			typ := reflect.TypeOf(row)
		
			for i := 1; i < val.NumField(); i++ {
				field := typ.Field(i)
				value := val.Field(i)
				Bot.Send(tgbotapi.NewMessage(id, fmt.Sprintf("%s = %v", field.Name, value.Interface())))
			}
		}
		backCheck(update, id)
		
		delete(userSteps, id)

		
	}

}

func UserAdded(id int64, user *User) {
	query, err := db.Used_sql_with_parms("../db/migrations/add_person.sql")
	if err != nil {
		Bot.Send(tgbotapi.NewMessage(id, "Ошибка при подготовке SQL запроса"))
	} else {
		err := db.Db.QueryRow(query, user.Prof, user.Name, user.Num, user.Tg).Scan(&user.Id)
		if err != nil {
			Bot.Send(tgbotapi.NewMessage(id, "Ошибка при записи в базу данных"))
		} else {
			Bot.Send(tgbotapi.NewMessage(id, "Вы успешно зарегистрированы!"))

			cmd := exec.Command("python3", "../bert/bert.py", user.Prof, fmt.Sprint(user.Id))
			err := cmd.Run()
			Error.GetErr(err)
		}

		delete(userSteps, id)
		delete(userData, id)
	}
}

func searchUser(proff string) []User{
	var user []User

	res, err := exec.Command("python3", "../bert/vector_return.py", proff).Output()
	Error.GetErr(err)
	fmt.Println(string(res))

	floatRow, err := utils.ParsePgvectorString(string(res))
	Error.GetErr(err)

	stringRow := utils.FormatVectorForSQL(floatRow)

	date, err := db.Used_sql_with_parms("../db/migrations/serch_user.sql")
	Error.GetErr(err)

	rows, err := db.Db.Query(date, stringRow)
	Error.GetErr(err)

	defer rows.Close()

	for rows.Next() {
		var u User
		var distance float64
		err = rows.Scan(&u.Name, &distance, &u.Num, &u.Prof, &u.Tg)
		Error.GetErr(err)
		user = append(user, u)
	}

	return user
}

