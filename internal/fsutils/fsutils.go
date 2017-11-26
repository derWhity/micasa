package fsutils

import (
	"fmt"
	"os"
	"syscall"

	"github.com/derWhity/micasa/internal/log"
)

// CheckAndCreateDir checks and tries to create the given directory recursively (or panics if this fails)
func CheckAndCreateDir(path string, logger log.Logger) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			logger.Info("Directory does not exist - trying to create...", log.FldPath, path)
			if err = os.MkdirAll(path, os.ModePerm); err != nil {
				logger.Crit("Failed to create directory", log.FldError, err)
				panic("Cannot continue")
			}
			logger.Info("Directory created successfully")
		} else {
			logger.Crit("Stat has failed", log.FldError, err.(*os.PathError).Err)
			panic("Cannot continue")
		}
	} else {
		if !fileInfo.IsDir() {
			logger.Crit(fmt.Sprintf("'%s' is not a directory. Remove the plain file if you want to continue", path))
			panic("Cannot continue")
		}
	}
}
