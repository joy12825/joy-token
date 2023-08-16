// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/joy12825/gf.

// Package gcode provides universal error code definition and common error codes implements.
package gcode

// Code is universal error code interface definition.
type Code interface {
	// Code returns the integer number of current error code.
	Code() int

	// Message returns the brief message for current error code.
	Message() string

	// Detail returns the detailed information of current error code,
	// which is mainly designed as an extension field for error code.
	Detail() interface{}
}

// ================================================================================================================
// Common error code definition.
// There are reserved internal error code by framework: code < 1000.
// ================================================================================================================

var (
	CodeNil                       = localCode{-1, "", nil}                                // No error code specified.
	CodeOK                        = localCode{10000, "OK", nil}                           // It is OK.
	CodeInternalError             = localCode{10050, "Internal Error", nil}               // An error occurred internally.
	CodeValidationFailed          = localCode{10051, "Validation Failed", nil}            // Data validation failed.
	CodeDbOperationError          = localCode{10052, "Database Operation Error", nil}     // Database operation error.
	CodeInvalidParameter          = localCode{10053, "Invalid Parameter", nil}            // The given parameter for current operation is invalid.
	CodeMissingParameter          = localCode{10054, "Missing Parameter", nil}            // Parameter for current operation is missing.
	CodeInvalidOperation          = localCode{10055, "Invalid Operation", nil}            // The function cannot be used like this.
	CodeInvalidConfiguration      = localCode{10056, "Invalid Configuration", nil}        // The configuration is invalid for current operation.
	CodeMissingConfiguration      = localCode{10057, "Missing Configuration", nil}        // The configuration is missing for current operation.
	CodeNotImplemented            = localCode{10058, "Not Implemented", nil}              // The operation is not implemented yet.
	CodeNotSupported              = localCode{10059, "Not Supported", nil}                // The operation is not supported yet.
	CodeOperationFailed           = localCode{10060, "Operation Failed", nil}             // I tried, but I cannot give you what you want.
	CodeNotAuthorized             = localCode{10061, "Not Authorized", nil}               // Not Authorized.
	CodeSecurityReason            = localCode{10062, "Security Reason", nil}              // Security Reason.
	CodeServerBusy                = localCode{10063, "Server Is Busy", nil}               // Server is busy, please try again later.
	CodeUnknown                   = localCode{10064, "Unknown Error", nil}                // Unknown error.
	CodeNotFound                  = localCode{10065, "Not Found", nil}                    // Resource does not exist.
	CodeInvalidRequest            = localCode{10066, "Invalid Request", nil}              // Invalid request.
	CodeNecessaryPackageNotImport = localCode{10067, "Necessary Package Not Import", nil} // It needs necessary package import.
	CodeBusinessValidationFailed  = localCode{10300, "Business Validation Failed", nil}   // Business validation failed.
)

// New creates and returns an error code.
// Note that it returns an interface object of Code.
func New(code int, message string, detail interface{}) Code {
	return localCode{
		code:    code,
		message: message,
		detail:  detail,
	}
}

// WithCode creates and returns a new error code based on given Code.
// The code and message is from given `code`, but the detail if from given `detail`.
func WithCode(code Code, detail interface{}) Code {
	return localCode{
		code:    code.Code(),
		message: code.Message(),
		detail:  detail,
	}
}
