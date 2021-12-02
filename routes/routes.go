package routes

import (
	"net/http"
	"path"
	"smartville-server/db"
	"smartville-server/middleware"
	"smartville-server/repository"
	"github.com/labstack/echo"
)

func InitRoute(ech *echo.Echo) {
	//Init db
	database := db.OpenDatabase()

	//Migrate
	db.Migrate(database)

	//Basic route
	ech.GET(path.Join("/"), func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Welcome to Smartville API")
	})

	//Admin Routes
	ech.GET("/admins", repository.GetAdminList(database))
	ech.GET("/adminbytoken", repository.GetAdminByToken(database), middlewares.IsLoggedIn())
	ech.POST("/admin/register", repository.RegisterAdmin(database))
	ech.POST("/admin/login", repository.LoginAdmin(database))
	ech.PUT("/admin/edit", repository.EditProfileAdmin(database), middlewares.IsLoggedIn())

	//User Routes
	ech.GET("/user-list", repository.GetUserList(database))
	ech.POST("/register", repository.Register(database))
	ech.POST("/login", repository.Login(database))
	ech.GET("/user-id/:id", repository.GetUserById(database), middlewares.IsLoggedIn())
	ech.GET("/userbytoken", repository.GetUserByToken(database), middlewares.IsLoggedIn())
	ech.PUT("/user/edit", repository.EditProfile(database), middlewares.IsLoggedIn())

	//Email verification
	ech.POST("/user/email-verif", repository.SendEmail(database))

	//Change password
	ech.PUT("/user/forgot-password", repository.ChangeForgotPassword(database))
	ech.PUT("/user/change-password", repository.ChangePasswordProfile(database), middlewares.IsLoggedIn())

	//News Routes
	ech.GET("/news", repository.GetAllNews(database))
	ech.GET("/news/:id", repository.GetNewsById(database))
	ech.POST("/news", repository.AddNews(database))
	ech.PUT("/news/:id", repository.EditNews(database))
	ech.DELETE("/news/:id", repository.DeleteNews(database))

	//Birth Registration Routes
	ech.GET("/birth-regis", repository.GetAllBirthRegistration(database))
	ech.GET("/birth-regis/:id", repository.GetBirthRegistrationById(database))
	ech.POST("/birth-regis", repository.AddBirthRegistration(database), middlewares.IsLoggedIn())
	ech.PUT("/birth-regis/:id", repository.EditBirthRegistration(database), middlewares.IsLoggedIn())
	ech.DELETE("/birth-regis/:id", repository.DeleteBirthRegistration(database), middlewares.IsLoggedIn())

	//Domicile Registration Routes
	ech.GET("/domicile-regis", repository.GetAllDomicileRegistration(database))
	ech.GET("/domicile-regis/:id", repository.GetDomicileRegistrationById(database))
	ech.POST("/domicile-regis", repository.AddDomicileRegistration(database), middlewares.IsLoggedIn())
	ech.PUT("/domicile-regis/:id", repository.EditDomicileRegistration(database), middlewares.IsLoggedIn())
	ech.DELETE("/domicile-regis/:id", repository.DeleteDomicileRegistration(database), middlewares.IsLoggedIn())

	//Introduction Mail Registration Routes
	ech.GET("/introductionmail", repository.GetAllIntroductionMail(database))
	ech.GET("/introductionmail/:id", repository.GetIntroductionMailById(database))
	ech.POST("/introductionmail", repository.AddIntroductionMail(database), middlewares.IsLoggedIn())
	ech.PUT("/introductionmail/:id", repository.EditIntroductionMail(database), middlewares.IsLoggedIn())
	ech.DELETE("/introductionmail/:id", repository.DeleteIntroductionMail(database), middlewares.IsLoggedIn())

	//Report Routes
	ech.GET("/report", repository.GetAllReports(database))
	ech.GET("/report/:id", repository.GetReportById(database))
	ech.POST("/report", repository.AddReport(database), middlewares.IsLoggedIn())
	ech.PUT("/report/:id", repository.EditReport(database), middlewares.IsLoggedIn())
	ech.DELETE("/report/:id", repository.DeleteReport(database), middlewares.IsLoggedIn())

	//Financial Help Routes
	ech.GET("/financialhelp", repository.GetAllFinancialHelp(database))
	ech.GET("/financialhelp/:id", repository.GetFinancialHelpById(database))
	ech.POST("/financialhelp", repository.AddFinancialHelp(database), middlewares.IsLoggedIn())
	ech.PUT("/financialhelp/:id", repository.EditFinancialHelp(database), middlewares.IsLoggedIn())
	ech.DELETE("/financialhelp/:id", repository.DeleteFinancialHelp(database), middlewares.IsLoggedIn())

	//Death Data Routes
	ech.GET("/deathdata", repository.GetAllDeathData(database))
	ech.GET("/deathdata/:id", repository.GetDeathDataById(database))
	ech.POST("/deathdata", repository.AddDeathData(database), middlewares.IsLoggedIn())
	ech.PUT("/deathdata/:id", repository.EditDeathData(database), middlewares.IsLoggedIn())
	ech.DELETE("/deathdata/:id", repository.DeleteDeathData(database), middlewares.IsLoggedIn())

	//History Routes
	ech.GET("/history", repository.GetAllHistory(database), middlewares.IsLoggedIn())
	ech.PUT("/history/:history_id", repository.EditStatusHistory(database), middlewares.IsLoggedIn())
}
