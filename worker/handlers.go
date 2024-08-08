package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func newBadRequestError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request: %s", err.Error()))
}

func handleSwiftCompile(c echo.Context) error {
	var req swiftCompileRequestData
	if err := c.Bind(&req); err != nil {
		return newBadRequestError(err)
	}
	if err := req.validate(); err != nil {
		return newBadRequestError(err)
	}

	res := execSwiftCompile(
		c.Request().Context(),
		req.Code,
		req.CodeHash,
		req.maxDuration(),
	)

	return c.JSON(http.StatusOK, res)
}

func handleWasmCompile(c echo.Context) error {
	var req wasmCompileRequestData
	if err := c.Bind(&req); err != nil {
		return newBadRequestError(err)
	}
	if err := req.validate(); err != nil {
		return newBadRequestError(err)
	}

	res := execWasmCompile(
		c.Request().Context(),
		req.CodeHash,
		req.maxDuration(),
	)

	return c.JSON(http.StatusOK, res)
}

func handleTestRun(c echo.Context) error {
	var req testRunRequestData
	if err := c.Bind(&req); err != nil {
		return newBadRequestError(err)
	}
	if err := req.validate(); err != nil {
		return newBadRequestError(err)
	}

	res := execTestRun(
		c.Request().Context(),
		req.CodeHash,
		req.Stdin,
		req.maxDuration(),
	)

	return c.JSON(http.StatusOK, res)
}
