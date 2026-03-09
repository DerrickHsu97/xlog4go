package xlog4go

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type ConfFileWriter struct {
	On                  bool   `json:"On"`
	LogPath             string `json:"LogPath"`
	RotateLogPath       string `json:"RotateLogPath"`
	WfLogPath           string `json:"WfLogPath"`
	RotateWfLogPath     string `json:"RotateWfLogPath"`
	PublicLogPath       string `json:"PublicLogPath"`
	RotatePublicLogPath string `json:"RotatePublicLogPath"`
}

type ConfConsoleWriter struct {
	On    bool `json:"On"`
	Color bool `json:"Color"`
}

type LogConfig struct {
	Level string            `json:"LogLevel"`
	FW    ConfFileWriter    `json:"FileWriter"`
	CW    ConfConsoleWriter `json:"ConsoleWriter"`
}

type NamedLogConfig struct {
	Name   string     `json:"Name"`
	Config LogConfig  `json:"Config"`
}

type MultiLogConfig struct {
	Default LogConfig        `json:"Default"`
	Loggers []NamedLogConfig `json:"Loggers"`
}

func SetupLogWithConf(file string) (err error) {
	var lc LogConfig

	cnt, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(cnt, &lc); err != nil {
		return err
	}

	return setupLogger(logger_default, lc)
}

func SetupMultiLogWithConf(file string) (err error) {
	var mlc MultiLogConfig

	cnt, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(cnt, &mlc); err != nil {
		return err
	}

	// Setup default logger
	if err = setupLogger(logger_default, mlc.Default); err != nil {
		return err
	}

	// Setup named loggers
	for _, namedConf := range mlc.Loggers {
		l := NewLoggerWithName(namedConf.Name)
		if err = setupLogger(l, namedConf.Config); err != nil {
			return err
		}
	}

	return nil
}

func setupLogger(l *Logger, lc LogConfig) error {
	if lc.FW.On {
		if len(lc.FW.LogPath) > 0 {
			w := NewFileWriter()
			w.SetFileName(lc.FW.LogPath)
			w.SetPathPattern(lc.FW.RotateLogPath)
			w.SetLogLevelFloor(TRACE)
			if len(lc.FW.WfLogPath) > 0 {
				w.SetLogLevelCeil(INFO)
			} else {
				w.SetLogLevelCeil(ERROR)
			}
			l.Register(w)
		}

		if len(lc.FW.WfLogPath) > 0 {
			wfw := NewFileWriter()
			wfw.SetFileName(lc.FW.WfLogPath)
			wfw.SetPathPattern(lc.FW.RotateWfLogPath)
			wfw.SetLogLevelFloor(WARNING)
			wfw.SetLogLevelCeil(ERROR)
			l.Register(wfw)
		}

		if len(lc.FW.PublicLogPath) > 0 {
			public := NewFileWriter()
			public.SetFileName(lc.FW.PublicLogPath)
			public.SetPathPattern(lc.FW.RotatePublicLogPath)
			public.SetLogLevelFloor(PUBLIC)
			public.SetLogLevelCeil(PUBLIC)
			l.Register(public)
		}
	}

	if lc.CW.On {
		w := NewConsoleWriter()
		w.SetColor(lc.CW.Color)
		l.Register(w)
	}

	switch lc.Level {
	case "public":
		l.SetLevel(PUBLIC)

	case "trace":
		l.SetLevel(TRACE)

	case "debug":
		l.SetLevel(DEBUG)

	case "info":
		l.SetLevel(INFO)

	case "warning":
		l.SetLevel(WARNING)

	case "error":
		l.SetLevel(ERROR)

	case "fatal":
		l.SetLevel(FATAL)

	default:
		return errors.New("Invalid log level")
	}
	return nil
}
