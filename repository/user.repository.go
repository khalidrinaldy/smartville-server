package repository

import (
	"net/http"
	"os"
	"smartville-server/entity"
	"smartville-server/helper"
	"strconv"
	"strings"
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

		//Check token is valid for user id
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")
		if user.Token != headerToken {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Invalid Token", ""))
		}
		
		user.Password = "hidden"
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Data Success", &user))
	}
}

func Register(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		user.Nik = c.FormValue("nik")
		user.Nama = c.FormValue("nama")
		user.Email = c.FormValue("email")
		user.Password = c.FormValue("password")
		user.Tgl_lahir,_ = time.Parse("20060102", c.FormValue("tgl_lahir"))
		user.Tempat_lahir = c.FormValue("tempat_lahir")
		user.Alamat = c.FormValue("alamat")
		user.Dusun = c.FormValue("dusun")
		user.Rt,_ = strconv.Atoi(c.FormValue("rt"))
		user.Rw,_ = strconv.Atoi(c.FormValue("rw"))
		user.Jenis_kelamin,_ = strconv.ParseBool(c.FormValue("jenis_kelamin"))
		user.No_hp = c.FormValue("no_hp")
		user.Profile_pic = c.FormValue("profile_pic")

		//Check NIK
		userNik := db.Where("nik = ?", user.Nik).Find(&user)
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
		user.Password = "hidden"
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

func ChangePassword(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		email := c.FormValue("email")
		password := c.FormValue("new_password")

		//Check Email
		userEmail := db.Where("email = ?", email).Find(&user)
		if userEmail.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", userEmail.Error))
		}
		if userEmail.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Email Belum Terdaftar", ""))
		}

		//Hashing Password
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", err.Error()))
		}
		password = string(hash)

		//Change Password
		setPassword := db.Exec("UPDATE users SET password = ? where email = ?", password, email)
		if setPassword !=nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured[setPassword]", setPassword.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Password berhasil diubah", ""))
	}
}