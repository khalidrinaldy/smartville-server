package repository

import (
	"fmt"
	"net/http"
	"os"
	"smartville-server/entity"
	"smartville-server/helper"
	"strings"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetAdminList(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admins []entity.Admin
		result := db.Raw("SELECT id, nama, email, profile_pic from admins").Scan(&admins)

		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}

		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Admin List Success", &admins))
	}
}

func GetAdminByToken(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Query
		result := db.First(&admin, "token = ?", headerToken)
		if result.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", result.Error))
		}
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Token Not Found", ""))
		}

		admin.Password = "hidden"
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Fetch Data Success", &admin))
	}
}

func RegisterAdmin(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin

		//Get Value From Body
		admin.Email = c.FormValue("email")
		admin.Nama = c.FormValue("nama")
		admin.Password = c.FormValue("password")

		//Check Email
		adminEmail := db.Where("email = ?", admin.Email).Find(&admin)
		if adminEmail.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured", adminEmail.Error))
		}
		if adminEmail.RowsAffected > 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Email Sudah Dipakai", ""))
		}

		//Hashing Password
		hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 5)
		if err != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Hashing Password", err))
		}
		admin.Password = string(hash)

		//Check file upload
		_, photoErr := c.FormFile("profile_pic")
		if photoErr != http.ErrMissingFile {
			//Upload profile picture
			imageURL, err := helper.UploadImage(c, admin.Email, "profile_pic", fmt.Sprintf("admins/%s/profile", admin.Email), "profile")
			if err != nil {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Upload Image", err.Error()))
			}
			admin.Profile_pic = imageURL
		}

		//Token
		admin.Token = helper.JwtGenerator(admin.Email, os.Getenv("SECRET_KEY"))

		//Post Register
		regisResult := db.Create(&admin)
		if regisResult.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", regisResult.Error))
		}
		admin.Password = "hidden"
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Register Success", &admin))
	}
}

func LoginAdmin(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var adminInput entity.Admin
		var adminResult entity.Admin

		//Get Value From Body
		adminInput.Email = c.FormValue("email")
		adminInput.Password = c.FormValue("password")

		//Check User Exist
		resLogin := db.Where("email = ?", adminInput.Email).Find(&adminResult)
		if resLogin.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resLogin.Error))
		}
		if resLogin.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Email Belum Pernah Didaftarkan", ""))
		}

		//Check Password
		checkPass := bcrypt.CompareHashAndPassword([]byte(adminResult.Password), []byte(adminInput.Password))
		if checkPass != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Password Salah", ""))
		}

		//Login Success
		adminResult.Password = "hidden"
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Login Success", &adminResult))
	}
}

func EditProfileAdmin(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		admin.Nama = c.FormValue("nama")
		admin.Email = c.FormValue("email")

		//Get token
		headerToken := c.Request().Header.Get("Authorization")
		headerToken = strings.ReplaceAll(headerToken, "Bearer", "")
		headerToken = strings.ReplaceAll(headerToken, " ", "")

		//Check upload photo
		_, photoErr := c.FormFile("profile_pic")
		if photoErr != http.ErrMissingFile {
			//Upload profile picture
			imageURL, err := helper.UploadImage(c, admin.Email, "profile_pic", fmt.Sprintf("admins/%s/profile", admin.Email), "profile")
			if err != nil {
				return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Upload Image", err.Error()))
			}
			admin.Profile_pic = imageURL
		}

		//Hashing Password
		hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 5)
		if err != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Hashing Password", err))
		}
		admin.Password = string(hash)

		//Put Request
		resultEdit := db.Model(&admin).Where("token = ?", headerToken).Updates(map[string]interface{}{
			"nama": admin.Nama,
			"email": admin.Email,
			"password": admin.Password,
			"profile_pic": admin.Profile_pic,
		})
		if resultEdit.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultEdit.Error))
		}
		if resultEdit.RowsAffected == 0 {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Token not found", ""))
		}
		admin.Password = "hidden"
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Edit Profile Success", &admin))
	}
}

func DeleteAdmin(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var admin entity.Admin
		resultDelete := db.Delete(&admin, c.Param("id"))
		if resultDelete.Error != nil {
			return c.JSON(http.StatusOK, helper.ResultResponse(true, "Error Occured While Querying SQL", resultDelete.Error))
		}
		return c.JSON(http.StatusOK, helper.ResultResponse(false, "Delete Admin Data Success", &admin))
	}
}