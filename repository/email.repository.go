package repository

import (
	"math/rand"
	"net/http"
	"smartville-server/entity"
	"smartville-server/helper"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func SendEmail(db *gorm.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		to := []string{c.FormValue("to")}
		subject := c.FormValue("subject")

		otp := []string{}
		for i := 0; i < 4; i++ {
			rand.Seed(time.Now().UnixNano())
			number := rand.Intn(9-0+1) + 0
			otp = append(otp, strconv.Itoa(number))
		}
		messageOTP := strings.Join(otp, "")

		//Check Email
		var user entity.User
		userEmail := db.Where("email = ?", to).Find(&user)
		if userEmail.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", userEmail.Error))
		}
		if userEmail.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Email Belum Terdaftar", ""))
		}

		//Send Email
		sendErr := helper.SendEmail(to, subject, "Kode OTP kamu : " + messageOTP)
		if sendErr != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error While Sending Email", ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Email Has Been Sent!", map[string]string{
			"otp" : messageOTP,
			"email" : user.Email,
		}))
	}
}