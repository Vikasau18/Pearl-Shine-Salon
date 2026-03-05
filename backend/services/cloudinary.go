package services

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService(cloudinaryURL string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}
	return &CloudinaryService{cld: cld}, nil
}

// UploadImage uploads a file to Cloudinary and returns the secure URL
func (s *CloudinaryService) UploadImage(ctx context.Context, file multipart.File, filename string) (string, error) {
	resp, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   "salons",
		PublicID: filename,
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}
