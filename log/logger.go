package LogWork

import (
	conf "tg_res/config"
	Error "tg_res/internal/err"
	"tg_res/internal/utils"

	"encoding/json"
	"io"
	"log/slog"
	"os"
)

type LogStruct struct {
	Logger slog.Logger
}

type readLog struct {
	Log   string  `json:"log"`
	Time  string  `json:"time"`
	Level string  `json:"level"`
}

func (logStruct *LogStruct) Info(message interface{}) {
	res := utils.FormInit(message)
	logStruct.Logger.Info(res)

	pars := readLog{
		Log: res,
		Time: utils.TimeFormat(),
		Level: "INFO",
	}

	ReadLog(pars)
}

func (logStruct *LogStruct) Fatal(message interface{}){
	res := utils.FormInit(message)
	logStruct.Logger.Error(res)
	
	defer os.Exit(1)

	pars := readLog{
		Log: res,
		Time: utils.TimeFormat(),
		Level: "FATAL",
	}

	ReadLog(pars)
}

func (logStruct *LogStruct) Debug(message interface{}) {
	res := utils.FormInit(message)
	logStruct.Logger.Debug(res)

	pars := readLog{
		Log: res,
		Time: utils.TimeFormat(),
		Level: "DEBUG",
	}

	ReadLog(pars)
}

func (logStruct *LogStruct) Warn(message interface{}) {
	res := utils.FormInit(message)
	logStruct.Logger.Warn(res)

	pars := readLog{
		Log: res,
		Time: utils.TimeFormat(),
		Level: "WARN",
	}

	ReadLog(pars)
}

func (logStruct *LogStruct) Error(message interface{}) {
	res := utils.FormInit(message)
	logStruct.Logger.Error(res)

	pars := readLog{
		Log: res,
		Time: utils.TimeFormat(),
		Level: "ERROR",
	}

	ReadLog(pars)
}

func ReadLog(pars readLog) {
	file, err := utils.FileEx(conf.Cfg.Path.LogPath)
	if err != nil {
		file, err = os.Create("../internal/utils/log/log_storage.json")
		if err != nil {
			Error.GetErr(err)
		}
	}

	defer file.Close()

	var logs []readLog
	data, err := io.ReadAll(file)
	if err != nil {
		Error.GetErr(err)
	}

	if len(data) > 0 {
		if unmarshalErr := json.Unmarshal(data, &logs); unmarshalErr != nil {
			Error.GetErr(unmarshalErr)
		}
	}

	logs = append([]readLog{pars}, logs...)

	utils.FileClear(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(logs)
	if err != nil {
		Error.GetErr(err)
	}
}


func LogInit() *LogStruct{
	var level slog.Level

	switch conf.Cfg.Env {
	case "local":
		level = slog.LevelDebug
	case "dev":
		level = slog.LevelInfo
	case "staging":
		level = slog.LevelWarn
	case "prod":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var log LogStruct
	log.Logger = *slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	return &log
}
