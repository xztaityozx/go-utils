package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	ColorMap      map[string]string
	IsMultiLogger bool
	Logfile       string
	PrintColor    string
	FatalColor    string
}

func New() *Logger {
	return &Logger{
		ColorMap: map[string]string{
			"Red":    "\033[0;31m",
			"RedB":   "\033[1;31m",
			"Green":  "\033[0;32m",
			"GreenB": "\033[1;32m",
			"Reset":  "\033[0;39m",
		},
		IsMultiLogger: false,
		Logfile:       "",
		PrintColor:    "\033[0;39m",
		FatalColor:    "\033[0;39m",
	}
}

func (lgr *Logger) AppendLoggerColor(name string, ansi string) {
	lgr.ColorMap[name] = ansi
}

func (lgr *Logger) ToggleMultiLog() (bool, error) {
	lgr.IsMultiLogger = !lgr.IsMultiLogger

	if lgr.IsMultiLogger && lgr.Logfile == "" {
		return lgr.IsMultiLogger, errors.New("LogFile is not set")
	}

	if lgr.IsMultiLogger {
		setOutFile(lgr.Logfile)
	} else {
		setSingle(os.Stderr)
	}

	return lgr.IsMultiLogger, nil
}

func setOutFile(f string) error {
	if _, err := os.Stat(f); err != nil {
		return err
	}

	fp, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	log.SetOutput(io.MultiWriter(fp, os.Stderr))
	log.SetFlags(log.Ldate | log.Ltime)
	return nil
}

func setSingle(w io.Writer) {
	log.SetOutput(w)
	log.SetFlags(log.Ldate | log.Ltime)
}

func (lgr *Logger) SetPrintColorDirect(color string) {
	lgr.PrintColor = color
}

func (lgr *Logger) SetPrintColor(name string) {
	lgr.SetPrintColorDirect(lgr.ColorMap[name])
}

func (lgr *Logger) SetFatalColorDirect(color string) {
	lgr.FatalColor = color
}

func (lgr *Logger) SetFatalColor(name string) {
	lgr.SetFatalColorDirect(lgr.ColorMap[name])
}

func (lgr *Logger) SetLogFile(f string) error {
	if _, err := os.Stat(f); err != nil {
		return err
	}

	lgr.Logfile = f

	return setOutFile(f)
}

func (lgr *Logger) TrySetLogFile(f string) error {
	if _, stE := os.Stat(f); stE != nil {
		if _, cE := os.Create(f); cE != nil {
			return nil
		}
	}

	return lgr.SetLogFile(f)
}

func (lgr Logger) Print(v ...interface{}) {
	log.Print(lgr.PrintColor, v, lgr.ColorMap["Reset"])
}

func (lgr Logger) Println(v ...interface{}) {
	lgr.Print(v, "\n")
}

func (lgr Logger) Printf(format string, v ...interface{}) {
	lgr.Print(fmt.Sprintf(format, v))
}

func (lgr Logger) Fatal(v ...interface{}) {
	log.Fatal(lgr.FatalColor, v, lgr.ColorMap["Reset"])
}

func (lgr Logger) Fatalln(v ...interface{}) {
	lgr.Fatal(v, "\n")
}

func (lgr Logger) Fatalf(format string, v ...interface{}) {
	lgr.Fatal(fmt.Sprintf(format, v))
}

func (lgr Logger) SwitchPrint(fatal bool, v ...interface{}) {
	if fatal {
		lgr.Fatal(v)
	} else {
		lgr.Print(v)
	}
}

func (lgr Logger) FatalExit(err error) {
	if err != nil {
		lgr.Fatal(err)
	}
}
