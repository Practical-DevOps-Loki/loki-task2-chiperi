package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	loggerConfig := middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}","method":"${method}","uri":"${uri}","status":${status}}\n`,
	}
	if err == nil {
		mw := io.MultiWriter(os.Stdout, logFile)
		e.Logger.SetOutput(mw)
		loggerConfig.Output = mw
		e.Use(middleware.LoggerWithConfig(loggerConfig))
	} else {
		e.Use(middleware.LoggerWithConfig(loggerConfig))
	}
	e.Use(middleware.Recover())
	e.Static("/", "public")

	e.GET("/", func(c echo.Context) error {
		return c.File("public/views/webapp.html")
	})
	go func() {
		time.Sleep(2 * time.Second)
		_, _ = http.Get("http://127.0.0.1:" + port + "/")
	}()
	e.Start(":" + port)
}
