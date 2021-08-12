package logger

import (
	"errors"
	"io"
	"io/ioutil"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/cxweoth/gin-api-server-template/internal/conf"
)

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterHook struct {
	Writer    io.Writer
	LogLevels []log.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []log.Level {
	return hook.LogLevels
}

// A function to make logger
func MakeLogger(cfg conf.IConf) (*logrus.Entry, error) {

	// Init logger
	var logger *logrus.Entry

	// Fetch config of logger
	loggerConf := cfg.LoggerCfg()

	apiServiceName := loggerConf.APIServiceName
	infoDebugLogPath := loggerConf.InfoDebugLogPath
	warnPanicLogPath := loggerConf.WarnPanicLogPath

	// Set to keep 5 min log and 60 mins history log for info and debug msg
	logInfoDebugWriter, err := rotatelogs.New(
		infoDebugLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(infoDebugLogPath),
		rotatelogs.WithMaxAge(time.Duration(12)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(1)*time.Hour),
	)

	if err != nil {
		return nil, errors.New("Log rotatelogs init failed" + err.Error())
	}

	// Set to keep 1 hour log and 168 hours (a week) history log for warn and panic msg
	logWarnPanicWriter, err := rotatelogs.New(
		warnPanicLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(warnPanicLogPath),
		rotatelogs.WithMaxAge(time.Duration(168)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(1)*time.Hour),
	)

	if err != nil {
		return nil, errors.New("Log rotatelogs init failed" + err.Error())
	}

	// Set fomatter
	log.SetFormatter(&log.JSONFormatter{})

	// Send all logs to nowhere by default
	log.SetOutput(ioutil.Discard)

	// Adds hooks to send logs to different destinations on info and debug level
	log.AddHook(&WriterHook{
		Writer: logInfoDebugWriter,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
		},
	})

	// Adds hooks to send logs to different destinations on warn and panic level
	log.AddHook(&WriterHook{
		Writer: logWarnPanicWriter,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		},
	})

	// Set fields memo in each log
	logger = log.WithFields(log.Fields{"api_service_name": apiServiceName})

	return logger, nil
}
