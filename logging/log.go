// Combine go-flags + logrus
package logging

import (
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
	"strings"
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

type Context struct {
	*logrus.Entry
}

func (context Context) Logger(field string, value interface{}) Logger {
	return Context{context.WithFields(logrus.Fields{"type": field, field: value})}
}

type Formatter struct {
}

func (formatter Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer
	var entryType string

	fmt.Fprintf(&buf, "%-5s ", strings.ToUpper(entry.Level.String()))

	if entryPackage := entry.Data["package"]; entryPackage != nil {
		fmt.Fprintf(&buf, "%20s:", entryPackage)
	} else {
		fmt.Fprintf(&buf, "%20s ", "")
	}

	if entryTypeInterface := entry.Data["type"]; entryTypeInterface == nil {

	} else if entryTypeString, ok := entryTypeInterface.(string); !ok {

	} else {
		entryType = entryTypeString
	}

	if entryType == "" {
		fmt.Fprintf(&buf, "%-30s ", "")
	} else if typeValue := entry.Data[entryType]; typeValue == nil {
		fmt.Fprintf(&buf, "%-30s ", entryType)
	} else {
		entry := fmt.Sprintf("%s<%s>", entryType, typeValue)

		fmt.Fprintf(&buf, "%-30s ", entry)
	}

	buf.WriteString(entry.Message)

	for key, value := range entry.Data {
		if key == "package" || key == "type" || key == entryType {
			continue
		}

		fmt.Fprintf(&buf, ": %s=%s", key, value)
	}

	buf.WriteString("\n")

	return buf.Bytes(), nil
}

var log *logrus.Logger
var Log Logger

func init() {
	log = logrus.StandardLogger()
	log.Formatter = Formatter{}

	Log = Context{log.WithFields(logrus.Fields{})}
}

func Setup(option Option) {
	log.Level = option.LogLevel
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
	logger.Level = log.Level
	logger.Formatter = log.Formatter

	if option.LogLevel != 0 {
		logger.Level = option.LogLevel
	}

	return Context{logger.WithFields(logrus.Fields{"package": option.Package, "type": field, field: value})}
}
