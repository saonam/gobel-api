package infrastructure

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bmf-san/gobel-api/app/usecases"
	"github.com/pkg/errors"
)

// Level is log levels.
type Level int8

const (
	// FatalLevel is an abend error.
	FatalLevel Level = iota
	// ErrorLevel is unexpected runtime error.
	ErrorLevel
	// WarnLevel is warning.
	WarnLevel
	// InfoLevel is something notable infomation.
	InfoLevel
)

// String converts constant to string.
func (l Level) String() string {
	switch l {
	case FatalLevel:
		return "trace"
	case ErrorLevel:
		return "debug"
	case WarnLevel:
		return "info"
	case InfoLevel:
		return "warn"
	}
	return ""
}

// A Logger represents a logger.
type Logger struct{}

// NewLogger creates a logger.
func NewLogger() usecases.Logger {
	// TODO: set option?
	return &Logger{}
}

// LogError writes a log for an error log.
func (l *Logger) LogError(e error) {
	set(os.Stderr)
	write()
}

// LogAccess writes a log for an access log.
func (l *Logger) LogAccess(r *http.Request) {
	set(os.Stdout)

	startTime := time.Now()
	writer := &LoggingWriter{
		ResponseWriter: rw,
		logRecord: LogRecord{
			Time:          startTime.UTC(),
			Ip:            ip,
			Method:        r.Method,
			Uri:           r.RequestURI,
			Username:      username,
			Protocol:      r.Proto,
			Host:          r.Host,
			Status:        0,
			Size:          0,
			ElapsedTime:   time.Duration(0),
			RequestHeader: r.Header,
		},
	}

	if h.logBefore {
		writer.SetCustomLogRecord("at", "before")
		h.logger.LogContext(r.Context(), writer.logRecord)
	}

	ctx := context.WithValue(r.Context(), ctxLoggerKey, writer)
	r = r.WithContext(ctx)
	h.handler.ServeHTTP(writer, r)
	finishTime := time.Now()

	writer.logRecord.Time = finishTime.UTC()
	writer.logRecord.ElapsedTime = finishTime.Sub(startTime)

	write()
}

// Set sets logger options.
func set(writer *os.File) {
	log.SetOutput(writer)
	// TODO: 後で不要なやつ削る
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.Lshortfile | log.LUTC)
}

// Write writes logs.
func write() {
	// err
	log.Printf("%+v\n", errors.WithStack(e))

	// access
	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	log.Printf()
}
