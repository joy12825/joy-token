// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/joy12825/gf.

package builtin

import (
	"errors"
	"reflect"

	"github.com/joy12825/gf/util/gconv"
)

// RuleRequired implements `required` rule.
// Format: required
type RuleRequired struct{}

func init() {
	Register(RuleRequired{})
}

func (r RuleRequired) Name() string {
	return "required"
}

func (r RuleRequired) Message() string {
	return "The {field} field is required"
}

func (r RuleRequired) Run(in RunInput) error {
	if isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}

func isRequiredEmpty(value interface{}) bool {
	reflectValue := reflect.ValueOf(value)
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectValue.Kind() {
	case reflect.String, reflect.Map, reflect.Array, reflect.Slice:
		return reflectValue.Len() == 0
	}
	return gconv.String(value) == ""
}
