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
		result := db.Raw("SELECT * FROM admins").Scan(&admins)
		
		if result.Error != nil {
			return result.Error
		}

		return c.JSON(http.StatusOK, &admins)
	}
}

func GetPostListJoinAdmin(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var results []entity.AdminsPost
		result := db.Raw(
		`select admins.id, admins.nama, posts.title
		from admins
		left join posts
		on posts.admin_id = admins.id;`).Scan(&results)

		if result.Error != nil {
			return result.Error
		}

		return c.JSON(http.StatusOK, &results)
	}
}