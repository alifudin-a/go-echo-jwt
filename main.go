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
	// e.Renderer = helpers.NewRenderer("./view/index.html", true)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "==> METHOD=${method}, URI=${uri}, STATUS=${status}, " +
			"HOST=${host}, ERROR=${error}, LATENCY_HUMAN=${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// serve file
	e.Static("/", "view")

	// Login route
	e.POST("/login", auth.Login)
	e.POST("/login_v2", auth.LoginV2)
	e.GET("/accessible", auth.Accessible)

	// access := os.Getenv("JWT_ACCESS_SECRET")
	// refresh := os.Getenv("JWT_REFRESH_SECRET")

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte(os.Getenv("JWT_ACCESS_SECRET"))))
	r.GET("", auth.Restricted)

	e.Logger.Fatal(e.Start(":8900"))
}
