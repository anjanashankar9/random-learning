package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Provider interface {
	GetServiceName() string
	GetHostname() string
	GetConfigPath() string
	GetDatacenter() string
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetInterface(key string, result interface{}) error
}

type provider struct {
	viper           *viper.Viper
	serviceName     string
	hostname        string
	configPath      string
	environmentKeys []string
}

func New(serviceName, hostname, configPath string, environmentKeys ...string) (Provider, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &provider{
		viper:           v,
		serviceName:     serviceName,
		hostname:        hostname,
		configPath:      configPath,
		environmentKeys: expandKeys(environmentKeys),
	}, nil
}

func (c *provider) GetServiceName() string {
	return c.serviceName
}

func (c *provider) GetHostname() string {
	return c.hostname
}

func (c *provider) GetConfigPath() string {
	return c.configPath
}

func (c *provider) GetDatacenter() string {
	if len(c.environmentKeys) == 0 {
		return ""
	}
	return c.environmentKeys[0]
}

func (c *provider) GetString(key string) string {
	for _, qualifiedKey := range c.qualifiedKeys(key) {
		if c.viper.IsSet(qualifiedKey) {
			return c.viper.GetString(qualifiedKey)
		}
	}
	return c.viper.GetString(key)
}

func (c *provider) GetInt(key string) int {
	for _, qualifiedKey := range c.qualifiedKeys(key) {
		if c.viper.IsSet(qualifiedKey) {
			return c.viper.GetInt(qualifiedKey)
		}
	}
	return c.viper.GetInt(key)
}

func (c *provider) GetBool(key string) bool {
	for _, qualifiedKey := range c.qualifiedKeys(key) {
		if c.viper.IsSet(qualifiedKey) {
			return c.viper.GetBool(qualifiedKey)
		}
	}
	return c.viper.GetBool(key)
}

func (c *provider) GetInterface(key string, result interface{}) error {
	for _, qualifiedKey := range c.qualifiedKeys(key) {
		if c.viper.IsSet(qualifiedKey) {
			return c.viper.UnmarshalKey(qualifiedKey, result)
		}
	}
	return c.viper.UnmarshalKey(key, result)
}

func (c *provider) qualifiedKeys(key string) []string {
	var qualifiedKeys []string
	for i := len(c.environmentKeys) - 1; i >= 0; i-- {
		qualifiedKeys = append(qualifiedKeys, fmt.Sprintf("%s.%s", c.environmentKeys[i], key))
	}

	return qualifiedKeys
}

func expandKeys(keys []string) []string {
	var expandedKeys []string

	var sb strings.Builder
	for _, key := range keys {
		sb.WriteString(key)
		expandedKeys = append(expandedKeys, sb.String())
		sb.WriteString(".")
	}

	return expandedKeys
}
