package models

import (
	"strconv"
)

func Error(msg string, code int) map[string]string {
	return map[string]string{
		"message": msg,
		"code":    strconv.Itoa(code),
	}
}

func ErrorPOSTResponse(inputValue string) *Data {
	return &Data{
		Value:  inputValue,
		Result: "ERROR",
	}
}
