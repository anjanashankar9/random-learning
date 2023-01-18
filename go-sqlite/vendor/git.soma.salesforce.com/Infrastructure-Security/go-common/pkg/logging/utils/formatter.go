package utils

import (
	"fmt"
	"strings"

	"github.com/go-stack/stack"
	"github.com/sirupsen/logrus"
)

const defaultTimestampFormat = "2006-01-02T15:04:05.999Z07:00"

// Implements logrus.Formatter
type loggingFormatter struct {
	logrus.Formatter
}

func NewTextFormatter(disableColors bool) *loggingFormatter {
	return &loggingFormatter{
		Formatter: &logrus.TextFormatter{
			DisableColors:   disableColors,
			FullTimestamp:   true,
			TimestampFormat: defaultTimestampFormat,
		},
	}
}

func NewJsonFormatter() *loggingFormatter {
	return &loggingFormatter{
		Formatter: &logrus.JSONFormatter{
			DataKey:         "data",
			TimestampFormat: defaultTimestampFormat,
		},
	}
}

func (s *loggingFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if !shouldAddStackTrace(entry) {
		return s.Formatter.Format(entry)
	}

	withTrace := entry.WithField("trace", trimmedStackTrace())
	withTrace.Level = entry.Level
	withTrace.Message = entry.Message
	return s.Formatter.Format(withTrace)
}

func shouldAddStackTrace(entry *logrus.Entry) bool {
	return entry.Level == logrus.ErrorLevel ||
		entry.Level == logrus.PanicLevel ||
		entry.Level == logrus.FatalLevel
}

func trimmedStackTrace() string {
	callStack := stack.Trace()
	trimmedStack := make([]string, 0)
	for _, call := range callStack {
		filePath := fmt.Sprintf("%#s", call)
		// Strip out leading $GOPATH/src
		subPaths := strings.SplitAfter(filePath, "/src/")
		if len(subPaths) > 1 {
			filePath = strings.Join(subPaths[1:], "")
		}
		line := fmt.Sprintf("%s:%d %n()", filePath, call, call)
		// Don't include calls in the logging stack
		if !strings.Contains(line, "/sirupsen/logrus/") &&
			!strings.Contains(line, "/logging/") {
			trimmedStack = append(trimmedStack, line)
		}
	}
	return strings.Join(trimmedStack, "\n")
}
