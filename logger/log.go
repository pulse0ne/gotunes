package logger

import (
	"log"
	"os"
	"strings"
)

type Level int

const (
	Ldebug Level = 1 << iota
	Linfo
	Lwarn
	Lerror
)

type loglevel struct {
	level  Level
	logger *log.Logger
}

type logging struct {
	level Level
	debug *loglevel
	info  *loglevel
	warn  *loglevel
	error *loglevel
}

func (lg *logging) Debug(v ...interface{}) {
	if lg.level <= lg.debug.level {
		lg.debug.logger.Println(v...)
	}
}

func (lg *logging) Info(v ...interface{}) {
	if lg.level <= lg.info.level {
		lg.info.logger.Println(v...)
	}
}

func (lg *logging) Warn(v ...interface{}) {
	if lg.level <= lg.warn.level {
		lg.warn.logger.Println(v...)
	}
}

func (lg *logging) Error(v ...interface{}) {
	if lg.level <= lg.error.level {
		lg.error.logger.Println(v...)
	}
}

func (lg *logging) Fatal(v ...interface{}) {
	log.Fatalln(v...)
}

func NewLoggerFromString(level string) *logging {
	switch strings.ToLower(level) {
	case "debug":
		return NewLogger(Ldebug)
	case "warn":
		return NewLogger(Lwarn)
	case "error":
		return NewLogger(Lerror)
	default:
		return NewLogger(Linfo)
	}
}

func NewLogger(level Level) *logging {
	flags := log.Ldate | log.Ltime
	return &logging{
		level: level,
		debug: &loglevel{Ldebug, log.New(os.Stdout, "DEBUG: ", flags)},
		info:  &loglevel{Linfo, log.New(os.Stdout, "INFO: ", flags)},
		warn:  &loglevel{Lwarn, log.New(os.Stdout, "WARN: ", flags)},
		error: &loglevel{Lerror, log.New(os.Stderr, "ERROR: ", flags)},
	}
}
