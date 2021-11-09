package routes

import (
	"net/http"
	"path"
	// "smartville-server/db"
	// "smartville-server/repository"

	"github.com/labstack/echo"
)

func InitRoute(ech *echo.Echo) {
	// //Init db
	// database := db.OpenDatabase()

	// //Migrate
	// db.Migrate(database)

	//Basic route
	ech.GET(path.Join("/"), func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Welcome to Smartville API")
	})

	// //Admin Routes
	// ech.GET("/admins", repository.GetAdminList(database))
	// ech.GET("/adminById/:id", repository.GetAdminById(database))

	// //User Routes
	// ech.POST("/register", repository.Register(database))
	// ech.POST("/login", repository.Login(database))
}