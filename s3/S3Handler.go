package s3

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/evoila/BPM-Client/model"
)

func UploadFile(filename, path string, body model.Destination) error {

	log.Println("Opening file at", path)
	file, err := os.Open(path)
	if err != nil {
		return errors.New("Failed to open file " + path + " due to '" + err.Error() + "'")
	}
	defer file.Close()
	log.Println("Successfully opened file at", path)

	// -- Creating session, service client and uploader --
	os.Setenv("AWS_ACCESS_KEY_ID", body.AuthKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", body.AuthSecret)

	//Clear credentials after use
	defer os.Setenv("AWS_ACCESS_KEY_ID", "")
	defer os.Setenv("AWS_SECRET_ACCESS_KEY", "")

	log.Println("Creating S3 session ...")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(body.Region)},
	)

	if err != nil {
		return errors.New("Unable to create a S3 session due to due to '" + err.Error() + "'")
	}
	log.Println("Successfully created S3 session")

	log.Println("Setting up S3 uploader")
	var uploader = s3manager.NewUploader(sess)
	var client = s3.New(sess)

	// -- Listing all objects of the given bucket --
	// Surely not needed later
	listObjectsOfBucket(body.Bucket, client)

	// -- Uploading the backup file to the given bucket --
	log.Println("Uploading", filename, "to", body.Bucket)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(body.Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})

	if err != nil {
		return errors.New("Failed to upload to S3 due to '" + err.Error() + "'")
	}
	log.Printf("Successfully uploaded %q to %q\n", filename, body.Bucket)

	listObjectsOfBucket(body.Bucket, client)

	return nil
}

func listObjectsOfBucket(bucket string, client *s3.S3) error {

	resp, err := client.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})
	if err != nil {
		return errors.New("Failed to list all buckets due to '" + err.Error() + "'")
	}

	log.Println("Objects of bucket", bucket, ":")

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")

	}

	return nil
}

func listAllBuckets(client *s3.S3) error {
	// -- Listing all buckets --
	// Could be removed later on

	log.Println("Sending request for the bucket list.")
	result, err := client.ListBuckets(nil)
	if err != nil {
		return errors.New("Unable to list buckets due to '" + err.Error() + "'")
	}
	log.Println("Listing all buckets.")
	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
	return nil
}
