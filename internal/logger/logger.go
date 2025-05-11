package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type SlogLogger struct {
	filePath string
	maxSize  int64 // in case uint, type conversion would be needed when comparing with file stat
	file     *os.File
	handler  *slog.Logger
	bckupNum uint8
}

func NewSlogLogger(filePath string) (*SlogLogger, error) {
	sl := &SlogLogger{
		filePath: filePath,
		maxSize:  5 * 1024 * 1024, // 5MB
		bckupNum: 0,
	}
	lFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	sl.file = lFile
	logger := slog.New(slog.NewTextHandler(sl.file, nil))

	sl.handler = logger
	slog.SetDefault(sl.handler)

	return sl, nil
}

//type Logger interface {
//	Info(msg, source, method, path, agent string, data any)
//	Error(msg, source, method, path, agent string, data any)
//}

func (l *SlogLogger) rotateFiles() error {
	if l.file != nil {
		stat, err := l.file.Stat()
		if err != nil {
			return err
		}
		if stat.Size() < l.maxSize {
			return nil
		}
		l.file.Close()
	}

    if l.bckupNum > 2 {
        // Think of a better way to handle this
        return nil
    }
    newFilePath := fmt.Sprintf("%s.%d", l.filePath, l.bckupNum+1)
    err := os.Rename(l.filePath, newFilePath)
    if err != nil {
        return err
    }

    lFile, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    l.file = lFile
    l.handler = slog.New(slog.NewTextHandler(l.file, nil))

	return nil
}

func (l *SlogLogger) Info(msg, source, method, path, agent string, data any) {
    l.rotateFiles()
	slog.Info(
		msg,
		"source", source,
		"method", method,
		"path", path,
		"user_agent", agent,
		"data", data,
	)
}

func (l *SlogLogger) Error(msg, source, method, path, agent string, data any) {
    l.rotateFiles()
	slog.Error(
		msg,
		"source", source,
		"method", method,
		"path", path,
		"user_agent", agent,
		"data", data,
	)
}
