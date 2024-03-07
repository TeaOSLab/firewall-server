package remotelogs

import (
	"github.com/iwind/TeaGo/logs"
)

// Debug 打印调试信息
func Debug(tag string, description string) {
	logs.Println("[" + tag + "]" + description)
}

// Println 打印普通信息
func Println(tag string, description string) {
	logs.Println("[" + tag + "]" + description)
}

// Warn 打印警告信息
func Warn(tag string, description string) {
	logs.Println("[" + tag + "]" + description)
}

// Error 打印错误信息
func Error(tag string, description string) {
	logs.Println("[" + tag + "]" + description)
}

// ErrorObject 打印错误对象
func ErrorObject(tag string, err error) {
	if err == nil {
		return
	}
	Error(tag, err.Error())
}
