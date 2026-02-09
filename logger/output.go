package logger

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/lfcypo/viperx"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	Stdout = iota + 1
	File
	FileAndStdout
)

func rotaryFile() *lumberjack.Logger {
	logFile := viperx.GetString("log.file", "./logs/scrkun.log")
	logFileSplit := strings.Split(logFile, string(os.PathSeparator))
	if !strings.Contains(logFileSplit[len(logFileSplit)-1], ".") {
		logFile = path.Join(logFile, "scrkun.log")
	}

	return &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    viperx.GetInt("log.maxSize", 100),
		MaxBackups: viperx.GetInt("log.maxBackups", 10),
		MaxAge:     viperx.GetInt("log.maxAge", 30),
		Compress:   viperx.GetBool("log.compress", true),
		LocalTime:  true,
	}
}

func Output() io.Writer {
	mode := viperx.GetInt("log.mode", FileAndStdout)
	if mode == FileAndStdout {
		return io.MultiWriter(os.Stdout, rotaryFile())
	} else if mode == File {
		return rotaryFile()
	} else if mode == Stdout {
		return os.Stdout
	}

	panic("log mode is invalid")
}
