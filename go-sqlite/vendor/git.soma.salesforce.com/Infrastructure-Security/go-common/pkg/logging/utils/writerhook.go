package utils

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Implements logrus.Hook methods
type writerHook struct {
	formatter logrus.Formatter
	writer    io.Writer
}

func NewWriterHook(formatter logrus.Formatter, writer io.Writer) *writerHook {
	return &writerHook{
		formatter: formatter,
		writer:    writer,
	}
}

func (w *writerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (w *writerHook) Fire(e *logrus.Entry) error {
	bytes, err := w.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
