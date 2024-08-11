package main

import (
	"errors"
	"time"
)

const (
	resultSuccess       = "success"
	resultCompileError  = "compile_error"
	resultRuntimeError  = "runtime_error"
	resultTimeout       = "timeout"
	resultInternalError = "internal_error"
)

var (
	errInvalidMaxDuration = errors.New("'max_duration_ms' must be positive")
)

type swiftCompileRequestData struct {
	MaxDurationMilliseconds int    `json:"max_duration_ms"`
	Code                    string `json:"code"`
	CodeHash                string `json:"code_hash"`
}

func (req *swiftCompileRequestData) maxDuration() time.Duration {
	return time.Duration(req.MaxDurationMilliseconds) * time.Millisecond
}

func (req *swiftCompileRequestData) validate() error {
	if req.MaxDurationMilliseconds <= 0 {
		return errInvalidMaxDuration
	}
	return nil
}

type swiftCompileResponseData struct {
	Status string `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type wasmCompileRequestData struct {
	MaxDurationMilliseconds int    `json:"max_duration_ms"`
	CodeHash                string `json:"code_hash"`
}

type wasmCompileResponseData struct {
	Status string `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func (req *wasmCompileRequestData) maxDuration() time.Duration {
	return time.Duration(req.MaxDurationMilliseconds) * time.Millisecond
}

func (req *wasmCompileRequestData) validate() error {
	if req.MaxDurationMilliseconds <= 0 {
		return errInvalidMaxDuration
	}
	return nil
}

type testRunRequestData struct {
	MaxDurationMilliseconds int    `json:"max_duration_ms"`
	CodeHash                string `json:"code_hash"`
	Stdin                   string `json:"stdin"`
}

func (req *testRunRequestData) maxDuration() time.Duration {
	return time.Duration(req.MaxDurationMilliseconds) * time.Millisecond
}

func (req *testRunRequestData) validate() error {
	if req.MaxDurationMilliseconds <= 0 {
		return errInvalidMaxDuration
	}
	return nil
}

type testRunResponseData struct {
	Status string `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}
