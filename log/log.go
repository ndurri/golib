package log

import (
	"log"
)

func Error(e error) {
	log.Printf("ERROR: %v", e)
}

func Info(info string) {
	log.Println("INFO: " + info)
}

func Infofmt(fmt string, args ...any) {
	log.Printf("INFO: "+fmt, args...)
}

func Warn(warning string) {
	log.Println("WARNING: " + warning)
}
