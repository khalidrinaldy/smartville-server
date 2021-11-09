package main

import (
	"fmt"
	"os"
	"smartville-server/config"
	"smartville-server/routes"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("port no ERROR")
		port = cfg.Port
	}

	ech.Logger.Fatal(ech.Start(":" + port))
}