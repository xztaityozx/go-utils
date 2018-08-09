package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Logger struct {
	ColorMap      map[string]string
	IsMultiLogger bool
	Logfile       string
	PrintColor    string
	FatalColor    string
	Logger        *log.Logger
}

func New() *Logger {
	return &Logger{
		ColorMap: map[string]string{
			"Red":    "\033[0;31m",
			"RedB":   "\033[1;31m",
			"Green":  "\033[0;32m",
			"GreenB": "\033[1;32m",
			"Bold":   "\033[1;39m",
			"Reset":  "\033[0;39m",
		},
		IsMultiLogger: false,
		Logfile:       "/dev/null",
		PrintColor:    "\033[0;39m",
		FatalColor:    "\033[0;39m",
		Logger:        log.New(os.Stderr, "", log.Ltime|log.Ldate),
	}
}

func (lgr *Logger) AppendLoggerColor(name string, ansi string) {
	lgr.ColorMap[name] = ansi
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

	return nil
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
	fp, err := os.OpenFile(lgr.Logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	defer fp.Close()
	if err != nil {
		return
	}
	lgr.Logger.SetOutput(io.MultiWriter(os.Stderr, fp))
	lgr.Logger.Print(lgr.PrintColor, v, lgr.ColorMap["Reset"])
}

func (lgr Logger) Println(v ...interface{}) {
	lgr.Print(v, "\n")
}

func (lgr Logger) Printf(format string, v ...interface{}) {
	lgr.Print(fmt.Sprintf(format, v))
}

func (lgr Logger) Fatal(v ...interface{}) {
	fp, err := os.OpenFile(lgr.Logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	defer fp.Close()
	if err != nil {
		return
	}
	lgr.Logger.SetOutput(io.MultiWriter(os.Stderr, fp))
	lgr.Logger.Fatal(lgr.FatalColor, v, lgr.ColorMap["Reset"])
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

func (lgr Logger) FatalStdErrOnly(v ...interface{}) {
	lgr.Logger.SetOutput(os.Stderr)
	lgr.Fatal(v)
}

func (lgr Logger) FatallnStdErrOnly(v ...interface{}) {
	lgr.FatalStdErrOnly(v, "\n")
}

func (lgr Logger) FatalfStdErrOnly(format string, v ...interface{}) {
	lgr.FatalStdErrOnly(fmt.Sprintf(format, v))
}
func (lgr Logger) PrintStdErrOnly(v ...interface{}) {
	lgr.Logger.SetOutput(os.Stderr)
	lgr.Print(v)
}

func (lgr Logger) PrintlnStdErrOnly(v ...interface{}) {
	lgr.PrintStdErrOnly(v, "\n")
}

func (lgr Logger) PrintfStdErrOnly(format string, v ...interface{}) {
	lgr.PrintStdErrOnly(fmt.Sprintf(format, v))
}

func (lgr Logger) FatalFileOnly(v ...interface{}) {
	fp, err := os.OpenFile(lgr.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer fp.Close()
	if err != nil {
		return
	}
	lgr.Logger.SetOutput(fp)
	lgr.Fatal(v)
}

func (lgr Logger) FatallnFileOnly(v ...interface{}) {
	lgr.FatalFileOnly(v, "\n")
}

func (lgr Logger) FatalfFileOnly(format string, v ...interface{}) {
	lgr.FatalFileOnly(fmt.Sprintf(format, v))
}
func (lgr Logger) PrintFileOnly(v ...interface{}) {
	fp, err := os.OpenFile(lgr.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer fp.Close()
	if err != nil {
		return
	}
	lgr.Logger.SetOutput(fp)
	lgr.Print(v)
}

func (lgr Logger) PrintlnFileOnly(v ...interface{}) {
	lgr.PrintFileOnly(v, "\n")
}

func (lgr Logger) PrintfFileOnly(format string, v ...interface{}) {
	lgr.PrintFileOnly(fmt.Sprintf(format, v))
}

func (lgr Logger) PrintSeparator(sep string, length int) {
	lgr.Logger.SetFlags(0)
	lgr.Printf("%s%s%s\n", lgr.ColorMap["Bold"], strings.Repeat(sep, length), lgr.ColorMap["Reset"])
	lgr.Logger.SetFlags(log.Ldate | log.Ltime)
}

func (lgr Logger) FatalSeparator(sep string, length int) {
	lgr.Logger.SetFlags(0)
	lgr.Fatalf("%s%s%s\n", lgr.ColorMap["RedB"], strings.Repeat(sep, length), lgr.ColorMap["Reset"])
	lgr.Logger.SetFlags(log.Ldate | log.Ltime)
}

func (lgr Logger) FatalExit(err error) {
	if err != nil {
		lgr.Fatal(err)
	}
}

func (lgr Logger) PrintFromFile(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	lgr.Print(string(b))
	return nil
}

func (lgr Logger) FatalFromFile(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	lgr.Fatal(string(b))
	return nil
}
