# logging

The `logging` package provides an opinionated logging framework. The logger can be configured to write to stdout, or a file, or both. If file logging is configured, the log files are automatically limited to 1GB in size, and rotated automatically to prevent unbounded growth.

Configuring the logger requires a logging configuration. You can either create the config manually:

```go
cfg := logging.Config {
    LogLevel:    "DEBUG",
    LogToStdout: true,
    LogToFile:   true,
    LogFile:     "/tmp/path/to/logs/myService.log" 
}
```

Or, using the common config provider:

Sample config:

```properties
logging.level=INFO
logging.stdout.enabled=false
logging.file.enabled=true
logging.file.path=/var/log/myService
# Note: You only need to provide the directory.
#   The log file will be constructed as: <logging.file.path>/<configProvider.GetServiceName()>.log 
```

```go
cfg := logging.ParseConfig(provider)
```

Configure the logging library:

```go
logging.Configure(cfg)
```

Create and use the logger:

```go
// The logSource will show up as a field in any logs emitted by the logger.
// It is just an opaque identifier to roughly locate the source of a log.
// The go file name is used by convention, but you can use any string.
logSource := "myService.go"
logger = logging.New(logSource)

logger.Info("message")
// Outputs: INFO[2021-02-10T10:23:47.553-06:00] message .source=myService.go

logger = logger.WithField("foo", bar)
logger.Warn("message")
// Outputs: WARN[2021-02-10T10:23:47.553-06:00] message .source=myService.go foo=bar
```