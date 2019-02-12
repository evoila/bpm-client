package s3

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/evoila/BPM-Client/model"
)

func UploadFile(path, depth string, body S3Permission) error {

	file, err := os.Open(path)
	if err != nil {
		return errors.New("Failed to open file " + path + " due to '" + err.Error() + "'")
	}
	defer file.Close()

	uploadCredentials := credentials.NewStaticCredentials(
		body.AuthKey, body.AuthSecret, body.SessionToken)

	s3Session, err := session.NewSession(&aws.Config{
		Region:      aws.String(body.Region),
		Credentials: uploadCredentials,}, )

	if err != nil {
		return errors.New("Unable to create a S3 s3Session due to due to '" + err.Error() + "'")
	}
	var uploader = s3manager.NewUploader(s3Session)

	var done = false

	go func() {
		_, ur := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(body.Bucket),
			Key:    aws.String(body.S3location),
			Body:   file,
		})

		if ur != nil {
			panic(ur)
		}
		done = true
	}()

	var chars = []string{"-", "/", "|", "\\"}
	var line = depth
	var i = 0
	for !done {
		fmt.Printf("\033[%dD", 1)
		fmt.Print(line)

		time.Sleep(time.Second / 2)

		line = depth + chars[i]

		i = (i + 1) % 4

		fmt.Print("\r\033[K")
	}

	if err != nil {
		return errors.New("Failed to upload to S3 due to '" + err.Error() + "'")
	}

	return nil
}

func DownloadFile(filename, depth string, body S3Permission) error {

	file, err := os.Create(filename + ".bpm")
	if err != nil {
		return errors.New("Failed to create file " + filename + "due to '" + err.Error() + "'")
	}
	defer file.Close()

	downloadCredentials := credentials.NewStaticCredentials(body.AuthKey, body.AuthSecret, body.SessionToken)

	downloadSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(body.Region),
		Credentials: downloadCredentials},
	)

	if err != nil {
		return errors.New("Unable to create a S3 downloadSession due to '" + err.Error() + "'")
	}

	var downloader = s3manager.NewDownloader(downloadSession)
	var done = false

	go func() {
		_, err = downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(body.Bucket),
				Key:    aws.String(body.S3location),
			})
		done = true
	}()

	var chars = []string{"-", "/", "|", "\\"}
	var line = depth
	var i = 0
	for !done {
		fmt.Printf("\033[%dD", 1)
		fmt.Print(line)

		time.Sleep(time.Second / 2)

		line = depth + chars[i]

		i = (i + 1) % 4

		fmt.Print("\r\033[K")
	}

	if err != nil {
		return errors.New("Failed to download the file " + filename + "  due to '" + err.Error() + "'")
	}

	return nil
}
