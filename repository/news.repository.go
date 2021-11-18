package repository

import (
	"net/http"
	"smartville-server/entity"
	"smartville-server/helper"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetAllNews(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var news []entity.News
		result := db.Raw("select * from news order by tanggal_terbit limit 5").Scan(&news)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error While Querying SQL", result.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch News Success", &news))
	}
}

func GetNewsById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var news entity.News
		result := db.First(&news, c.Param("id"))
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error While Querying SQL", result.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch News By Id Success", &news))
	}
}

func AddNews(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var news entity.News
		news.Judul_berita = c.FormValue("judul_berita")
		news.Deskripsi_berita = c.FormValue("deskripsi_berita")
		news.Tanggal_terbit = time.Now()

		//Upload news image
		imageURL, err := helper.UploadImage(c, news.Judul_berita, "foto_berita", "news", news.Judul_berita)
		if err != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", err.Error()))
		}
		news.Foto_berita = imageURL

		//Post news
		resultAdd := db.Create(&news)
		if resultAdd.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error While Querying SQL", resultAdd.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Add News Success", &news))
	}
}

func EditNews(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var news entity.News
		news.Judul_berita = c.FormValue("judul_berita")
		news.Deskripsi_berita = c.FormValue("deskripsi_berita")
		news.Tanggal_terbit = time.Now()
		_, err := c.FormFile("foto_berita")

		//check file empty or not
		if err != http.ErrMissingFile {
			//Upload news image
			imageURL, err := helper.UploadImage(c, news.Judul_berita, "foto_berita", "news", news.Judul_berita)
			if err != nil {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", err.Error()))
			}
			news.Foto_berita = imageURL

			//PUT request
			resultEdit := db.Model(&news).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{
				"judul_berita": news.Judul_berita,
				"deskripsi_berita": news.Deskripsi_berita,
				"foto_berita": news.Foto_berita,
			})
			if resultEdit.Error != nil {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error While Querying SQL", resultEdit.Error))
			}
			return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit News Success", &news))
		}

	
		//PUT request without photo
		resultEdit := db.Model(&news).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{
			"judul_berita": news.Judul_berita,
			"deskripsi_berita": news.Deskripsi_berita,
			"tanggal_terbit": news.Tanggal_terbit,
		})
		if resultEdit.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error While Querying SQL", resultEdit.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit News Success", &news))
	}
}

func DeleteNews(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var news entity.News
		resultDelete := db.Delete(&news, c.Param("id"))
		if resultDelete.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error While Querying SQL", resultDelete.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Delete News Success", &news))
	}
}
