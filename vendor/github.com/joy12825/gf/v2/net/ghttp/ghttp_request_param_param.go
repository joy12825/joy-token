// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/joy12825/gf.

package ghttp

import "github.com/joy12825/gf/v2/container/gvar"

// SetParam sets custom parameter with key-value pairs.
func (r *Request) SetParam(key string, value interface{}) {
	if r.paramsMap == nil {
		r.paramsMap = make(map[string]interface{})
	}
	r.paramsMap[key] = value
}

// SetParamMap sets custom parameter with key-value pair maps.
func (r *Request) SetParamMap(data map[string]interface{}) {
	if r.paramsMap == nil {
		r.paramsMap = make(map[string]interface{})
	}
	for k, v := range data {
		r.paramsMap[k] = v
	}
}

// GetParam returns custom parameter with a given name `key`.
// It returns `def` if `key` does not exist.
// It returns nil if `def` is not passed.
func (r *Request) GetParam(key string, def ...interface{}) *gvar.Var {
	if len(r.paramsMap) > 0 {
		if value, ok := r.paramsMap[key]; ok {
			return gvar.New(value)
		}
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}
