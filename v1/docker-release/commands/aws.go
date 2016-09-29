package commands

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// s3URL used to define the release URL endpoints
var s3URL, s3DomainURL string

// Create constants for all of the release endpoints
const (
	s3Region            = "us-east-1"
	s3StageBucket       = "editions-stage-us-east-1-150610-005505"
	s3LiveBucket        = "editions-us-east-1"
	s3DomainTemplate    = "https://%s.s3.amazonaws.com/%s"
	stageDomainTemplate = "https://download-stage.docker.com/%s"
	liveDomainTemplate  = "https://download.docker.com/%s"
	s3Profile           = "editions"
)

var s3Sess = &session.Session{}

// getS3Session returns an existing session or creates a new one
func getS3Session() *session.Session {
	if s3Sess.Config == nil {
		awsOpts := &aws.Config{Region: aws.String(awsRegion)}
		logrus.Debugf("Key: %s - Region: %s", awsKey, awsRegion)
		if awsKey == "" {
			awsOpts.Credentials = credentials.NewSharedCredentials("", awsProfile)
		} else {
			awsOpts.Credentials = credentials.NewStaticCredentials(awsKey, awsSecret, "")
		}
		s3Sess = session.New(awsOpts)
		_, err := s3Sess.Config.Credentials.Get()
		if err != nil {
			logrus.Fatal("Could not create S3 session: ", err)
		}
	}
	return s3Sess
}

func getObject(objKey string) *s3.GetObjectOutput {
	logrus.Debug("Checking: ", objKey)
	sess := getS3Session()
	svc := s3.New(sess)
	params := &s3.GetObjectInput{
		Bucket: aws.String(awsBucket), // Required
		Key:    aws.String(objKey),    // Required
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		logrus.Error("Build not found!")
		logrus.Error(err.Error())
	}

	logrus.Debug("Object: ", resp)
	return resp
}

// pushData push matching file to bucket from dir - should return success of failure.
func pushData(filePath string, destination string, public bool) string {
	logrus.Debugf("Uploading %s to %s make public? %t", filePath, destination, public)
	sess := getS3Session()
	uploader := s3manager.NewUploader(sess)
	s3UploadInput := &s3manager.UploadInput{
		Bucket: aws.String(awsBucket),
		Key:    aws.String(destination),
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader([]byte("")),
	}

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			logrus.Fatal("Failed to open file ", err)
		}
		defer file.Close()
		s3UploadInput.Body = file
		switch path.Base(destination) {
		case "appcast.xml":
			logrus.Debug("Setting XML content-type")
			s3UploadInput.ContentType = aws.String("application/xml; charset=utf-8")
		case "NOTES", "Docker.dmg.sha256sum", "InstallDocker.msi.sha256sum":
			logrus.Debug("Setting TEXT content-type")
			s3UploadInput.ContentType = aws.String("text/plain; charset=utf-8")
		}
	}
	result, err := uploader.Upload(s3UploadInput)
	if err != nil {
		logrus.Fatal("Failed to upload ", err)
	}
	logrus.Debug(result)
	uploadURL, _ := url.QueryUnescape(result.Location)
	// swap bucket url with live or stage url
	logrus.Debugf("Uploaded to: %s", uploadURL)
	fmt.Println("Successfully uploaded to ", s3URL+destination)
	return s3URL + destination
}

func createLatest(srcKey string, targetKey string) *s3.CopyObjectOutput {
	sess := getS3Session()
	svc := s3.New(sess)
	params := &s3.CopyObjectInput{
		Bucket:     aws.String(awsBucket),                // Required
		Key:        aws.String(targetKey),                // Required
		CopySource: aws.String(awsBucket + "/" + srcKey), // Required
		ACL:        aws.String("public-read"),
	}
	resp, err := svc.CopyObject(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		logrus.Error(err.Error())
	}
	logrus.Debug("Object: ", resp)
	fmt.Println("Successfully released latest build: ", s3URL+targetKey)
	return resp
}
