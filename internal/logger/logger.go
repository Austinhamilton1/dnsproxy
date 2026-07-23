package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var Log *log.Logger

type Level int

const (
	ERROR Level = iota
	WARNING
	INFO
	DEBUG
)

var SetLevel Level = 0

func Init() {
	Log = log.New(
		os.Stdout,
		"",
		log.LstdFlags|log.Lmicroseconds,
	)
}

func writeArgs(level Level, args ...any) {
	if level <= SetLevel {
		var sb strings.Builder

		switch level {
		case ERROR:
			sb.WriteString("[ERROR]")
		case WARNING:
			sb.WriteString("[WARNING]")
		case INFO:
			sb.WriteString("[INFO]")
		case DEBUG:
			sb.WriteString("[DEBUG]")
		default:
			sb.WriteString("[UNKNOWN]")
		}

		for _, arg := range args {
			fmt.Fprintf(&sb, " %v", arg)
		}

		Log.Println(sb.String())
	}
}

func Error(args ...any) {
	writeArgs(ERROR, args)
	os.Exit(1)
}

func Warning(args ...any) {
	writeArgs(WARNING, args)
}

func Info(args ...any) {
	writeArgs(INFO, args)
}

func Debug(args ...any) {
	writeArgs(DEBUG, args)
}
