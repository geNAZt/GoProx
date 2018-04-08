package log

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"strconv"
	"strings"
)

func createHeader() string {
	// Add time
	header := "[" + time.Now().Format("02.01. 15:04:05.999999999") + "] "

	// Append file and line
	_, file, line, ok := runtime.Caller(2)
	if ok {
		// Skip file until we hit "gomint.io"
		index := strings.Index(file, "gomint.io")
		header += "[" + file[index:] + ":" + strconv.Itoa(line) + "] "
	}

	return header
}

func Info(format string, v ...interface{}) {
	if *DebugEnabled {
		fmt.Printf(createHeader()+"[INFO] "+format+"\n", v...)
	}
}

func Warn(format string, v ...interface{}) {
	if *DebugEnabled {
		fmt.Printf(createHeader()+"[WARN] "+format+"\n", v...)
	}
}

func Debug(format string, v ...interface{}) {
	if *DebugEnabled {
		fmt.Printf(createHeader()+"[DEBUG] "+format+"\n", v...)
	}
}

func Fatal(format string, v ...interface{}) {
	fmt.Printf(createHeader()+"[FATAL] "+format+"\n", v...)
	os.Exit(-1)
}
