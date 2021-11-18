package helper

import (
	"context"
	"io"
	"os"
	"smartville-server/config"

	"github.com/labstack/echo"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploadImage(c echo.Context, formName, folderName, fileName string) (string, error) {
	cld, err := config.CloudConfig()
	if err != nil {
		return "Config Cloud Error", err
	}
	// Source
	file, err := c.FormFile(formName)
	if err != nil {
		return "Form Error", err
	}
	src, err := file.Open()
	if err != nil {
		return "Error While Processing File", err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return "Error While Processing File", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "Error While Processing File", err
	}

	uploadResult, err := cld.Upload.Upload(
		context.Background(),
		file.Filename,
		uploader.UploadParams{
			Folder:           folderName,
			UseFilename:      true,
			UniqueFilename:   false,
			FilenameOverride: fileName,
		},
	)
	if err != nil {
		return "Error While Uploading File", err
	}
	return uploadResult.URL, err
	// return c.JSON(http.StatusOK, map[string]interface{}{
	// 	"message": "Upload success",
	// 	"data":    uploadResult.URL,
	// })
}
