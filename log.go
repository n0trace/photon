package photon

import (
	"io"
)

type (
	// Logger defines the logging interface.
	Logger interface {
		Output() io.Writer
		SetOutput(w io.Writer)
		Prefix() string
		SetPrefix(p string)
		Print(i ...interface{})
		Printf(format string, args ...interface{})
		Debug(i ...interface{})
		Debugf(format string, args ...interface{})
		Info(i ...interface{})
		Infof(format string, args ...interface{})
		Warn(i ...interface{})
		Warnf(format string, args ...interface{})
		Error(i ...interface{})
		Errorf(format string, args ...interface{})
		Fatal(i ...interface{})
		Fatalf(format string, args ...interface{})
		Panic(i ...interface{})
		Panicf(format string, args ...interface{})
	}
)
