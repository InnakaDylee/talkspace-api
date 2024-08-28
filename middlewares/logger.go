package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Logger(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339_nano}, remote_ip=${remote_ip}, user_agent=${user_agent}, method=${method}, host=${host}, uri=${uri}, status=${status}, error=${error}, latency=${latency_human}\n",
	}))
}