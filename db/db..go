package db

import (
	"database/sql"
	"io"
	conf "tg_res/config"
	Error "tg_res/internal/err"
	"tg_res/internal/utils"
	LogWork "tg_res/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL драйвер
)

var Db *sqlx.DB

func ProcessingDB() {
	logger := LogWork.LogInit()
	var err error

	Db, err = sqlx.Connect("postgres", conf.Cfg.DbToken)
	if Error.GetErr(err){
		logger.Fatal("ошибка при чтении бд, проверьте ключ")
	}
	
	date, _ := Used_sql_not_parms("../db/migrations/check_table.sql")

	if date == nil {
		logger.Debug("таблица не существует")
		logger.Info("создание таблицы")

		_, err := Used_sql_not_parms("../db/migrations/create_table.sql")
		if Error.GetErr(err){
			logger.Fatal("ошибка при создании таблицы")
		}

		logger.Info("таблица создана")
	} 
		
	logger.Info("база данных в норме")
}


func Used_sql_not_parms(path string) (sql.Result, error){
	logger := LogWork.LogInit()

	file,err := utils.FileEx(path)
	if Error.GetErr(err){
		logger.Debug("ошибка при чтении файла")
		return nil, err
	}

	defer file.Close()

	body, err := io.ReadAll(file)
	if Error.GetErr(err){
		logger.Debug("ошибка при расшифровки файла")
		return nil, err
	}

	date := string(body)

	res, err := Db.Exec(date)
	if Error.GetErr(err){
		logger.Debug("ошибка при отправке запроса")
		return nil, err
	}

	return res, nil
}

func Used_sql_with_parms(path string) (string, error){
	logger := LogWork.LogInit()

	file,err := utils.FileEx(path)
	if Error.GetErr(err){
		logger.Debug("ошибка при чтении файла")
		return "", err
	}

	defer file.Close()

	body, err := io.ReadAll(file)
	if Error.GetErr(err){
		logger.Debug("ошибка при расшифровки файла")
		return "", err
	}

	return string(body), nil
}

