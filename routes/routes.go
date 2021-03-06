package routes

import (
	"net/http"
	"path"
	"smartville-server/db"
	"smartville-server/repository"
	"smartville-server/middleware"
	"github.com/labstack/echo"
)

func InitRoute(ech *echo.Echo) {
	//Init db
	database := db.OpenDatabase()

	//Migrate
	db.Migrate(database)

	//Basic route
	ech.GET(path.Join("/"), func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Welcome to Smartville API")
	})

	//Admin Routes
	ech.GET("/admins", repository.GetAdminList(database))
	ech.GET("/adminById/:id", repository.GetAdminById(database))

	//User Routes
	ech.GET("/user-list", repository.GetUserList(database))
	ech.POST("/register", repository.Register(database))
	ech.POST("/login", repository.Login(database))
	ech.GET("user-id/:id", repository.GetUserById(database), middlewares.IsLoggedIn())

	//Email verification
	ech.POST("/email-verif", repository.SendEmail(database))
}