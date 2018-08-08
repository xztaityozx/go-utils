package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAllLogger(t *testing.T) {
	home := os.Getenv("HOME")
	os.MkdirAll(filepath.Join(home, "WorkSpace"), 0644)

	lgr := New()

	t.Run("001_SetLogFile", func(t *testing.T) {
		f := filepath.Join(home, "WorkSpace", "out.log")
		if err := lgr.SetLogFile(f); err == nil {
			t.Fatal(err)
		}

		if lgr.Logfile != "" {
			t.Fail()
		}

		if err := lgr.TrySetLogFile(f); err != nil {
			t.Fail()
		}

		if lgr.Logfile != f {
			t.Fail()
		}

	})

	t.Run("002_SetPrintColor", func(t *testing.T) {
		if lgr.PrintColor != "\033[0;39m" {
			t.Fail()
		}

		lgr.SetPrintColorDirect("ABC")
		if lgr.PrintColor != "ABC" {
			t.Fail()
		}

		lgr.SetPrintColor("Red")
		if lgr.PrintColor != "\033[0;31m" {
			t.Fail()
		}
	})
	t.Run("003_SetFatalColor", func(t *testing.T) {
		if lgr.FatalColor != "\033[0;39m" {
			t.Fail()
		}

		lgr.SetFatalColorDirect("ABC")
		if lgr.FatalColor != "ABC" {
			t.Fail()
		}

		lgr.SetFatalColor("Red")
		if lgr.FatalColor != "\033[0;31m" {
			t.Fail()
		}
	})

	t.Run("004_AppendLoggerColor", func(t *testing.T) {
		lgr.AppendLoggerColor("A", "A")
		if lgr.ColorMap["A"] != "A" {
			t.Fail()
		}
	})

	os.Remove(lgr.Logfile)
}
