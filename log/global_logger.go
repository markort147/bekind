package log

import (
	"io"
	"os"

	"github.com/bekind/bekindfrontend/config"
	"github.com/labstack/gommon/log"
)

/*
=== GLOBAL LOGGER CONFIGURATION ===
This file is used to configure the global logger for the application.
The global logger is used to log messages that are not specific to a particular package.
==================================
*/

var Logger = log.New("global")

func Init() {
	Logger.SetLevel(parseLevel(config.GetConfig().Log.Level))
	Logger.SetOutput(parseOutput(config.GetConfig().Log.Output))
}

func Test() {
	Logger.SetLevel(log.DEBUG)
	Logger.SetOutput(os.Stdout)
}

func parseLevel(level string) log.Lvl {
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

func parseOutput(output string) io.Writer {
	switch output {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		panic("invalid log output")
	}
}
