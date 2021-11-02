package repository

import (
	"net/http"
	"smartville-server/entity"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetAdminList(db *gorm.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		var admins []entity.Admin
		result := db.Find(&admins)
		
		if result.Error != nil {
			return result.Error
		}

		return c.JSON(http.StatusOK, &admins)
	}
}

func GetAdminById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		result := db.First(&admin, c.Param("id"))

		if result.Error != nil {
			return result.Error
		}

		return c.JSON(http.StatusOK, &admin)
	}
}