package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/sirupsen/logrus"
	"github.com/evalphobia/logrus_sentry"
	"time"
)

var _ ILogger = (*FileLogger)(nil)
// FileLogger file logger
type FileLogger struct {
	*logrus.Logger
}

// NewFileLogger providers a file logger based on logrus
func NewFileLogger(filename string, logLevel string) (ILogger, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("can't get file abs path: filename = %v, err = %v", filename, err)
	}

	dirPath := filepath.Dir(absPath)
	if _, err := os.Stat(dirPath); err != nil {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("can't mkdirall directory: path = %v, err = %v", absPath, err)
		}
	}

	f, err := os.OpenFile(absPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("can't open file: path = %v, err = %v", absPath, err)
	}

	l := &logrus.Logger{
		Out:       f,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
	}

	switch logLevel {
	case LevelDebug:
		l.Level = logrus.DebugLevel
	case LevelInfo:
		l.Level = logrus.InfoLevel
	case LevelError:
		l.Level = logrus.ErrorLevel
	case LevelFatal:
		l.Level = logrus.FatalLevel
	}

	tags := map[string]string{
		"type": "go.job",
	}

	hook, err := logrus_sentry.NewWithTagsSentryHook(SENTRY_DSN, tags, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.InfoLevel,
	})
	hook.Timeout = 1 * time.Second
	hook.StacktraceConfiguration.Enable = true

	if err == nil {
		l.Hooks.Add(hook)
	}

	return &FileLogger{
		l,
	}, nil
}

func (l *FileLogger) Debug(format string, args ...interface{}) {
	l.WithFields(logrus.Fields{
		"fingerprint": []string{format},
	}).Debugf(format, args...)
}

func (l *FileLogger) Info(format string, args ...interface{}) {
	l.WithFields(logrus.Fields{
		"fingerprint": []string{format},
	}).Infof(format, args...)
}

func (l *FileLogger) Error(format string, args ...interface{}) {
	l.WithFields(logrus.Fields{
		"fingerprint": []string{format},
	}).Errorf(format, args...)
}

func (l *FileLogger) Fatal(format string, args ...interface{}) {
	l.WithFields(logrus.Fields{
		"fingerprint": []string{format},
	}).Fatalf(format, args...)
}

func (l *FileLogger) SetFormatter(format logrus.Formatter) {
	l.Formatter = format
}
