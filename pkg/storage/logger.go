package storage

import (
	"github.com/rs/zerolog"
	xormlog "xorm.io/xorm/log"
)

// implement xorm Logger interface
type xormLogger struct {
	logger  *zerolog.Logger
	showSql bool
}

func (l *xormLogger) Debug(v ...interface{}) {
	l.logger.Debug().Msgf("%v", v...)
}

func (l *xormLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

func (l *xormLogger) Error(v ...interface{}) {
	l.logger.Error().Msgf("%v", v...)
}

func (l *xormLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}

func (l *xormLogger) Info(v ...interface{}) {
	l.logger.Info().Msgf("%v", v...)
}

func (l *xormLogger) Infof(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}

func (l *xormLogger) Warn(v ...interface{}) {
	l.logger.Warn().Msgf("%v", v...)
}

func (l *xormLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Msgf(format, v...)
}

func (l *xormLogger) Level() xormlog.LogLevel {
	switch l.logger.GetLevel() {
	case zerolog.DebugLevel, zerolog.TraceLevel:
		return xormlog.LOG_DEBUG
	case zerolog.InfoLevel:
		return xormlog.LOG_INFO
	case zerolog.WarnLevel:
		return xormlog.LOG_WARNING
	case zerolog.ErrorLevel:
		return xormlog.LOG_ERR
	case zerolog.Disabled, zerolog.FatalLevel, zerolog.PanicLevel:
		return xormlog.LOG_OFF
	default:
		return xormlog.LOG_UNKNOWN
	}
}

func (l *xormLogger) SetLevel(level xormlog.LogLevel) {
	switch level {
	case xormlog.LOG_DEBUG:
		l.logger.Level(zerolog.DebugLevel)
	case xormlog.LOG_INFO:
		l.logger.Level(zerolog.InfoLevel)
	case xormlog.LOG_WARNING:
		l.logger.Level(zerolog.WarnLevel)
	case xormlog.LOG_ERR:
		l.logger.Level(zerolog.ErrorLevel)
	case xormlog.LOG_OFF:
		l.logger.Level(zerolog.Disabled)
	default:
		l.logger.Level(zerolog.NoLevel)
	}
}

func (l *xormLogger) ShowSQL(show ...bool) {
	l.showSql = show[0]
}

func (l *xormLogger) IsShowSQL() bool {
	return l.showSql
}
