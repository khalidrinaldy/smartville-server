package main

import (
	"smartville-server/config"
	"smartville-server/routes"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
)

func main()  {
	ech := echo.New()

	//Routes
	routes.InitRoute(ech)

	//Config
	cfg,_ := config.NewConfig(".env")

	//use CORS
	ech.Use(middleware.CORS())

	//Set PORT
	port := cfg.Port
	if cfg.Env == "production" {
		port = "8080"
	}

	ech.Logger.Fatal(ech.Start(":" + port))
}