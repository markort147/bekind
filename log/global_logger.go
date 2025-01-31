package log

import (
	"os"

	"github.com/labstack/gommon/log"
)

var Logger = log.New("global")

func Init() {
	Logger.SetLevel(log.DEBUG)
}

func Test() {
	Logger.SetLevel(log.DEBUG)
	Logger.SetOutput(os.Stdout)
}