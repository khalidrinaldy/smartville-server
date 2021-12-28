package repository

import (
	"net/http"
	"smartville-server/entity"
	"smartville-server/helper"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetAllFinancialHelp(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var financialHelps []entity.FinancialQuery

		//Query
		result := db.Raw(`select fh.id , fh.nama_bantuan , fh.jenis_bantuan , fh.jumlah_dana , fh.alokasi_dana , fh.dana_terealisasi , fh.sisa_dana_bantuan ,fh.history_id, h.status 
		from financial_helps fh 
		left join histories h 
		on fh.history_id = h.id
		order by fh.id desc`).Scan(&financialHelps)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch FinancialHelp Data Success", &financialHelps))
	}
}

func GetFinancialHelpById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var financialHelp entity.FinancialHelp

		//Query
		result := db.First(&financialHelp, c.Param("id"))
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "FinancialHelp Id Not Found", ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch FinancialHelp Data Success", &financialHelp))
	}
}

func AddFinancialHelp(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		var financialHelp entity.FinancialHelp

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
		financialHelp.UserNik = user.Nik
		financialHelp.Nama_bantuan = c.FormValue("nama_bantuan")
		financialHelp.Jenis_bantuan = c.FormValue("jenis_bantuan")
		financialHelp.Jumlah_dana,_ = strconv.Atoi(c.FormValue("jumlah_dana"))
		financialHelp.Alokasi_dana,_ = strconv.Atoi(c.FormValue("alokasi_dana"))
		financialHelp.Dana_terealisasi,_ = strconv.Atoi(c.FormValue("dana_terealisasi"))
		financialHelp.Sisa_dana_bantuan,_ = strconv.Atoi(c.FormValue("sisa_dana_bantuan"))

		//Post History Financial Help
		postHistory, postHistoryErr := AddHistory(
			db,
			user,
			"Permohonan Bantuan Dana",
			c.FormValue("registration_token"),
		)
		if postHistoryErr != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, postHistoryErr.Error(), ""))
		}

		//POST Request
		financialHelp.HistoryId,_ = strconv.Atoi(postHistory)
		resultAdd := db.Create(&financialHelp)
		if resultAdd.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, resultAdd.Error.Error(), ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Add FinancialHelp Success", &financialHelp))
	}
}

func EditFinancialHelp(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		var financialHelp entity.FinancialHelp

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
		financialHelp.Nama_bantuan = c.FormValue("nama_bantuan")
		financialHelp.Jenis_bantuan = c.FormValue("jenis_bantuan")
		financialHelp.Jumlah_dana,_ = strconv.Atoi(c.FormValue("jumlah_dana"))
		financialHelp.Alokasi_dana,_ = strconv.Atoi(c.FormValue("alokasi_dana"))
		financialHelp.Dana_terealisasi,_ = strconv.Atoi(c.FormValue("dana_terealisasi"))
		financialHelp.Sisa_dana_bantuan,_ = strconv.Atoi(c.FormValue("sisa_dana_bantuan"))

		//PUT Request
		resultEdit := db.Model(&financialHelp).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{
			"nama_bantuan":    financialHelp.Nama_bantuan,
			"jenis_bantuan":   financialHelp.Jenis_bantuan,
			"jumlah_dana":          financialHelp.Jumlah_dana,
			"alokasi_dana": financialHelp.Alokasi_dana,
			"dana_terealisasi":    financialHelp.Dana_terealisasi,
			"sisa_dana_bantuan":    financialHelp.Sisa_dana_bantuan,
		})
		if resultEdit.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, resultEdit.Error.Error(), ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit FinancialHelp Data Success", &financialHelp))
	}
}

func DeleteFinancialHelp(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		var financialHelp entity.FinancialHelp

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
		resultDelete := db.Delete(&financialHelp, c.Param("id"))
		if resultDelete.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, resultDelete.Error.Error(), ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Delete FinancialHelp Data Success", &financialHelp))
	}
}
