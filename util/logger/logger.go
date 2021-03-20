package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type loggerKey string

const contextKey = loggerKey("_loggerKey_context")

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(newSourceHook())
}

// SetLevel set log level
func SetLevel(level string) error {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lv)
	return nil
}

// Entry get entry
func Entry() *logrus.Entry {
	return logrus.NewEntry(logrus.StandardLogger())
}

// FromContext get entry frome context
func FromContext(c context.Context) *logrus.Entry {
	logger := c.Value(contextKey)
	if logger == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return logger.(*logrus.Entry)
}

// GinMiddleware set logrus into context
func GinMiddleware(keys ...string) gin.HandlerFunc {
	logrus.AddHook(newContextHook(keys...))
	return func(c *gin.Context) {
		// set headers into context
		for _, k := range keys {
			if v := c.GetHeader(k); v != "" {
				c.Set(k, v)
				continue
			}
			if v, ok := c.GetQuery(k); ok {
				c.Set(k, v)
				continue
			}
		}
		// set context values into logrus entry.Data
		entry := FromContext(c).WithContext(c)

		// set logrus entry into context
		c.Set(string(contextKey), entry)
		c.Next()
	}
}
