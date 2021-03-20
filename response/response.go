package response

import (
	"net/http"
	"rate-limit-request/util/ginger"
)

// PrefixCode response code prefix
var PrefixCode *int

type response struct {
	Meta struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}

// NewWrap response struct
func NewWrap() ginger.Wrap {
	return &response{}
}

func (r *response) WithData(data interface{}) {
	r.Data = data
}

func (r *response) HandelResponse(httpStatus int, msg string) {
	r.Meta.Code = r.formatCode(httpStatus)
	r.Meta.Status = statusString(httpStatus)
	r.Meta.Message = msg
}

func (r *response) HandelErrorResponse(err error, httpStatus int, msg string) {
	switch e := err.(type) {
	case CustomError:
		r.Meta.Code = e.Code
		r.Meta.Status = statusString(e.httpStatus)
		r.Meta.Message = e.Message
	default:
		r.HandelResponse(httpStatus, msg)
	}
}

func statusString(httpStatus int) string {
	if httpStatus < http.StatusBadRequest {
		return "success"
	}
	return "fail"
}

func (r *response) formatCode(code int) int {
	if code >= 1000 {
		// Log.Errorf("Invalid Error Code %d", code)
	}
	if PrefixCode != nil {
		code = (*PrefixCode * 1000) + code
	}
	r.Meta.Code = code
	return code
}
