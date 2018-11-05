package s3

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/evoila/BPM-Client/model"
)

func UploadFile(path string, body S3Permission) error {

	file, err := os.Open(path)
	if err != nil {
		return errors.New("Failed to open file " + path + " due to '" + err.Error() + "'")
	}
	defer file.Close()

	// -- Creating session, service client and uploader --
	os.Setenv("AWS_ACCESS_KEY_ID", body.AuthKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", body.AuthSecret)

	//Clear credentials after use
	defer os.Setenv("AWS_ACCESS_KEY_ID", "")
	defer os.Setenv("AWS_SECRET_ACCESS_KEY", "")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(body.Region)},
	)

	if err != nil {
		return errors.New("Unable to create a S3 session due to due to '" + err.Error() + "'")
	}

	var uploader = s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(body.Bucket),
		Key:    aws.String(body.S3location),
		Body:   file,
	})

	if err != nil {
		return errors.New("Failed to upload to S3 due to '" + err.Error() + "'")
	}

	return nil
}

func DownloadFile(filename string, body S3Permission) error {

	log.Println("Creating file at", filename+".bpm")
	file, err := os.Create(filename + ".bpm")
	if err != nil {
		return errors.New("Failed to create file " + filename + "due to '" + err.Error() + "'")
	}
	defer file.Close()

	// -- Creating downloadSession, service client and uploader --
	os.Setenv("AWS_ACCESS_KEY_ID", body.AuthKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", body.AuthSecret)

	//Clear credentials after use
	defer os.Setenv("AWS_ACCESS_KEY_ID", "")
	defer os.Setenv("AWS_SECRET_ACCESS_KEY", "")

	downloadSession, err := session.NewSession(&aws.Config{
		Region: aws.String(body.Region)},
	)

	if err != nil {
		return errors.New("Unable to create a S3 downloadSession due to '" + err.Error() + "'")
	}

	var downloader = s3manager.NewDownloader(downloadSession)
	//var client = s3.New(downloadSession)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(body.Bucket),
			Key:    aws.String(body.S3location),
		})
	if err != nil {
		return errors.New("Failed to download the file " + filename + "  due to '" + err.Error() + "'")
	}

	log.Println("Successfully downloaded", file.Name(), "(", numBytes, "bytes )")

	return nil
}