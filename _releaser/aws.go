package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsCmd struct {
	S3UpdateConfig AwsS3UpdateConfigCmd `kong:"cmd,name=s3-update-config"`
	LambdaInvoke   AwsLambdaInvokeCmd   `kong:"cmd,name=lambda-invoke"`
}

type AwsS3UpdateConfigCmd struct {
	Region   string `kong:"name='region',env='AWS_REGION'"`
	S3Bucket string `kong:"name='s3-bucket',env='AWS_S3_BUCKET'"`
	S3Config string `kong:"name='s3-website-config',env='AWS_S3_CONFIG'"`
}

func (s *AwsS3UpdateConfigCmd) Run() error {
	file, err := ioutil.ReadFile(s.S3Config)
	if err != nil {
		return fmt.Errorf("failed to read s3 config file %s: %w", s.S3Config, err)
	}

	data := s3.WebsiteConfiguration{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", s.S3Config, err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.Region)},
	)

	svc := s3.New(sess)

	// Create SetBucketWebsite parameters based on the JSON file input
	params := s3.PutBucketWebsiteInput{
		Bucket:               aws.String(s.S3Bucket),
		WebsiteConfiguration: &data,
	}

	// Set the website configuration on the bucket.
	// Replacing any existing configuration.
	_, err = svc.PutBucketWebsite(&params)
	if err != nil {
		return fmt.Errorf("unable to set bucket %q website configuration: %w", s.S3Bucket, err)
	}

	log.Printf("INFO: successfully set bucket %q website configuration\n", s.S3Bucket)
	return nil
}

type AwsLambdaInvokeCmd struct {
	Region         string `kong:"name='region',env='AWS_REGION'"`
	LambdaFunction string `kong:"name='lambda-function',env='AWS_LAMBDA_FUNCTION'"`
}

func (s *AwsLambdaInvokeCmd) Run() error {
	svc := lambda.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})), &aws.Config{
		Region: aws.String(s.Region),
	})

	_, err := svc.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(s.LambdaFunction),
	})
	if err != nil {
		return err
	}

	log.Printf("INFO: lambda function %q invoked successfully\n", s.LambdaFunction)
	return nil
}
