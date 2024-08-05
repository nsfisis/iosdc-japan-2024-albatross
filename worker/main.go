package main

import (
	"log"
	"net/http"

	// echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := prepareDirectories(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// TODO: temporarily disabled
	// e.Use(echojwt.WithConfig(echojwt.Config{
	// 	SigningKey: []byte("TODO"),
	// }))

	e.POST("/api/swiftc", handleSwiftCompile)
	e.POST("/api/wasmc", handleWasmCompile)
	e.POST("/api/testrun", handleTestRun)

	if err := e.Start(":80"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
