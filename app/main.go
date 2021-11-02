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

	ech.Logger.Fatal(ech.Start(":" + cfg.Port))
}