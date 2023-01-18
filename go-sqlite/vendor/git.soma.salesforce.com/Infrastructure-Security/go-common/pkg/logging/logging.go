package logging

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"git.soma.salesforce.com/Infrastructure-Security/go-common/pkg/logging/utils"
)

const (
	sourceKey = ".source" // Note . is a hack to make it the first output
)

var rootContainer LogContainer = &logContainer{
	rootLogger: logrus.New(), // Default config if logger is used without configuring (tests)
}

type LogContainer interface {
	New(logSource string) *logrus.Entry
	IsLevelEnabled(level logrus.Level) bool
	GetConfig() Config
}

type logContainer struct {
	rootLogger    *logrus.Logger
	loggingConfig Config
}

func Init(cfg Config) LogContainer {
	return innerInit(cfg)
}

func Configure(cfg Config) {
	lc := innerInit(cfg)

	// Set the default output for libraries that use logrus directly via logrus.Info
	logrus.SetFormatter(lc.rootLogger.Formatter)
	logrus.SetOutput(lc.rootLogger.Out)

	for _, hooks := range lc.rootLogger.Hooks {
		for _, hook := range hooks {
			logrus.AddHook(hook)
		}
	}

	rootContainer = lc
}

func (lc *logContainer) New(logSource string) *logrus.Entry {
	logger := lc.rootLogger.WithField(sourceKey, logSource)
	logger.Level = lc.rootLogger.Level
	return logger
}

func New(logSource string) *logrus.Entry {
	return rootContainer.New(logSource)
}

func (lc *logContainer) IsLevelEnabled(level logrus.Level) bool {
	return lc.rootLogger.IsLevelEnabled(level)
}

func IsLevelEnabled(level logrus.Level) bool {
	return rootContainer.IsLevelEnabled(level)
}

func (lc *logContainer) GetConfig() Config {
	return lc.loggingConfig
}

func GetConfig() Config {
	return rootContainer.GetConfig()
}

func innerInit(cfg Config) *logContainer {
	logger := logrus.New()
	logger.Level, _ = logrus.ParseLevel(cfg.Level)

	if cfg.StdoutEnabled {
		logger.Out = os.Stdout
		logger.Formatter = newFormatter(cfg.StdoutFormat, false)
	} else {
		logger.Out = ioutil.Discard
	}

	if cfg.FileEnabled {
		lumberjackWriterHook := utils.NewWriterHook(
			newFormatter(cfg.FileFormat, true), // Disable colors on the text formatter when writing to a file
			newLumberjackWriter(cfg.FileName))
		logger.AddHook(lumberjackWriterHook)
	}

	return &logContainer{
		rootLogger:    logger,
		loggingConfig: cfg,
	}
}

func newFormatter(format string, disableColors bool) logrus.Formatter {
	if format == FormatText {
		return utils.NewTextFormatter(disableColors)
	} else {
		return utils.NewJsonFormatter()
	}
}

func newLumberjackWriter(logFile string) *lumberjack.Logger {
	// Max hard disk space is ~1GB
	const maxSizeMB = 100
	const maxBackupFiles = 10
	const maxAgeDays = 30
	return &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackupFiles,
		MaxAge:     maxAgeDays,
		Compress:   true,
	}
}
