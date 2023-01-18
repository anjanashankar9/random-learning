package logging

import (
	"fmt"
	"path"

	"github.com/sirupsen/logrus"

	"git.soma.salesforce.com/Infrastructure-Security/go-common/pkg/config"
)

const (
	KeyLogLevel      = "logging.level"
	KeyStdoutEnabled = "logging.stdout.enabled"
	KeyStdoutFormat  = "logging.stdout.format"
	KeyFileEnabled   = "logging.file.enabled"
	KeyFileFormat    = "logging.file.format"
	KeyFilePath      = "logging.file.path"

	FormatText = "text"
	FormatJson = "json"
)

type Config struct {
	Level         string
	StdoutEnabled bool
	StdoutFormat  string
	FileEnabled   bool
	FileFormat    string
	FileName      string
	FilePath      string
}

func ParseConfig(provider config.Provider) Config {
	cfg := Config{
		Level:         provider.GetString(KeyLogLevel),
		StdoutEnabled: provider.GetBool(KeyStdoutEnabled),
		FileEnabled:   provider.GetBool(KeyFileEnabled),
	}

	_, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		cfg.Level = logrus.InfoLevel.String()
	}

	if cfg.StdoutEnabled {
		cfg.StdoutFormat = provider.GetString(KeyStdoutFormat)
		if !isValidFormat(cfg.StdoutFormat) {
			cfg.StdoutFormat = FormatJson
		}
	}

	if cfg.FileEnabled {
		cfg.FileFormat = provider.GetString(KeyFileFormat)
		if !isValidFormat(cfg.FileFormat) {
			cfg.FileFormat = FormatJson
		}

		// This is useful for other configs that want to build a
		// path to another log file in the standard logging directory
		cfg.FilePath = provider.GetString(KeyFilePath)
		cfg.FileName = path.Join(cfg.FilePath, fmt.Sprintf("%s.log", provider.GetServiceName()))
	}

	return cfg
}

func isValidFormat(format string) bool {
	return format == FormatText || format == FormatJson
}
