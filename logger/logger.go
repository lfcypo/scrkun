package logger

import (
	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

func init() {
	//color.ForceOpenColor()
}

func DisableColor() {
	color.Disable()
}

func New(name string) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(Output())
	logger.SetFormatter(NewFormatter(name))

	logger.SetLevel(logrus.DebugLevel)

	return logger
}
