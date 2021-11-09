package repository

import (
	"net/http"
	"smartville-server/config"
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
		user.Tgl_lahir,_ = time.Parse("2020-05-11", c.FormValue("tgl_lahir"))
		user.Alamat = c.FormValue("alamat")
		user.Dusun = c.FormValue("dusun")
		user.Rt,_ = strconv.Atoi(c.FormValue("rt"))
		user.Rw,_ = strconv.Atoi(c.FormValue("rw"))
		user.Jenis_kelamin,_ = strconv.ParseBool(c.FormValue("jenis_kelamin"))
		user.No_hp = c.FormValue("no_hp")
		user.Profile_pic = c.FormValue("profile_pic")

		//Check user exist
		userExist := db.Where("email = ?", user.Email).Find(&user)
		if userExist.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", userExist.Error))
		}
		if userExist.RowsAffected > 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "email already used", ""))
		}

		//Hashing Password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
		if err != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", err.Error()))
		}

		//Token
		cfg,_ := config.NewConfig(".env")
		user.Password = string(hash)
		user.Token = helper.JwtGenerator(user.Nik, cfg.JWTConfig.SecretKey)

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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Invalid Username", ""))
		}

		//Check Password
		checkPass := bcrypt.CompareHashAndPassword([]byte(userInput.Password), []byte(userResult.Password))
		if checkPass != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Wrong Password", ""))
		}

		//Login Success
		userResult.Password = "hidden"
		return c.JSON(http.StatusOK, &userResult)
	}
}