package main

import (
	"os"

	"github.com/alifudin-a/go-echo-jwt/services/auth"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, " +
			"host=${host}, error=${error}, latency_human=${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// api := e.Group("/api")

	// e.File("/", "static/login.html")

	// Login route
	e.POST("/login", auth.Login)
	e.GET("/", auth.Accessible)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))
	r.GET("", auth.Restricted)

	e.Logger.Fatal(e.Start(":8900"))
}
