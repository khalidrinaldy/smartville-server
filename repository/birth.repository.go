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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Birth Data Success", &births))
	}
}

func GetBirthRegistrationById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var birth entity.BirthRegistration
		result := db.First(&birth, c.Param("id"))
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
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

		//Post History Birth
		postHistory, postHistoryErr := AddHistory(
			db,
			user,
			"Pendataan Kelahiran",
			c.FormValue("registration_token"),
		)
		if postHistoryErr != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, postHistory, postHistoryErr.Error()))
		}

		//Post birth registration
		birth.HistoryId,_ = strconv.Atoi(postHistory)
		resultBirth := db.Create(&birth)
		if resultBirth.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultBirth.Error.Error()))
		}

		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Add Birth Registration Success", &birth))
	}
}

func EditBirthRegistration(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		var birth entity.BirthRegistration

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
		birth.Nama_bayi = c.FormValue("nama_bayi")
		birth.Jenis_kelamin,_ = strconv.ParseBool(c.FormValue("jenis_kelamin"))
		birth.Nama_ayah = c.FormValue("nama_ayah")
		birth.Nama_ibu = c.FormValue("nama_ibu")
		birth.Anak_ke,_ = strconv.Atoi(c.FormValue("anak_ke"))
		birth.Tanggal_kelahiran,_ = time.Parse(
			"2006-01-02T15:04:05+0700", 
			fmt.Sprintf("%sT%s+0700", c.FormValue("tgl_lahir"), c.FormValue("waktu_lahir")))
		birth.Alamat_kelahiran = c.FormValue("alamat_kelahiran")

		//PUT request
		resultEdit := db.Model(&birth).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{
			"nama_bayi": birth.Nama_bayi,
			"jenis_kelamin": birth.Jenis_kelamin,
			"nama_ayah": birth.Nama_ayah,
			"nama_ibu": birth.Nama_ibu,
			"anak_ke": birth.Anak_ke,
			"tanggal_kelahiran": birth.Tanggal_kelahiran,
			"alamat_kelahiran": birth.Alamat_kelahiran,
		})
		if resultEdit.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultEdit.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit Birth Registration Data Success", &birth))
	}
}

func DeleteBirthRegistration(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		var birth entity.BirthRegistration

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
		resultDelete := db.Delete(&birth, c.Param("id"))
		if resultDelete.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultDelete.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Delete Birth Registration Data Success", &birth))
	}
}