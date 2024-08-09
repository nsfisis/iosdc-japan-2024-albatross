package main

import (
	"log"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	jwtSecret := os.Getenv("ALBATROSS_JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("ALBATROSS_JWT_SECRET is not set")
	}

	if err := prepareDirectories(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecret),
	}))

	e.POST("/api/swiftc", handleSwiftCompile)
	e.POST("/api/wasmc", handleWasmCompile)
	e.POST("/api/testrun", handleTestRun)

	if err := e.Start(":80"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
