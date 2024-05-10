// Package logging is a wrapper around the selected logging package, which provides simple API and
// predefined configurations applied to the logger e.g. dev mode or production mode loggers
package logging

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// LogProp represents a key-value log properties to be provided by the user
type LogProp struct {
	Key   string
	Value string
}

var logger *zerolog.Logger
var initializeOnce sync.Once

func init() {
	initLogger()
}

func Trace(msg string, props []LogProp) {
	applyProps(logger.Trace(), props).Msg(msg)
}

func Debug(msg string, props []LogProp) {
	applyProps(logger.Debug(), props).Msg(msg)
}

func Info(msg string, props []LogProp) {
	applyProps(logger.Info(), props).Msg(msg)
}

func Warning(err error, msg string, props []LogProp) {
	if err != nil {
		if props == nil {
			props = make([]LogProp, 1)
		}
		props = append(props, LogProp{Key: zerolog.ErrorFieldName, Value: err.Error()})
	}
	applyProps(logger.Warn(), props).Msg(msg)
}

func Error(err error, msg string, props []LogProp) {
	if err == nil {
		//Applying logger.Err with a nil error will apply INFO severity
		applyProps(logger.Error(), props).Msg(msg)
		return
	}
	applyProps(logger.Err(err), props).Msg(msg)
}

func Panic(err error, msg string, props []LogProp) {
	if err != nil {
		if props == nil {
			props = make([]LogProp, 1)
		}
		props = append(props, LogProp{Key: zerolog.ErrorFieldName, Value: err.Error()})
	}
	applyProps(logger.Panic(), props).Msg(msg)
}

func applyProps(logEvent *zerolog.Event, props []LogProp) *zerolog.Event {
	if props == nil {
		return logEvent
	}
	for _, prop := range props {
		logEvent.Str(prop.Key, prop.Value)
	}
	return logEvent
}

func initLogger() {
	initializeOnce.Do(func() {
		instance := initializeDevModeLogger()
		logger = &instance
	})
}

func initializeDevModeLogger() zerolog.Logger {
	instance := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel).Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
			FormatTimestamp: func(t interface{}) string {
				if strTime, ok := t.(string); ok {
					parse, _ := time.Parse(time.RFC3339, strTime)
					return parse.Format("2006-01-02 15:04:05.1234")
				}
				return fmt.Sprint(t)
			},
			FormatLevel: func(i interface{}) string {
				if l, ok := i.(string); ok {
					return strings.ToUpper(l)
				}
				return ""
			},
			FormatMessage: func(i interface{}) string {
				return i.(string)
			},
			//Do not include custom log properties
			FormatFieldName: func(i interface{}) string {
				return ""
			},
			FormatFieldValue: func(i interface{}) string {
				return ""
			},
		})
	return instance
}
