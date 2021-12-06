package repository

import (
	"net/http"
	"smartville-server/entity"
	"smartville-server/helper"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetAllDeathData(db *gorm.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		var deaths []entity.Death

		//Query
		result := db.Find(&deaths)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Death Data Success", &deaths))
	}
}

func GetDeathDataById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var death entity.Death

		//Query
		result := db.First(&death, c.Param("id"))
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Death Registration Id Not Found", result.RowsAffected))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Death Data Success", &death))
	}
}

func AddDeathData(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var death entity.Death
		var user entity.User

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Query user first
		result := db.First(&user, "token = ?", headerToken)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "User Token Not Found", ""))
		}

		//Get Value From Body
		death.UserNik = user.Nik
		death.Nik = c.FormValue("nik")
		death.Nama = c.FormValue("nama")
		death.Jenis_kelamin,_ = strconv.ParseBool(c.FormValue("jenis_kelamin"))
		death.Usia,_ = strconv.Atoi(c.FormValue("usia"))
		death.Tgl_wafat,_ = time.Parse("20060102", c.FormValue("tgl_wafat"))
		death.Alamat = c.FormValue("alamat")

		//Post History Death
		postHistory, postHistoryErr := AddHistory(
			db,
			user,
			"Pendataan Kematian",
			c.FormValue("registration_token"),
		)
		if postHistoryErr != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, postHistory, postHistoryErr.Error()))
		}

		//Post Death registration
		death.HistoryId,_ = strconv.Atoi(postHistory)
		resultAdd := db.Create(&death)
		if resultAdd.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultAdd.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Add Death Data Success", &death))
	}
}

func EditDeathData(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var death entity.Death
		var admin entity.Admin

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Check Is Admin
		result := db.First(&admin, "token = ?", headerToken)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Admin Token Not Found", ""))
		}
		if admin.Role != "admin" {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Akun Anda Bukan Akun Admin", ""))
		}

		//Get Value From Body
		death.Nik = c.FormValue("nik")
		death.Nama = c.FormValue("nama")
		death.Jenis_kelamin,_ = strconv.ParseBool(c.FormValue("jenis_kelamin"))
		death.Usia,_ = strconv.Atoi(c.FormValue("usia"))
		death.Tgl_wafat,_ = time.Parse("20060102", c.FormValue("tgl_wafat"))
		death.Alamat = c.FormValue("alamat")

		//PUT Request
		resultEdit := db.Model(&death).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{
			"nik": death.Nik,
			"nama": death.Nama,
			"jenis_kelamin": death.Jenis_kelamin,
			"usia": death.Usia,
			"tgl_wafat": death.Tgl_wafat,
			"alamat": death.Alamat,
		})
		if resultEdit.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultEdit.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit Death Data Success", &death))
	}
}

func DeleteDeathData(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var death entity.Death
		var admin entity.Admin

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Check Is Admin
		result := db.First(&admin, "token = ?", headerToken)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Admin Token Not Found", ""))
		}
		if admin.Role != "admin" {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Akun Anda Bukan Akun Admin", ""))
		}

		//DELETE
		resultDelete := db.Delete(&death, c.Param("id"))
		if resultDelete.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultDelete.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Delete Death Data Success", &death))
	}
}