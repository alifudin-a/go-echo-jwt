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
	e.File("/", "view/index.html")

	// Login route
	e.POST("/login", auth.Login)
	e.GET("/accessible", auth.Accessible)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))
	r.GET("", auth.Restricted)

	e.Logger.Fatal(e.Start(":8900"))
}
