package ginger

import (
	"fmt"
	"net/http"
	"rate-limit-request/util/logger"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Wrap response struct
type Wrap interface {
	WithData(data interface{})
	HandelResponse(httpStatus int, message string)
	HandelErrorResponse(err error, httpStatus int, message string)
}

// Ginger ginger
type Ginger struct {
	*gin.Context
	wrap Wrap
	err  error
}

// ExtendResponse extend Gin context with response feature
func ExtendResponse(parent *gin.Context, wrap Wrap) *Ginger {
	return &Ginger{
		Context: parent,
		wrap:    wrap,
	}
}

// Response serializes the given struct as JSON into the response body.
func (c *Ginger) Response(httpStatus int, message string) {

	if c.err == nil {
		c.wrap.HandelResponse(httpStatus, message)
	} else {
		c.wrap.HandelErrorResponse(c.err, httpStatus, message)
		c.logError(httpStatus, message)
	}

	c.JSON(httpStatus, c.wrap)
	c.Abort()
}

// WithData set response struct with data
func (c *Ginger) WithData(data interface{}) *Ginger {
	c.wrap.WithData(data)
	return c
}

// WithError set error
func (c *Ginger) WithError(err error) *Ginger {
	c.err = err
	return c
}

// Done invoke parent Done method
func (c *Ginger) Done() <-chan struct{} {
	return c.Context.Done()
}

// Err invoke parent Err method
func (c *Ginger) Err() error {
	return c.Context.Err()
}

// Value invoke parent Value method
func (c *Ginger) Value(key interface{}) interface{} {
	return c.Context.Value(key)
}

func (c *Ginger) logError(httpStatus int, message string) {
	_, file, line, _ := runtime.Caller(2)
	l := logger.FromContext(c.Context).WithFields(logrus.Fields{
		"status":     httpStatus,
		"respSource": fmt.Sprintf("%s:%d", file, line),
	})
	if httpStatus >= http.StatusInternalServerError {
		l.WithError(c.err).Error(message)
	} else {
		l.WithField("err", c.err).Warning(message)
	}
}
