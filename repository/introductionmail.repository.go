package repository

import (
	"net/http"
	"smartville-server/entity"
	"smartville-server/helper"
	"strings"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetAllIntroductionMail(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var introductionMail []entity.IntroductionMail

		//Query
		result := db.Find(&introductionMail)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch IntroductionMail Data Success", &introductionMail))
	}
}

func GetIntroductionMailById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var introductionMail entity.IntroductionMail
		result := db.First(&introductionMail, c.Param("id"))
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "IntroductionMail Id Not Found", result.RowsAffected))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch IntroductionMail Data Success", &introductionMail))
	}
}

func AddIntroductionMail(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		var introductionMail entity.IntroductionMail

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
		introductionMail.UserNik = user.Nik
		introductionMail.Nik_pemohon = c.FormValue("nik_pemohon")
		introductionMail.Nama_pemohon = c.FormValue("nama_pemohon")
		introductionMail.No_hp = c.FormValue("no_hp")
		introductionMail.Alamat_pemohon = c.FormValue("alamat_pemohon")
		introductionMail.Jenis_surat = c.FormValue("jenis_surat")

		//POST Request
		resultAdd := db.Create(&introductionMail)
		if resultAdd.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultAdd.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Add IntroductionMail Success", &introductionMail))
	}
}

func EditIntroductionMail(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		var introductionMail entity.IntroductionMail

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
		introductionMail.Nik_pemohon = c.FormValue("nik_pemohon")
		introductionMail.Nama_pemohon = c.FormValue("nama_pemohon")
		introductionMail.No_hp = c.FormValue("no_hp")
		introductionMail.Alamat_pemohon = c.FormValue("alamat_pemohon")
		introductionMail.Jenis_surat = c.FormValue("jenis_surat")

		//PUT Request
		resultEdit := db.Model(&introductionMail).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{
			"nik_pemohon":    introductionMail.Nik_pemohon,
			"nama_pemohon":   introductionMail.Nama_pemohon,
			"no_hp":          introductionMail.No_hp,
			"alamat_pemohon": introductionMail.Alamat_pemohon,
			"jenis_surat":    introductionMail.Jenis_surat,
		})
		if resultEdit.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultEdit.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit IntroductionMail Data Success", &introductionMail))
	}
}

func DeleteIntroductionMail(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		var introductionMail entity.IntroductionMail

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
		resultDelete := db.Delete(&introductionMail, c.Param("id"))
		if resultDelete.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultDelete.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Delete IntroductionMail Data Success", &introductionMail))
	}
}