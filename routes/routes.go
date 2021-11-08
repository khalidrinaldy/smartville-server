package routes

import (
	"smartville-server/db"
	"smartville-server/repository"
	"github.com/labstack/echo"
)

func InitRoute(ech *echo.Echo) {
	//Init db
	database := db.OpenDatabase()

	//Migrate
	db.Migrate(database)

	//Admin Routes
	ech.GET("/admins", repository.GetAdminList(database))
	ech.GET("/adminById/:id", repository.GetAdminById(database))

	//User Routes
	ech.POST("/register", repository.Register(database))
	ech.POST("/login", repository.Login(database))
}