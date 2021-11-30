package config

import (
	"os"

	"github.com/cloudinary/cloudinary-go"
)

func CloudConfig() (*cloudinary.Cloudinary, error) {
    cld, err := cloudinary.NewFromParams(
        os.Getenv("CLOUD_NAME"), 
        os.Getenv("CLOUD_API_KEY"), 
        os.Getenv("CLOUD_SECRET_KEY"),
    )
    return cld,err
}