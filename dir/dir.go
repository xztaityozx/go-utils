package dirs

import (
	"os"

	"github.com/xztaityozx/go-utils/logger"
)

type Dirs struct {
	Logger *logger.Logger
}

func New() *Dirs {
	return &Dirs{
		Logger: logger.New(),
	}
}

func (dirs Dirs) TryMkDir(f string, perm os.FileMode) error {
	if _, statE := os.Stat(f); statE != nil {
		if mkE := os.MkdirAll(f, perm); mkE != nil {
			return mkE
		}
		dirs.Logger.Println("Mkdir : ", f)
	}
	return nil
}

func (dirs Dirs) TryMkDirAuto(f string) error {
	return dirs.TryMkDir(f, 0644)
}

func (dirs Dirs) TryMkDirSuppress(f string, perm os.FileMode) error {
	if _, statE := os.Stat(f); statE != nil {
		if mkE := os.MkdirAll(f, perm); mkE != nil {
			return mkE
		}
	}
	return nil
}

func (dirs Dirs) TryMkDirSuppressAuto(f string) error {
	return dirs.TryMkDirSuppress(f, 0644)
}

func (dirs Dirs) TryChDir(f string, perm os.FileMode) error {
	if _, statE := os.Stat(f); statE != nil {
		if mkE := dirs.TryMkDir(f, perm); mkE != nil {
			return mkE
		}

		if chE := os.Chdir(f); chE != nil {
			return chE
		}

		dirs.Logger.Println("Chdir : ", f)
	}
	return nil
}
func (dirs Dirs) TryChDirSuppress(f string, perm os.FileMode) error {
	if _, statE := os.Stat(f); statE != nil {
		if mkE := dirs.TryMkDirSuppress(f, perm); mkE != nil {
			return mkE
		}

		if chE := os.Chdir(f); chE != nil {
			return chE
		}

	}
	return nil
}

func (dirs Dirs) TryChDirAuto(f string) error {
	return dirs.TryChDir(f, 0644)
}

func (dirs Dirs) TryChDirSuppressAuto(f string) error {
	return dirs.TryChDirSuppress(f, 0644)
}
