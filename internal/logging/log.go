package logging

import (
	"github.com/rs/zerolog"
	"io"
)

var Logger InternalLogger

func Crit(msg string, args ...interface{}) {
	if len(args) > 0 {
		Logger.Logger.Error().Msgf(msg, args...)
	} else {
		Logger.Logger.Error().Msg(msg)
	}
}

func Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		Logger.Logger.Warn().Msgf(msg, args...)
	} else {
		Logger.Logger.Warn().Msg(msg)
	}
}

func Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		Logger.Logger.Debug().Msgf(msg, args...)
	} else {
		Logger.Logger.Debug().Msg(msg)
	}
}

func Configure(dst io.Writer, level zerolog.Level) {
	Logger = InternalLogger{
		zerolog.New(dst).With().Timestamp().Logger(),
	}
	zerolog.SetGlobalLevel(level)
}

type InternalLogger struct {
	zerolog.Logger
}

func (l *InternalLogger) GetSubLogger(k string, v string) *InternalLogger {
	return &InternalLogger{
		Logger: l.Logger.With().Str(k, v).Logger(),
	}
}

func (l *InternalLogger) Crit(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.Logger.Error().Msgf(msg, args...)
	} else {
		l.Logger.Error().Msg(msg)
	}
}

func (l *InternalLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.Logger.Warn().Msgf(msg, args...)
	} else {
		l.Logger.Warn().Msg(msg)
	}
}

func (l *InternalLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.Logger.Debug().Msgf(msg, args...)
	} else {
		l.Logger.Debug().Msg(msg)
	}
}
