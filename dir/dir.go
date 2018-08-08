package dirs

import (
	"os"

	"github.com/xztaityozx/go-utils/logger"
)

var Logger *logger.Logger = logger.New()

func TryMkDir(f string, perm os.FileMode) error {
	if _, statE := os.Stat(f); statE != nil {
		if mkE := os.MkdirAll(f, perm); mkE != nil {
			return mkE
		}
		Logger.Println("Mkdir : ", f)
	}
	return nil
}

func TryMkDirAuto(f string) error {
	return TryMkDir(f, 0644)
}

func TryMkDirSuppress(f string, perm os.FileMode) error {
	if _, statE := os.Stat(f); statE != nil {
		if mkE := os.MkdirAll(f, perm); mkE != nil {
			return mkE
		}
	}
	return nil
}

func TryMkDirSuppressAuto(f string) error {
	return TryMkDirSuppress(f, 0644)
}

func TryChDir(f string, perm os.FileMode) error {
	if _, statE := os.Stat(f); statE != nil {
		if mkE := TryMkDir(f, perm); mkE != nil {
			return mkE
		}

		if chE := os.Chdir(f); chE != nil {
			return chE
		}

		Logger.Println("Chdir : ", f)
	}
	return nil
}
func TryChDirSuppress(f string, perm os.FileMode) error {
	if _, statE := os.Stat(f); statE != nil {
		if mkE := TryMkDirSuppress(f, perm); mkE != nil {
			return mkE
		}

		if chE := os.Chdir(f); chE != nil {
			return chE
		}

	}
	return nil
}

func TryChDirAuto(f string) error {
	return TryChDir(f, 0644)
}

func TryChDirSuppressAuto(f string) error {
	return TryChDirSuppress(f, 0644)
}
