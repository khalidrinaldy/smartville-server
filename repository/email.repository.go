package repository

import (
	"math/rand"
	"net/http"
	"smartville-server/helper"
	"strings"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func SendEmail(db *gorm.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		to := []string{c.FormValue("to")}
		subject := c.FormValue("subject")

		otp := []string{}
		for i := 0; i < 4; i++ {
			number := rand.Intn(9-0+1) + 0
			otp = append(otp, string(number))
		}
		message := strings.Join(otp, "")

		sendErr := helper.SendEmail(to, subject, message)
		if sendErr != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error While Sending Email", ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Email Has Been Sent!", message))
	}
}