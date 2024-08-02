package main

import (
	"errors"
	"time"
)

const (
	resultSuccess       = "success"
	resultFailure       = "failure"
	resultTimeout       = "timeout"
	resultInternalError = "internal_error"
)

var (
	errInvalidMaxDuration = errors.New("'max_duration_ms' must be positive")
)

type swiftCompileRequestData struct {
	MaxDurationMilliseconds int    `json:"max_duration_ms"`
	Code                    string `json:"code"`
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
	Result string `json:"result"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type wasmCompileRequestData struct {
	MaxDurationMilliseconds int    `json:"max_duration_ms"`
	Code                    string `json:"code"`
}

type wasmCompileResponseData struct {
	Result string `json:"result"`
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
	Code                    string `json:"code"`
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
	Result string `json:"result"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}
