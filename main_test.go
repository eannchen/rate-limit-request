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

	for i := 1; i <= model.RateLimitMaximum+5; i++ {
		r.GET("/app").Run(app.RouterEngine(), func(tResp gofight.HTTPResponse, tReq gofight.HTTPRequest) {
			if i <= model.RateLimitMaximum {
				if assert.Equal(t, tResp.Code, http.StatusOK) {
					var body struct {
						Meta struct {
							Code    int    `json:"code"`
							Status  string `json:"status"`
							Message string `json:"message"`
						} `json:"meta"`
						Data model.RateLimit `json:"data"`
					}
					assert.Nil(t, json.Unmarshal([]byte(tResp.Body.String()), &body))
					fmt.Printf("TEST %03d: %+v\n", i, body)
					assert.Equal(t, body.Data.Count, i)
				}
			}
			if i > model.RateLimitMaximum {
				if assert.Equal(t, tResp.Code, http.StatusTooManyRequests) {
					var body struct {
						Meta struct {
							Code    int    `json:"code"`
							Status  string `json:"status"`
							Message string `json:"message"`
						} `json:"meta"`
						Data interface{} `json:"data"`
					}
					assert.Nil(t, json.Unmarshal([]byte(tResp.Body.String()), &body))
					fmt.Printf("TEST %03d: %+v\n", i, body)
				}
			}
		})
	}
}
