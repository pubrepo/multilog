package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

var (
	defaultLogFile string        = "./logs/log_info_%Y%m%d.log"
	defaultMaxDay  time.Duration = 24 * 30
	defaultChannel string        = "_default"
	defaultTimeFormat string ="2006-01-02T15:04:05.000"
	logger         *Logging
)

type Logging struct {
	loggers map[string]*zerolog.Logger
}

func init() {
	zerolog.TimeFieldFormat = defaultTimeFormat
	logger = &Logging{loggers: make(map[string]*zerolog.Logger)}
	AddChannel(defaultChannel)
}

func AddChannel(channel string) *zerolog.Logger {
	logFile := defaultLogFile

	if channel != defaultChannel {
		logFile = "./logs/" + channel
	}

	writer, _ := rotatelogs.New(logFile, rotatelogs.WithMaxAge(time.Hour*defaultMaxDay))

	output := zerolog.ConsoleWriter{Out: writer, TimeFormat: defaultTimeFormat,NoColor: true}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-2s|", i.(string)[0:1]))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	output.FormatCaller = func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			idx := strings.LastIndexByte(c, '/')
			if idx == -1 {
				return fmt.Sprintf(" %-20s |", c[idx+1:])
			}
			return fmt.Sprintf(" %-20s |", c[idx+1:])
		}
		return  fmt.Sprintf(" %-20s |", c)
	}

	l := zerolog.New(output).With().Timestamp().Logger()

	logger.loggers[channel] = &l

	return &l
}

func Info(args ...interface{}) {
	logger.loggers[defaultChannel].Info().Caller(1).Msg(fmt.Sprint(args...))
}

func Infof(channel string, args ...interface{}) {
	if channel == "" {
		logger.loggers[defaultChannel].Info().Caller(1).Msg(fmt.Sprint(args...))
		return
	}
	if logger, ok := logger.loggers[channel]; ok {
		logger.Info().Caller(1).Msg(fmt.Sprint(args...))
		return
	}
	AddChannel(channel).Info().Caller(1).Msg(fmt.Sprint(args...))
}

func Error(args ...interface{}) {
	logger.loggers[defaultChannel].Error().Caller(1).Msg(fmt.Sprint(args...))
}

func Errorf(channel string, args ...interface{}) {
	if channel == "" {
		logger.loggers[defaultChannel].Error().Msg(fmt.Sprint(args...))
		return
	}
	if logger, ok := logger.loggers[channel]; ok {
		logger.Error().Caller(1).Msg(fmt.Sprint(args...))
		return
	}
	AddChannel(channel).Error().Caller(1).Msg(fmt.Sprint(args...))
}
