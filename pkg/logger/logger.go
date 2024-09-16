package logger

import "golang.org/x/exp/slog"

type ILogger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func NewLogger(level string) ILogger {
	if level == "" {
		level = envLocal
	}

	log := loggerWrapper{
		slog: newSlogLogger(level),
	}

	return &log
}

func (l *loggerWrapper) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}

func (l *loggerWrapper) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

func (l *loggerWrapper) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

func (l *loggerWrapper) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}

func Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

func String(key string, value string) slog.Attr {
	return slog.String(key, value)
}

func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

func Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}

func Error(err error) slog.Attr {
	return slog.Attr{
		Key:   "Error",
		Value: slog.StringValue(err.Error()),
	}
}

func With(l ILogger, args ...any) ILogger {
	switch v := l.(type) {
	case *loggerWrapper:
		return &loggerWrapper{
			slog: v.slog.With(args...),
		}
	default:
		return l
	}
}
