// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/joy12825/gf.

package ghttp

import (
	"context"
	"net/http"

	"github.com/joy12825/gf/v2/errors/gcode"
	"github.com/joy12825/gf/v2/errors/gerror"
	"github.com/joy12825/gf/v2/util/gutil"
)

// middleware is the plugin for request workflow management.
type middleware struct {
	served         bool     // Is the request served, which is used for checking response status 404.
	request        *Request // The request object pointer.
	handlerIndex   int      // Index number for executing sequence purpose for handler items.
	handlerMDIndex int      // Index number for executing sequence purpose for bound middleware of handler item.
}

// Next calls the next workflow handler.
// It's an important function controlling the workflow of the server request execution.
func (m *middleware) Next() {
	var item *HandlerItemParsed
	var loop = true
	for loop {
		// Check whether the request is excited.
		if m.request.IsExited() || m.handlerIndex >= len(m.request.handlers) {
			break
		}
		item = m.request.handlers[m.handlerIndex]
		// Filter the HOOK handlers, which are designed to be called in another standalone procedure.
		if item.Handler.Type == HandlerTypeHook {
			m.handlerIndex++
			continue
		}
		// Current router switching.
		m.request.Router = item.Handler.Router

		// Router values switching.
		m.request.routerMap = item.Values

		var ctx = m.request.Context()
		gutil.TryCatch(ctx, func(ctx context.Context) {
			// Execute bound middleware array of the item if it's not empty.
			if m.handlerMDIndex < len(item.Handler.Middleware) {
				md := item.Handler.Middleware[m.handlerMDIndex]
				m.handlerMDIndex++
				niceCallFunc(func() {
					md(m.request)
				})
				loop = false
				return
			}
			m.handlerIndex++

			switch item.Handler.Type {
			// Service object.
			case HandlerTypeObject:
				m.served = true
				if m.request.IsExited() {
					break
				}
				if item.Handler.InitFunc != nil {
					niceCallFunc(func() {
						item.Handler.InitFunc(m.request)
					})
				}
				if !m.request.IsExited() {
					m.callHandlerFunc(item.Handler.Info)
				}
				if !m.request.IsExited() && item.Handler.ShutFunc != nil {
					niceCallFunc(func() {
						item.Handler.ShutFunc(m.request)
					})
				}

			// Service handler.
			case HandlerTypeHandler:
				m.served = true
				if m.request.IsExited() {
					break
				}
				niceCallFunc(func() {
					m.callHandlerFunc(item.Handler.Info)
				})

			// Global middleware array.
			case HandlerTypeMiddleware:
				niceCallFunc(func() {
					item.Handler.Info.Func(m.request)
				})
				// It does not continue calling next middleware after another middleware done.
				// There should be a "Next" function to be called in the middleware in order to manage the workflow.
				loop = false
			}
		}, func(ctx context.Context, exception error) {
			if gerror.HasStack(exception) {
				// It's already an error that has stack info.
				m.request.error = exception
			} else {
				// Create a new error with stack info.
				// Note that there's a skip pointing the start stacktrace
				// of the real error point.
				m.request.error = gerror.WrapCodeSkip(gcode.CodeInternalError, 1, exception, "")
			}
			m.request.Response.WriteStatus(http.StatusInternalServerError, exception)
			loop = false
		})
	}
	// Check the http status code after all handlers and middleware done.
	if m.request.IsExited() || m.handlerIndex >= len(m.request.handlers) {
		if m.request.Response.Status == 0 {
			if m.request.Middleware.served {
				m.request.Response.WriteHeader(http.StatusOK)
			} else {
				m.request.Response.WriteHeader(http.StatusNotFound)
			}
		}
	}
}

func (m *middleware) callHandlerFunc(funcInfo handlerFuncInfo) {
	niceCallFunc(func() {
		funcInfo.Func(m.request)
	})
}
