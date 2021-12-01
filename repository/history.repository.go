package repository

import (
	"net/http"
	"smartville-server/entity"
	"smartville-server/helper"
	"strings"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetAllHistory(db *gorm.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		var history []entity.History
		var user entity.User

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Query user first
		result := db.First(&user, "token = ?", headerToken)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "User Token Not Found", ""))
		}

		//Query
		resultHistory := db.Where("user_nik = ?", user.Nik).Find(&history)
		if resultHistory.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultHistory.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch History Data Success", &history))
	}
}