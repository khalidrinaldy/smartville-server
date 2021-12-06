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

func GetAllReports(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var reports []entity.Report
		
		//Query
		result := db.Find(&reports)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Reports Data Success", &reports))
	}
}

func GetReportById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var report entity.Report

		//Query
		result := db.First(&report, c.Param("id"))
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Report Id Not Found", result.RowsAffected))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Report Data Success", &report))
	}
}

func AddReport(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		var report entity.Report

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
		report.UserNik = user.Nik
		report.Nama_pelapor = c.FormValue("nama_pelapor")
		report.Deskripsi = c.FormValue("deskripsi")
		report.Jenis_laporan = c.FormValue("jenis_laporan")
		report.Tgl_laporan,_ = time.Parse("20060102", c.FormValue("tgl_laporan"))
		report.No_hp = c.FormValue("no_hp")
		report.Alamat = c.FormValue("alamat")

		//Check file upload
		_, photoErr := c.FormFile("foto_kejadian")
		if photoErr != http.ErrMissingFile {
			//Upload profile picture
			imageURL, err := helper.UploadImage(c, "", "foto_kejadian", fmt.Sprintf("laporan/%v/%v", report.Tgl_laporan.Year(), report.Tgl_laporan.Month()), "profile")
			if err != nil {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Upload Image", err.Error()))
			}
			report.Foto_kejadian = imageURL
		}

		//Post History Report
		postHistory, postHistoryErr := AddHistory(
			db,
			user,
			"Pelaporan Kejadian",
			c.FormValue("registration_token"),
		)
		if postHistoryErr != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, postHistory, postHistoryErr.Error()))
		}

		//POST Request
		report.HistoryId,_ = strconv.Atoi(postHistory)
		resultAdd := db.Create(&report)
		if resultAdd.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultAdd.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Add Report Success", &report))
	}
}

func EditReport(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		var report entity.Report

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
		report.Nama_pelapor = c.FormValue("nama_pelapor")
		report.Deskripsi = c.FormValue("deskripsi")
		report.Jenis_laporan = c.FormValue("jenis_laporan")
		report.Tgl_laporan,_ = time.Parse("20060102", c.FormValue("tgl_laporan"))
		report.No_hp = c.FormValue("no_hp")
		report.Alamat = c.FormValue("alamat")

		//Check file upload
		_, photoErr := c.FormFile("foto_kejadian")
		if photoErr != http.ErrMissingFile {
			//Upload profile picture
			imageURL, err := helper.UploadImage(c, "", "foto_kejadian", fmt.Sprintf("laporan/%v/%v", report.Tgl_laporan.Year(), report.Tgl_laporan.Month()), "profile")
			if err != nil {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Upload Image", err.Error()))
			}
			report.Foto_kejadian = imageURL
		}

		//PUT Request
		resultEdit := db.Model(&report).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{
			"nama_pelapor": report.Nama_pelapor,
			"deskripsi": report.Deskripsi,
			"jenis_laporan": report.Jenis_laporan,
			"tgl_laporan": report.Tgl_laporan,
			"no_hp": report.No_hp,
			"alamat": report.Alamat,
		})
		if resultEdit.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultEdit.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit Report Data Success", &report))
	}
}

func DeleteReport(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		var report entity.Report

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
		resultDelete := db.Delete(&report, c.Param("id"))
		if resultDelete.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultDelete.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Delete Report Data Success", &report))
	}
}