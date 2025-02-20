package log

import (
	"io"
	"os"

	"github.com/labstack/gommon/log"
	"github.com/markort147/bekind/internal/config"
)

/*
=== GLOBAL LOGGER CONFIGURATION ===
This file is used to configure the global logger for the application.
The global logger is used to log messages that are not specific to a particular package.
==================================
*/

var (
	Logger  = log.New("global")
	logFile *os.File
)

func Init(cfg *config.Config) {
	Logger.SetLevel(ParseLevel(cfg.Log.Level))
	Logger.SetOutput(ParseOutput(cfg.Log.Output))
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func Test() {
	Logger.SetLevel(log.DEBUG)
	Logger.SetOutput(os.Stdout)
}

func ParseLevel(level string) log.Lvl {
	switch level {
	case "debug":
		return log.DEBUG
	case "info":
		return log.INFO
	case "warn":
		return log.WARN
	case "error":
		return log.ERROR
	case "off":
		return log.OFF
	default:
		panic("invalid log level")
	}
}

func ParseOutput(output string) io.Writer {
	switch output {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("failed to open log file: " + err.Error())
		}
		logFile = file
		return file
	}
}
