package main

import (
	"os"
	//"smartville-server/config"
	"smartville-server/routes"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main()  {
	ech := echo.New()

	//Routes
	routes.InitRoute(ech)

	//Config
	//cfg,_ := config.NewConfig(".env")

	//use CORS
	ech.Use(middleware.CORS())

	//Set PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ech.Logger.Fatal(ech.Start(":" + port))
}