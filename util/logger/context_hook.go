package logger

import (
	"github.com/sirupsen/logrus"
)

type contextHook struct {
	keys []string
}

func newContextHook(keys ...string) logrus.Hook {
	return &contextHook{keys}
}

// Levels implement levels
func (hook contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implement fire
func (hook contextHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		for _, k := range hook.keys {
			if v := entry.Context.Value(k); v != nil {
				entry.Data[k] = v
			}
		}
	}
	return nil
}
