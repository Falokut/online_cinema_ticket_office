package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

var entry *logrus.Entry

type Mode int

const (
	Test Mode = iota
	FileOutput
	ConsoleOutput
	FileAndConsoleOutput
)

type Logger struct {
	*logrus.Entry
}

// Before calling this, need call NewEntry for entry initializing
func GetLogger() Logger {
	return Logger{entry}
}
func GetNullLogger() Logger {
	l, _ := test.NewNullLogger()
	return Logger{logrus.NewEntry(l)}
}
func (l *Logger) GetLoggerWithField(k string, v interface{}) Logger {
	return Logger{l.WithField(k, v)}
}

func NewEntry(mode Mode) {
	l := logrus.New()
	l.SetReportCaller(true)

	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	switch mode {
	case Test, ConsoleOutput:
		{
			l.SetOutput(os.Stdout)
		}
	case FileAndConsoleOutput:
		{
			err := os.MkdirAll("logs", 0755)

			if err != nil || os.IsExist(err) {
				panic("can't create log dir. no configured logging to files")
			}

			logfile, err := os.OpenFile("logs/all_logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
			if err != nil {
				panic(fmt.Sprintf("[Error]: %s", err))
			}

			l.SetOutput(io.MultiWriter(logfile, os.Stdout))
		}
	case FileOutput:
		{
			err := os.MkdirAll("logs", 0755)

			if err != nil || os.IsExist(err) {
				panic("can't create log dir. no configured logging to files")
			}

			logfile, err := os.OpenFile("logs/all_logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
			if err != nil {
				panic(fmt.Sprintf("[Error]: %s", err))
			}
			l.SetOutput(logfile)
		}
	}

	entry = logrus.NewEntry(l)
}
