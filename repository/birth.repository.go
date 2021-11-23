package repository

import (
	"fmt"
	"net/http"
	"smartville-server/entity"
	"smartville-server/helper"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetAllBirthRegistration(db *gorm.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		var births []entity.BirthRegistration
		result := db.Find(&births)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Birth Data Success", &births))
	}
}

func GetBirthRegistrationById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var birth entity.BirthRegistration
		result := db.First(&birth, c.Param("id"))
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "BirthRegistration Id Not Found", result.RowsAffected))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Birth Data Success", &birth))
	}
}

func AddBirthRegistration(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		var birth entity.BirthRegistration

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

		//Get Value From Body
		birth.UserNik = user.Nik
		birth.Nama_bayi = c.FormValue("nama_bayi")
		birth.Jenis_kelamin,_ = strconv.ParseBool(c.FormValue("jenis_kelamin"))
		birth.Nama_ayah = c.FormValue("nama_ayah")
		birth.Nama_ibu = c.FormValue("nama_ibu")
		birth.Anak_ke,_ = strconv.Atoi(c.FormValue("anak_ke"))
		birth.Tanggal_kelahiran,_ = time.Parse(
			"2006-01-02T15:04:05+0700", 
			fmt.Sprintf("%sT%s+0700", c.FormValue("tgl_lahir"), c.FormValue("waktu_lahir")))
		birth.Alamat_kelahiran = c.FormValue("alamat_kelahiran")

		//Post birth registration
		resultBirth := db.Create(&birth)
		if resultBirth.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultBirth.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Add Birth Registration Success", &birth))
	}
}