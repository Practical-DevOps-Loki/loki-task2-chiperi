package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"webapp/pkg/config"
)

func main() {
	port := config.GetEnv("PORT", "3000")

	e := echo.New()

	logDir := "/tmp/webapp-logs"
	_ = os.MkdirAll(logDir, 0755)
	logFile, err := os.OpenFile(filepath.Join(logDir, "app.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, logFile)
		e.Logger.SetOutput(mw)
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: mw}))
	} else {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.Static("/", "public")

	e.GET("/", func(c echo.Context) error {
		return c.File("public/views/webapp.html")
	})
	e.Start(":" + port)
}
