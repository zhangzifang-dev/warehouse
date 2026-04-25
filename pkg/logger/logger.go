package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	infoLog = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	errLog  = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
)

func Info(format string, args ...interface{}) {
	infoLog.Println(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	errLog.Println(fmt.Sprintf(format, args...))
}
