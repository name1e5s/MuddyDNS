package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// Default log format will output [INFO]: 2019-06-04T15:04:05+08:00 - Log message
	defaultLogFormat       = "[%lvl%]: %time% - %msg%\n"
	defaultTimestampFormat = time.RFC3339
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

// Formatter implements logrus.Formatter interface.
type Formatter struct {
	TimestampFormat string
	// Available standard keys: time, msg, lvl
	// Also can include custom fields but limited to strings.
	// All of fields need to be wrapped inside %% i.e %time% %msg%
	LogFormat string
}

var levelMap = map[string]string{
	"TRACE":   "TRAC",
	"DEBUG":   "DEBG",
	"INFO":    "INFO",
	"WARNING": "WARN",
	"ERROR":   "EROR",
	"FATAL":   "FATA",
	"PANIC":   "PNIC",
}

// Format log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	output = strings.Replace(output, "%msg%", entry.Message, 1)

	level := fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, levelMap[strings.ToUpper(entry.Level.String())])
	output = strings.Replace(output, "%lvl%", level, 1)
	for k, v := range entry.Data {
		if s, ok := v.(string); ok {
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}

	return []byte(output), nil
}
