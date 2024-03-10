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

func (l *Logger) Infof(format string, a ...any) {
	if l.level >= INFO {
		printf("INFO", format, a...)
	}
}

func (l *Logger) Info(s string) {
	if l.level >= INFO {
		print("INFO", s)
	}
}

func (l *Logger) Warnf(format string, a ...any) {
	if l.level >= WARN {
		printf("WARN", format, a...)
	}
}

func (l *Logger) Warn(s string) {
	if l.level >= WARN {
		print("WARN", s)
	}
}

func (l *Logger) Fatalf(format string, a ...any) {
	if l.level >= FATAL {
		printf("FATAL", format, a...)
	}
}

func (l *Logger) Fatal(s string) {
	if l.level >= FATAL {
		print("FATAL", s)
	}
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
