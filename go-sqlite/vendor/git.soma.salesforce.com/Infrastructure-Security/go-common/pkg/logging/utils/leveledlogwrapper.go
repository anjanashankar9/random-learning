package utils

import "github.com/sirupsen/logrus"

// Implements retryablehttp.LeveledLogger
type LeveledLogWrapper struct {
	wrapped *logrus.Entry
}

func NewLeveledLogWrapper(wrapped *logrus.Entry) *LeveledLogWrapper {
	return &LeveledLogWrapper{
		wrapped: wrapped,
	}
}

func (l *LeveledLogWrapper) Printf(msg string, keysAndValues ...interface{}) {
	l.Info(msg, keysAndValues...)
}

func (l *LeveledLogWrapper) Debug(msg string, keysAndValues ...interface{}) {
	l.ExtractKeyValues(keysAndValues...).Debug(msg)
}

func (l *LeveledLogWrapper) Info(msg string, keysAndValues ...interface{}) {
	l.ExtractKeyValues(keysAndValues...).Info(msg)
}

func (l *LeveledLogWrapper) Warn(msg string, keysAndValues ...interface{}) {
	l.ExtractKeyValues(keysAndValues...).Warn(msg)
}

func (l *LeveledLogWrapper) Error(msg string, keysAndValues ...interface{}) {
	l.ExtractKeyValues(keysAndValues...).Error(msg)
}

func (l *LeveledLogWrapper) ExtractKeyValues(keysAndValues ...interface{}) *logrus.Entry {
	result := l.wrapped
	for i := 0; i < len(keysAndValues)-1; i += 2 {
		k := keysAndValues[i]
		v := keysAndValues[i+1]
		result = result.WithField(k.(string), v)
	}

	return result
}
