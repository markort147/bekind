package log

import (
	"os"

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
	Logger.SetLevel(log.DEBUG)
}

func Test() {
	Logger.SetLevel(log.DEBUG)
	Logger.SetOutput(os.Stdout)
}
