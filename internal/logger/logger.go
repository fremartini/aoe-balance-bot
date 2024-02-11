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
	printf("INFO", format, a...)
}

func (*Logger) Info(s string) {
	print("INFO", s)
}

func (*Logger) Warnf(format string, a ...any) {
	printf("WARN", format, a...)
}

func (*Logger) Warn(s string) {
	print("WARN", s)
}

func (*Logger) Fatalf(format string, a ...any) {
	printf("FATAL", format, a...)
}

func (*Logger) Fatal(s string) {
	print("FATAL", s)
}

func printf(prefix, format string, a ...any) {
	fmt.Printf(fmt.Sprintf("[%s\t%s]\t%s\n", timestamp(), prefix, format), a...)
}

func print(prefix, s string) {
	fmt.Printf("[%s\t%s]\t%s\n", timestamp(), prefix, s)
}

func timestamp() string {
	return time.Now().Format(time.RFC850)
}
