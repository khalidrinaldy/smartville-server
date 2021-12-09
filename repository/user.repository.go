package repository

import (
	"fmt"
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
		if result.Error != nil {
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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
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

func GetUserByToken(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Query
		result := db.First(&user, "token = ?", headerToken)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Token Not Found", ""))
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
		user.Tgl_lahir, _ = time.Parse("20060102", c.FormValue("tgl_lahir"))
		user.Tempat_lahir = c.FormValue("tempat_lahir")
		user.Alamat = c.FormValue("alamat")
		user.Dusun = c.FormValue("dusun")
		user.Rt, _ = strconv.Atoi(c.FormValue("rt"))
		user.Rw, _ = strconv.Atoi(c.FormValue("rw"))
		user.Jenis_kelamin, _ = strconv.ParseBool(c.FormValue("jenis_kelamin"))
		user.No_hp = c.FormValue("no_hp")

		//Check NIK
		userNik := db.Where("nik = ?", user.Nik).Find(&user)
		if userNik.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, userNik.Error.Error(), ""))
		}
		if userNik.RowsAffected > 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "NIK Sudah Didaftarkan", ""))
		}

		//Check Email
		userEmail := db.Where("email = ?", user.Email).Find(&user)
		if userEmail.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, userEmail.Error.Error(), ""))
		}
		if userEmail.RowsAffected > 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Email Sudah Dipakai", ""))
		}

		//Hashing Password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
		if err != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, err.Error(), ""))
		}
		user.Password = string(hash)

		//Token
		user.Token = helper.JwtGenerator(user.Nik, os.Getenv("SECRET_KEY"))

		//Check file upload
		_, photoErr := c.FormFile("profile_pic")
		if photoErr != http.ErrMissingFile {
			//Upload profile picture
			imageURL, err := helper.UploadImage(c, user.Nik, "profile_pic", fmt.Sprintf("users/%s/profile", user.Nik), "profile")
			if err != nil {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, err.Error(), ""))
			}
			user.Profile_pic = imageURL

			//Post Register
			regisResult := db.Create(&user)
			if regisResult.Error != nil {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, regisResult.Error.Error(), ""))
			}
			user.Password = "hidden"
			return c.JSON(http.StatusOK, helper.ResultResponse(false, "Register Success", &user))
		}

		//Post Register Without Photo
		regisResult := db.Create(&user)
		if regisResult.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, regisResult.Error.Error(), ""))
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
			return c.JSON(http.StatusOK, helper.ResultResponse(true, resLogin.Error.Error(), ""))
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

func ChangeForgotPassword(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		email := c.FormValue("email")
		password := c.FormValue("new_password")

		//Check Email
		userEmail := db.Where("email = ?", email).Find(&user)
		if userEmail.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, userEmail.Error.Error(), ""))
		}
		if userEmail.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Email Belum Terdaftar", ""))
		}

		//Hashing Password
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, err.Error(), ""))
		}
		password = string(hash)

		//Change Password
		setPassword := db.Exec("UPDATE users SET password = ? where email = ?", password, email)
		if setPassword.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, setPassword.Error.Error(), ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Password berhasil diubah", ""))
	}
}

func EditProfile(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User
		var userCheck entity.User
		user.Nama = c.FormValue("nama")
		user.Alamat = c.FormValue("alamat")
		user.No_hp = c.FormValue("no_hp")
		user.Email = c.FormValue("email")
		user.Jenis_kelamin, _ = strconv.ParseBool(c.FormValue("jenis_kelamin"))

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Check Email
		userEmail := db.Where("email = ?", user.Email).Find(&userCheck)
		if userEmail.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, userEmail.Error.Error(), ""))
		}
		if userEmail.RowsAffected > 0 {
			if userCheck.Token != headerToken {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, "Email sudah digunakan", ""))
			}
		}

		//Check upload photo
		_, photoErr := c.FormFile("profile_pic")
		if photoErr != http.ErrMissingFile {
			//Upload profile picture
			imageURL, err := helper.UploadImage(c, user.Nik, "profile_pic", fmt.Sprintf("users/%s/profile", user.Nik), "profile")
			if err != nil {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, err.Error(), ""))
			}
			user.Profile_pic = imageURL

			//Update profile
			resultEdit := db.Model(&user).Where("token = ?", headerToken).Updates(map[string]interface{}{
				"nama":     user.Nama,
				"alamat":   user.Alamat,
				"no_hp":    user.No_hp,
				"email":    user.Email,
				"jenis_kelamin": user.Jenis_kelamin,
				"profile_pic": user.Profile_pic,
			})
			if resultEdit.Error != nil {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, resultEdit.Error.Error(), ""))
			}
			if resultEdit.RowsAffected == 0 {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, "Token not found", ""))
			}
			return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit Profile Success", &user))
		}

		//Update profile without photo
		resultEdit := db.Model(&user).Where("token = ?", headerToken).Updates(map[string]interface{}{
			"nama":     user.Nama,
			"alamat":   user.Alamat,
			"no_hp":    user.No_hp,
			"email":    user.Email,
			"jenis_kelamin": user.Jenis_kelamin,
		})
		if resultEdit.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, resultEdit.Error.Error(), ""))
		}
		if resultEdit.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Token not found", ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit Profile Success", &user))
	}
}

// func CheckPassword(db *gorm.DB) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var user entity.User

// 		//Get token
// 		headerToken := c.Request().Header.Get("Authorization")
// 		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
// 		headerToken = strings.ReplaceAll(headerToken, " ", "")

// 		//Get user first
// 		result := db.First(&user, "token = ?", headerToken)
// 		if result.Error != nil {
// 			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
// 		}
// 		if result.RowsAffected == 0 {
// 			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Token Not Found", ""))
// 		}

// 		//Compare password
// 		checkPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.FormValue("password")))
// 		if checkPass != nil {
// 			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Password Salah", ""))
// 		}
// 		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Password benar", map[string]interface{}{
// 			"check": true,
// 		}))
// 	}
// }

func ChangePasswordProfile(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user entity.User

		//Value body new password
		oldPassword := c.FormValue("old_password")
		newPassword := c.FormValue("new_password")

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Get user first
		result := db.First(&user, "token = ?", headerToken)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, result.Error.Error(), ""))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Token Not Found", ""))
		}

		//Compare password
		checkPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
		if checkPass != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Password Salah", ""))
		}

		//Hashing Password
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 5)
		if err != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, err.Error(), ""))
		}
		newPassword = string(hash)

		//Change Password
		setPassword := db.Exec("UPDATE users SET password = ? where token = ?", newPassword, headerToken)
		if setPassword.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, setPassword.Error.Error(), ""))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Password berhasil diubah", ""))
	}
}
