package routes

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"smartville-server/config"
	"smartville-server/db"
	"smartville-server/middleware"
	"smartville-server/repository"
	"github.com/cloudinary/cloudinary-go/api/uploader"

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
	ech.POST("/user/email-verif", repository.SendEmail(database))

	//Change password
	ech.POST("/user/change-password", repository.ChangePassword(database))

	//Test Upload image
	ech.POST("/image", func(c echo.Context) error {
		cld, err := config.CloudConfig()
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "Error occured cloud config",
				"data":    err,
			})
		}
		// Source
		file, err := c.FormFile("image")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		uploadResult, err := cld.Upload.Upload(
			context.Background(),
			file.Filename,
			uploader.UploadParams{
				Folder: "test",
				UseFilename: true,
				//DiscardOriginalFilename: true,
				FilenameOverride: "Asuna",
			},
		)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "Error occured upload image",
				"data":    err.Error(),
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Upload success",
			"data":    uploadResult.URL,
		})
	})
}
