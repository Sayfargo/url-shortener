package core_slogger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrNoLevelOrFile = errors.New("no level or file in config")
)

type Slogger struct {
	*slog.Logger

	file *os.File
}

func MustNew(level string, dir string) *Slogger {
	l, err := New(level, dir)
	if err != nil {
		panic(err.Error())
	}
	return l
}

func New(level string, dir string) (*Slogger, error) {

	if level == "" || dir == "" {
		return nil, ErrNoLevelOrFile
	}

	var slogLevel slog.Level

	if err := slogLevel.UnmarshalText([]byte(strings.ToUpper(level))); err != nil {
		return nil, fmt.Errorf("failed to unmarshal level : %w", err)
	}

	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	})

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("make log dir: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")

	logFilePath := filepath.Join(
		dir,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	fileHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slogLevel,
	})

	multiHandler := slog.NewMultiHandler(stdoutHandler, fileHandler)

	slogger := &Slogger{
		Logger: slog.New(multiHandler),
		file:   logFile,
	}

	return slogger, nil

}

func (l *Slogger) With(args ...any) *Slogger {
	return &Slogger{
		Logger: l.Logger.With(args...),
		file:   l.file,
	}
}

func (l *Slogger) Close() error {
	if err := l.file.Close(); err != nil {
		return fmt.Errorf("failed to close log file : %w", err)
	}
	return nil
}
