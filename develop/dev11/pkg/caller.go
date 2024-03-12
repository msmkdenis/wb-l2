package pkg

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
)

// Caller returns file name and line number of function call
func Caller() string {
	_, file, lineNo, ok := runtime.Caller(1)
	if !ok {
		return "runtime.Caller() failed"
	}

	fileName := path.Base(file)
	dir := filepath.Base(filepath.Dir(file))
	return fmt.Sprintf("%s/%s:%d", dir, fileName, lineNo)
}
