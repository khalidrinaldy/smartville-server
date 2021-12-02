package repository

import (
	"fmt"
	"net/http"
	"smartville-server/entity"
	"smartville-server/helper"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetAllHistory(db *gorm.DB) echo.HandlerFunc {
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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error.Error()))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "User Token Not Found", ""))
		}

		//Query
		resultHistory := db.Order("updated_at desc").Where("user_nik = ?", user.Nik).Find(&history)
		if resultHistory.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultHistory.Error.Error()))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch History Data Success", &history))
	}
}

func AddHistory(db *gorm.DB, user entity.User, perihal string, registration_token string) (string, error) {
	var history entity.History

	//Get Value From Body
	history.UserNik = user.Nik
	history.Perihal = perihal
	history.Status = "terkirim"
	history.Deskripsi = fmt.Sprintf("Formulir %s telah berhasil %s", history.Perihal, history.Status)
	history.Registration_token = registration_token

	//Query
	resultAdd := db.Create(&history)
	if resultAdd.Error != nil {
		return "Error Occured While Querying SQL", resultAdd.Error
	}

	//Send Notification
	sendNotif, err := helper.SendNotification(history.Registration_token, history.Deskripsi)
	if err != nil {
		return sendNotif,err
	}

	return strconv.Itoa(history.Id),nil
}

// func AddHistory2(db *gorm.DB, user entity.User, perihal string, status string) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var history entity.History

// 		//Get Value From Body
// 		history.UserNik = user.Nik
// 		history.Perihal = perihal
// 		history.Status = "terkirim"
// 		history.Deskripsi = fmt.Sprintf("Formulir %s telah berhasil %s", history.Perihal, history.Status)
// 		history.Registration_token = c.FormValue("registration_token")

// 		//Send Notification
// 		sendNotif, err := helper.SendNotification(history.Registration_token, history.Deskripsi)
// 		if err != nil {
// 			return c.JSON(http.StatusOK, helper.ResultResponse(true, sendNotif, err))
// 		}

// 		//Query
// 		resultAdd := db.Create(&history)
// 		if resultAdd.Error != nil {
// 			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultAdd.Error))
// 		}
// 		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Add History Data Success", &history))
// 	}
// }

func EditStatusHistory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var history entity.History
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
		history.Id, _ = strconv.Atoi(c.Param("history_id"))

		//Get history by id
		resHistory := db.First(&history, history.Id)
		if resHistory.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resHistory.Error.Error()))
		}
		if resHistory.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Invalid History Id", ""))
		}

		//History Deskripsi
		history.Status = c.FormValue("status")
		if history.Status == "ditolak" {
			history.Deskripsi = fmt.Sprintf("Formulir %s telah %s, mohon datang ke kantor untuk penjelasan hal ini", history.Perihal, history.Status)
		} else {
			history.Deskripsi = fmt.Sprintf("Formulir %s telah berhasil %s, silahkan datang ke kantor untuk mengambil berkas", history.Perihal, history.Status)
		}

		//Query
		resultEdit := db.Model(&history).Where("id = ?", history.Id).Updates(map[string]interface{}{
			"deskripsi": history.Deskripsi,
			"status":    history.Status,
		})
		if resultEdit.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultEdit.Error.Error()))
		}

		//Send Notification
		sendNotif, err := helper.SendNotification(history.Registration_token, history.Deskripsi)
		if err != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, sendNotif, err.Error()))
		}

		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit Status History Success", map[string]interface{}{
			"deskripsi": history.Deskripsi,
			"status":    history.Status,
		}))
	}
}
