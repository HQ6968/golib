package logger

import (
	"github.com/verystar/golib/date"
	"time"
	"os"
)

type Config struct {
	LogName     string
	LogMode     string
	LogLevel    string
	LogMaxFiles int
	LogPath     string
}

var (
	// std is the name of the standard logger in stdlib `log`
	std           ILogger
	defaultConfig *Config
)

func init() {
	defaultConfig = &Config{
		LogName:  "app",
		LogMode:  "std",
		LogLevel: "info",
	}
	std = NewLogger()
}

func DefaultLogger(options ...func(*Config)) {
	std = NewLogger(options...)
}

func NewLogger(options ...func(*Config)) ILogger {
	conf := *defaultConfig

	for _, option := range options {
		option(&conf)
	}

	var log ILogger

	if conf.LogMode == "file" {
		d := date.Format("yyyy-MM-dd", time.Now())
		if conf.LogMaxFiles > 0 {
			delDate := date.Format("yyyy-MM-dd", time.Now().AddDate(0, 0, -conf.LogMaxFiles))
			os.Remove(conf.LogPath + "logs/" + conf.LogName + "-" + delDate + ".log")
		}

		log, _ = NewFileLogger(conf.LogPath+"logs/"+conf.LogName+"-"+d+".log", conf.LogLevel)
	} else {
		log = NewStdLogger()
	}
	return log
}

func Debug(str string, args ...interface{}) {
	std.Debug(str, args...)
}

func Info(str string, args ...interface{}) {
	std.Info(str, args...)
}

func Error(str string, args ...interface{}) {
	std.Error(str, args...)
}

func Fatal(str string, args ...interface{}) {
	std.Fatal(str, args...)
}
