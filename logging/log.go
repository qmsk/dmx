// Combine go-flags + logrus
package logging

import (
	"github.com/Sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(fmt string, args ...interface{})
	Info(args ...interface{})
	Infof(fmt string, args ...interface{})
	Warn(args ...interface{})
	Warnf(fmt string, args ...interface{})
	Error(args ...interface{})
	Errorf(fmt string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(fmt string, args ...interface{})

	Logger(field string, value interface{}) Logger
}

type Option struct {
	Package  string
	LogLevel logrus.Level
}

func (option *Option) UnmarshalFlag(value string) error {
	if logLevel, err := logrus.ParseLevel(value); err != nil {
		return err
	} else {
		option.LogLevel = logLevel
	}

	return nil
}

func (option Option) Logger(field string, value interface{}) Logger {
	logger := logrus.New()
	logger.Level = option.LogLevel

	return Context{logger.WithFields(logrus.Fields{"package": option.Package, field: value})}
}

type Context struct {
	*logrus.Entry
}

func (context Context) Logger(field string, value interface{}) Logger {
	return Context{context.WithField(field, value)}
}
