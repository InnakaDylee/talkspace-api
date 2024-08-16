package cloud

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"talkspace-api/app/configs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func UploadImageToS3(image *multipart.FileHeader) (string, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		logrus.Error("failed to load configuration:", err)
		return "", err
	}

	awsAccessKeyID := config.CLOUDSTORAGE.AWS_ACCESS_KEY_ID
	awsSecretAccessKey := config.CLOUDSTORAGE.AWS_SECRET_ACCESS_KEY
	awsRegion := config.CLOUDSTORAGE.AWS_REGION
	bucketName := config.CLOUDSTORAGE.AWS_BUCKET_NAME

	maxUploadSize := int64(10 * 1024 * 1024)
	if image.Size > maxUploadSize {
		return "", errors.New("file size exceeds the maximum allowed size of 10MB")
	}

	extension := filepath.Ext(image.Filename)
	allowedExtensions := map[string]bool{".jpg": true, ".png": true, ".jpeg": true}

	if !allowedExtensions[strings.ToLower(extension)] {
		return "", errors.New("invalid image file format. supported formats: .jpg, .jpeg, .png")
	}

	imagePath := uuid.New().String() + extension

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		logrus.Error("failed to create AWS session:", err)
		return "", err
	}

	svc := s3.New(sess)

	file, err := image.Open()
	if err != nil {
		logrus.Error("failed to open file:", err)
		return "", err
	}
	defer file.Close()

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imagePath),
		Body:   file,
	}

	_, err = svc.PutObject(params)
	if err != nil {
		logrus.Error("failed to upload file to S3:", err)
		return "", err
	}

	imageURL := "https://" + bucketName + ".s3.amazonaws.com/" + imagePath

	return imageURL, nil
}
