package routes

import (
	"net/http"
	"path"
	"smartville-server/db"
	"smartville-server/middleware"
	"smartville-server/repository"
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
	ech.GET("/user-id/:id", repository.GetUserById(database), middlewares.IsLoggedIn())
	ech.GET("/userbytoken", repository.GetUserByToken(database), middlewares.IsLoggedIn())
	ech.PUT("/user/edit", repository.EditProfile(database), middlewares.IsLoggedIn())

	//Email verification
	ech.POST("/user/email-verif", repository.SendEmail(database))

	//Change password
	ech.PUT("/user/forgot-password", repository.ChangeForgotPassword(database))
	ech.PUT("/user/change-password", repository.ChangePasswordProfile(database))

	//News Routes
	ech.GET("/news", repository.GetAllNews(database))
	ech.GET("/news/:id", repository.GetNewsById(database))
	ech.POST("/news", repository.AddNews(database))
	ech.PUT("/news/:id", repository.EditNews(database))
	ech.DELETE("/news/:id", repository.DeleteNews(database))
}
