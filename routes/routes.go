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
	ech.GET("/adminsPost", repository.GetPostListJoinAdmin(database))
}