package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Logger struct {
	*log.Logger
	LogFileName string
	LogFile     *os.File
	LogLevel    int
	OutputMode  int
}

var (
	DEBUG int = -1
	INFO  int = 0 // default
	WARN  int = 1
	ERROR int = 2

	STDERR_ONLY int = 0
	FILE_ONLY   int = 1
	BOTH        int = 2
)

var LOG_LEVEL = map[string]int{
	"debug": DEBUG,
	"info":  INFO,
	"warn":  WARN,
	"error": ERROR,
}

var OUTPUT_MODE = map[string]int{
	"stderr": STDERR_ONLY,
	"file":   FILE_ONLY,
	"both":   BOTH,
}

var logger *Logger

func (l *Logger) Write(p []byte) (int, error) {
	if l.OutputMode == STDERR_ONLY {
		return os.Stderr.Write(p)
	}
	if l.OutputMode == FILE_ONLY {
		return l.LogFile.Write(p)
	}

	_, err := l.LogFile.Write(p)
	if err != nil {
		return os.Stderr.Write([]byte(err.Error()))
	}
	return os.Stderr.Write(p)
}

func SetupLogger(logFileName string, outputMode string, logLevel string) (*Logger, error) {
	lgr := new(Logger)
	lgr.Logger = log.New(lgr, "", log.LstdFlags)

	lgr.LogFileName = logFileName
	if err := os.MkdirAll(filepath.Dir(logFileName), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Creating log file directory \"%s\" was failed.\nError: %v\n", filepath.Dir(logFileName), err)
		return nil, err
	}

	if _, is := OUTPUT_MODE[outputMode]; !is {
		fmt.Fprintf(os.Stderr, "log output mode \"%d\" is invalid. Valid value is [0-2].\n", outputMode)
		return nil, errors.New("Specified log output mode is invaild.")
	}
	lgr.OutputMode = OUTPUT_MODE[outputMode]

	if lgr.OutputMode > 0 {
		file, err := os.OpenFile(lgr.LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Creating log file \"%s\" was failed.\nError: %v\n", logFileName, err)
			return nil, err
		}
		lgr.LogFile = file
	}

	if _, exist := LOG_LEVEL[logLevel]; !exist {
		fmt.Fprintf(os.Stderr, "log level \"%s\" is invalid.\n", logLevel)
		return nil, errors.New("Specified log level is invaild.")
	}
	lgr.LogLevel = LOG_LEVEL[logLevel]

	return lgr, nil
}

func Infof(format string, a ...interface{}) {
	loging(INFO, "[info]", format, a...)
}

func Warnf(format string, a ...interface{}) {
	loging(WARN, "[warn]", format, a...)
}

func Errorf(format string, a ...interface{}) {
	loging(ERROR, "[error]", format, a...)
}

func Debugf(format string, a ...interface{}) {
	loging(DEBUG, "[debug]", format, a...)
}

func loging(logLevel int, header string, format string, a ...interface{}) {
	if logger == nil {
		fmt.Printf(format+" (no logger configure)\n", a...)
		return
	}
	if logger.LogLevel > logLevel {
		return
	}
	logger.Printf(header+" "+format, a...)
}
