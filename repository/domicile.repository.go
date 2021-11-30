package repository

import (
	"net/http"
	"smartville-server/entity"
	"smartville-server/helper"
	"strings"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetAllDomicileRegistration(db *gorm.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		var domiciles []entity.DomicileRegistration
		result := db.Find(&domiciles)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Domicile Data Success", &domiciles))
	}
}

func GetDomicileRegistrationById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var domicile entity.DomicileRegistration
		result := db.First(&domicile, c.Param("id"))
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Domicile Registration Id Not Found", result.RowsAffected))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Domicile Data Success", &domicile))
	}
}

func AddDomicileRegistration(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		var domicile entity.DomicileRegistration

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
		domicile.UserNik = user.Nik
		domicile.Nik_pemohon = c.FormValue("nik_pemohon")
		domicile.Nama_pemohon = c.FormValue("nama_pemohon")
		domicile.Tgl_lahir,_ = time.Parse("20060102", c.FormValue("tgl_lahir"))
		domicile.Asal_domisili = c.FormValue("asal_domisili")
		domicile.Tujuan_domisili = c.FormValue("tujuan_domisili")

		//Post Domicile registration
		resultDomicile := db.Create(&domicile)
		if resultDomicile.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultDomicile.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Add Domicile Registration Success", &domicile))
	}
}

func EditDomicileRegistration(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		var domicile entity.DomicileRegistration

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Check Is Admin
		result := db.First(&admin, "token = ?", headerToken)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Admin Token Not Found", ""))
		}
		if admin.Role != "admin" {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Akun Anda Bukan Akun Admin", ""))
		}

		//Get Value From Body
		domicile.Nik_pemohon = c.FormValue("nik_pemohon")
		domicile.Nama_pemohon = c.FormValue("nama_pemohon")
		domicile.Tgl_lahir,_ = time.Parse("20060102", c.FormValue("tgl_lahir"))
		domicile.Asal_domisili = c.FormValue("asal_domisili")
		domicile.Tujuan_domisili = c.FormValue("tujuan_domisili")

		//PUT Request
		resultEdit := db.Model(&domicile).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{
			"nik_pemohon": domicile.Nik_pemohon,
			"nama_pemohon": domicile.Nama_pemohon,
			"tgl_lahir": domicile.Tgl_lahir,
			"asal_domisili": domicile.Asal_domisili,
			"tujuan_domisili": domicile.Tujuan_domisili,
		})
		if resultEdit.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultEdit.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit Domicile Registration Data Success", &domicile))
	}
}

func DeleteDomicileRegistration(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		var domicile entity.DomicileRegistration

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Check Is Admin
		result := db.First(&admin, "token = ?", headerToken)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Admin Token Not Found", ""))
		}
		if admin.Role != "admin" {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Akun Anda Bukan Akun Admin", ""))
		}

		//DELETE
		resultDelete := db.Delete(&domicile, c.Param("id"))
		if resultDelete.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultDelete.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Delete Domicile Registration Data Success", &domicile))
	}
}