package repository

import (
	"net/http"
	"os"
	"smartville-server/entity"
	"smartville-server/helper"
	"strconv"
	"time"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)



func GetUserList(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []entity.UserList
		result := db.Raw("SELECT nik, nama from users").Scan(&users)
		if result.Error!=nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch User List Success", &users))
	}
}

func GetUserById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		result := db.First(&user, c.Param("id"))

		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", result.Error))
		}

		header := c.Request().Header.Get("Authorization")
		
		return c.JSON(http.StatusOK, helper.ResultResponse(false, header, &user))
		// return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Data Success", &user))
	}
}

func Register(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		user.Nik = c.FormValue("nik")
		user.Nama = c.FormValue("nama")
		user.Email = c.FormValue("email")
		user.Password = c.FormValue("password")
		user.Tgl_lahir,_ = time.Parse("2020-05-11", c.FormValue("tgl_lahir"))
		user.Alamat = c.FormValue("alamat")
		user.Dusun = c.FormValue("dusun")
		user.Rt,_ = strconv.Atoi(c.FormValue("rt"))
		user.Rw,_ = strconv.Atoi(c.FormValue("rw"))
		user.Jenis_kelamin,_ = strconv.ParseBool(c.FormValue("jenis_kelamin"))
		user.No_hp = c.FormValue("no_hp")
		user.Profile_pic = c.FormValue("profile_pic")

		//Check NIK
		userNik := db.Where("email = ?", user.Email).Find(&user)
		if userNik.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", userNik.Error))
		}
		if userNik.RowsAffected > 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "NIK Sudah Didaftarkan", ""))
		}

		//Check Email
		userEmail := db.Where("email = ?", user.Email).Find(&user)
		if userEmail.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", userEmail.Error))
		}
		if userEmail.RowsAffected > 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Email Sudah Dipakai", ""))
		}

		//Hashing Password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
		if err != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", err.Error()))
		}
		user.Password = string(hash)

		//Token
		user.Token = helper.JwtGenerator(user.Nik, os.Getenv("SECRET_KEY"))

		//Post Register
		regisResult := db.Create(&user)
		if regisResult.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", regisResult.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Register Success", &user))
	}
}

func Login(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var userInput entity.User
		var userResult entity.User
		userInput.Email = c.FormValue("email")
		userInput.Password = c.FormValue("password")

		//Login
		resLogin := db.Where("email = ?", userInput.Email).Find(&userResult)
		if resLogin.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", resLogin.Error))
		}

		//Check user exist
		if resLogin.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Email Belum Pernah Didaftarkan", ""))
		}

		//Check Password
		checkPass := bcrypt.CompareHashAndPassword([]byte(userResult.Password), []byte(userInput.Password))
		if checkPass != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Password Salah", ""))
		}

		//Login Success
		userResult.Password = "hidden"
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Login Success", &userResult))
	}
}