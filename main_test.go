package main

import (
	"net/http"
	"rate-limit-request/app"
	"rate-limit-request/model"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	r := gofight.New()

	var (
		successExpect = model.RateLimitMaximum
		failExpect    = 5
		finishCount   = successExpect + failExpect
		finishCh      = make(chan bool, finishCount)
		errCh         = make(chan struct{})
	)
	for i := 1; i <= finishCount; i++ {

		go func() {
			r.GET("/app").Run(app.RouterEngine(), func(tResp gofight.HTTPResponse, tReq gofight.HTTPRequest) {
				switch tResp.Code {
				case http.StatusOK:
					body := tResp.Body.String()
					if assert.Regexp(t, "^[0-9]+$", body) {
						finishCh <- true
						// fmt.Printf("%+v\n", body)
					} else {
						errCh <- struct{}{}
					}

				case http.StatusTooManyRequests:
					body := tResp.Body.String()
					if assert.Equal(t, "Error", body) {
						finishCh <- false
						// fmt.Printf("%+v\n", body)
					} else {
						errCh <- struct{}{}
					}
				}
			})
		}()
	}

	var (
		successResult int
		failResult    int
	)
	for {
		select {
		case successReq := <-finishCh:
			if successReq {
				successResult++
			} else {
				failResult++
			}
			if successResult+failResult == finishCount {
				assert.Equal(t, successExpect, successResult)
				assert.Equal(t, failExpect, failResult)
				return
			}
		case <-errCh:
			return
		}
	}
}
