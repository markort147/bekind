package log

import (
	"github.com/labstack/gommon/log"
)

var Logger = log.New("global")

func Init() {
	Logger.SetLevel(log.DEBUG)
}