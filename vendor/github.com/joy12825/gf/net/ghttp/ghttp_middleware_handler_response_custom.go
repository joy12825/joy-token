package ghttp

import (
	"net/http"

	"github.com/joy12825/gf/errors/gcode"
	"github.com/joy12825/gf/errors/gerror"
)

// CustomHandlerResponse is the custom implementation of HandlerResponse.
type CustomHandlerResponse struct {
	RetCode int         `json:"RetCode"    dc:"Error code"`
	Message string      `json:"Message" dc:"Error message"`
	Data    interface{} `json:"Data"    dc:"Result data for certain request according API definition"`
}

// MiddlewareCustomHandlerResponse is the CustomHandlerResponse middleware handling handler response object and its error.
func MiddlewareCustomHandlerResponse(r *Request) {
	r.Middleware.Next()

	// There's custom buffer content, it then exits current handler.
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		msg  string
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
	)
	if err != nil {
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		msg = err.Error()
	} else {
		if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
			msg = http.StatusText(r.Response.Status)
			switch r.Response.Status {
			case http.StatusNotFound:
				code = gcode.CodeNotFound
			case http.StatusForbidden:
				code = gcode.CodeNotAuthorized
			default:
				code = gcode.CodeUnknown
			}
			// It creates error as it can be retrieved by other middlewares.
			err = gerror.NewCode(code, msg)
			r.SetError(err)
		} else {
			code = gcode.CodeOK
		}
	}

	r.Response.WriteJson(CustomHandlerResponse{
		RetCode: code.Code(),
		Message: msg,
		Data:    res,
	})
}
