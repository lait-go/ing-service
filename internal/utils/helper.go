package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func FileEx(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}

	return file, nil
}

func FileClear(file *os.File) bool {
	if _, seekErr := file.Seek(0, 0); seekErr != nil {
		log.Println("Ошибка seek файла:", seekErr)
		return true
	}

	if truncateErr := file.Truncate(0); truncateErr != nil {
		log.Println("Ошибка очистки файла:", truncateErr)
		return true
	}
	return false
}

func TimeFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func FormInit(date any) string {
	switch v := date.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return ""
	}
}

func ParsePgvectorString(s string) ([]float64, error) {
	s = strings.Trim(s, "[]")
	parts := strings.Split(s, ",")
	vec := make([]float64, 0, len(parts))

	for _, p := range parts {
		var f float64
		_, err := fmt.Sscanf(strings.TrimSpace(p), "%f", &f)
		if err != nil {
			return nil, fmt.Errorf("не удалось распарсить '%s': %v", p, err)
		}
		vec = append(vec, f)
	}
	return vec, nil
}

func FormatVectorForSQL(vec []float64) string {
	parts := make([]string, len(vec))
	for i, v := range vec {
		parts[i] = fmt.Sprintf("%f", v)
	}
	return fmt.Sprintf("ARRAY[%s]::vector", strings.Join(parts, ","))
}