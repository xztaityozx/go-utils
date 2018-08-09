package logger

import (
	"fmt"
	"io"
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
	fp            *os.File
}

func (lgr *Logger) Close() error {
	if err := lgr.fp.Close(); err != nil {
		return err
	}

	return os.Stderr.Close()
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

func (lgr *Logger) ToggleMultiLog() (bool, error) {
	if lgr.IsMultiLogger {
		lgr.OffMultiLogger()
	} else {
		if err := lgr.OnMultiLogger(); err != nil {
			return false, err
		}
	}

	return lgr.IsMultiLogger, nil
}

func (lgr *Logger) OnMultiLogger() error {
	if lgr.IsMultiLogger {
		return nil
	}

	lgr.IsMultiLogger = true
	return setOutFile(lgr, lgr.Logfile)
}

func (lgr *Logger) OffMultiLogger() {
	if !lgr.IsMultiLogger {
		return
	}

	lgr.IsMultiLogger = false
	setSingle(lgr, os.Stderr)
}

func setOutFile(lgr *Logger, f string) error {
	if _, err := os.Stat(f); err != nil {
		return err
	}

	fp, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	lgr.fp = fp

	setSingle(lgr, io.MultiWriter(fp, os.Stderr))
	return nil
}

func setSingle(lgr *Logger, w io.Writer) {
	lgr.Logger.SetOutput(w)
	lgr.Logger.SetFlags(log.Ldate | log.Ltime)
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

	return setOutFile(lgr, f)
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
	setSingle(&lgr, os.Stderr)
	lgr.Fatal(v)
	setOutFile(&lgr, lgr.Logfile)
}

func (lgr Logger) FatallnStdErrOnly(v ...interface{}) {
	lgr.FatalStdErrOnly(v, "\n")
}

func (lgr Logger) FatalfStdErrOnly(format string, v ...interface{}) {
	lgr.FatalStdErrOnly(fmt.Sprintf(format, v))
}

func (lgr Logger) PrintStdErrOnly(v ...interface{}) {
	setSingle(&lgr, os.Stderr)
	lgr.Print(v)
	setOutFile(&lgr, lgr.Logfile)
}

func (lgr Logger) PrintlnStdErrOnly(v ...interface{}) {
	lgr.PrintStdErrOnly(v, "\n")
}

func (lgr Logger) PrintfStdErrOnly(format string, v ...interface{}) {
	lgr.PrintStdErrOnly(fmt.Sprintf(format, v))
}

func (lgr Logger) FatalFileOnly(v ...interface{}) {
	fp, _ := os.OpenFile(lgr.Logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	*lgr.fp = *fp
	w := io.MultiWriter(os.Stderr, fp)
	setSingle(&lgr, fp)

	lgr.Print(v)

	setSingle(&lgr, w)
}

func (lgr Logger) FatallnFileOnly(v ...interface{}) {
	lgr.FatallnFileOnly(v, "\n")
}

func (lgr Logger) FatalfFileOnly(format string, v ...interface{}) {
	lgr.FatallnFileOnly(fmt.Sprintf(format, v))
}
func (lgr Logger) PrintFileOnly(v ...interface{}) {
	fp, _ := os.OpenFile(lgr.Logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	*lgr.fp = *fp
	w := io.MultiWriter(os.Stderr, fp)
	setSingle(&lgr, fp)

	lgr.Print(v)

	setSingle(&lgr, w)
}

func (lgr Logger) PrintlnFileOnly(v ...interface{}) {
	lgr.PrintlnFileOnly(v, "\n")
}

func (lgr Logger) PrintfFileOnly(format string, v ...interface{}) {
	lgr.PrintlnFileOnly(fmt.Sprintf(format, v))
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
