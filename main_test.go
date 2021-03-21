package main

import (
	"encoding/json"
	"fmt"
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
					var body struct {
						Meta struct {
							Code    int    `json:"code"`
							Status  string `json:"status"`
							Message string `json:"message"`
						} `json:"meta"`
						Data model.RateLimit `json:"data"`
					}
					if assert.Nil(t, json.Unmarshal([]byte(tResp.Body.String()), &body)) {
						finishCh <- true
						fmt.Printf("%+v\n", body)
					} else {
						errCh <- struct{}{}
					}

				case http.StatusTooManyRequests:
					var body struct {
						Meta struct {
							Code    int    `json:"code"`
							Status  string `json:"status"`
							Message string `json:"message"`
						} `json:"meta"`
						Data interface{} `json:"data"`
					}
					if assert.Nil(t, json.Unmarshal([]byte(tResp.Body.String()), &body)) {
						finishCh <- false
						fmt.Printf("%+v\n", body)
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
