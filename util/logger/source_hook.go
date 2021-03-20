package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type sourceHook struct {
	Field string
	Skip  int
}

func newSourceHook() logrus.Hook {
	return &sourceHook{
		Field: "source",
		Skip:  5,
	}
}

// Levels implement levels
func (hook sourceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implement fire
func (hook sourceHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
