package repository

import (
	"log"
	"net/http"
	"smartville-server/entity"
	"smartville-server/helper"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetAllIntroductionMail(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var introductionMail []entity.IntroductionMailQuery

		//Query
		result := db.Raw(`select im.id , im.nik_pemohon , im.nama_pemohon , im.alamat_pemohon , im.no_hp , im.jenis_surat,h.status 
		from introduction_mails im 
		left join histories h 
		on im.history_id = h.id
		order by im.id desc`).Scan(&introductionMail)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch IntroductionMail Data Success", &introductionMail))
	}
}

func GetIntroductionMailById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var introductionMail entity.IntroductionMail
		result := db.First(&introductionMail, c.Param("id"))
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "IntroductionMail Id Not Found", ""))
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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
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

		//Post History Introduction Mail
		postHistory, postHistoryErr := AddHistory(
			db,
			user,
			"Surat Pengantar",
			c.FormValue("registration_token"),
		)
		if postHistoryErr != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, postHistoryErr.Error(), ""))
		}

		//POST Request
		introductionMail.HistoryId,_ = strconv.Atoi(postHistory)
		resultAdd := db.Create(&introductionMail)
		if resultAdd.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, resultAdd.Error.Error(), ""))
		}
		log.Printf("token : %s", c.FormValue("registration_token"))
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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, resultEdit.Error.Error(), ""))
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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, resultDelete.Error.Error(), ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Delete IntroductionMail Data Success", &introductionMail))
	}
}
