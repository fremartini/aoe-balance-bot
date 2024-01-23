package logger

import (
	"fmt"
	"time"
)

type Logger struct {
	level uint
}

func New(level uint) *Logger {
	return &Logger{
		level: level,
	}
}

func (*Logger) Infof(format string, a ...any) {
	printf(format, a)
}

func (*Logger) Info(s string) {
	print(s)
}

func (*Logger) Fatalf(format string, a ...any) {
	printf(format, a)
}

func (*Logger) Fatal(s string) {
	print(s)
}

func printf(format string, a ...any) {
	fmt.Printf(fmt.Sprintf("[%s]\t%s\n", timestamp(), format), a...)
}

func print(s string) {
	fmt.Printf("[%s]\t%s\n", timestamp(), s)
}

func timestamp() string {
	return time.Now().Format(time.RFC850)
}
