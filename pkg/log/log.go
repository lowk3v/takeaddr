package log

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
	symbol SymbolConfig
}

type SymbolConfig struct {
	Success string
	Error   string
	Info    string
	Debug   string
}

func New(level string) *Logger {
	l := &Logger{
		Logger: logrus.New(),
		symbol: SymbolConfig{
			Success: color.GreenString("≠"),
			Error:   color.RedString("¿"),
			Info:    color.BlueString("ℹ"),
			Debug:   color.MagentaString("☣"),
		},
	}
	l.SetReportCaller(false)
	l.SetLevel(level)
	l.Formatter = &logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	return l
}

func (l *Logger) SetLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	l.Logger.SetLevel(lvl)
}
