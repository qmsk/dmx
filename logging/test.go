package logging

import (
	"github.com/Sirupsen/logrus"
)

func New(pkg string) Logger {
	logger := logrus.New()
	logger.Level = log.Level
	logger.Formatter = log.Formatter

	return Context{logger.WithFields(logrus.Fields{"package": pkg})}
}
