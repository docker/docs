package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsCmd struct {
	S3UpdateConfig   AwsS3UpdateConfigCmd   `kong:"cmd,name=s3-update-config"`
	LambdaInvoke     AwsLambdaInvokeCmd     `kong:"cmd,name=lambda-invoke"`
	CloudfrontUpdate AwsCloudfrontUpdateCmd `kong:"cmd,name=cloudfront-update"`
}

type AwsS3UpdateConfigCmd struct {
	Region   string `kong:"name='region',env='AWS_REGION'"`
	S3Bucket string `kong:"name='s3-bucket',env='AWS_S3_BUCKET'"`
	S3Config string `kong:"name='s3-website-config',env='AWS_S3_CONFIG'"`
}

func (s *AwsS3UpdateConfigCmd) Run() error {
	file, err := os.ReadFile(s.S3Config)
	if err != nil {
		return fmt.Errorf("failed to read s3 config file %s: %w", s.S3Config, err)
	}

	data := s3.WebsiteConfiguration{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", s.S3Config, err)
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: awsCredentials(),
		Region:      aws.String(s.Region),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

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
		Credentials: awsCredentials(),
		Region:      aws.String(s.Region),
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

type AwsCloudfrontUpdateCmd struct {
	Region        string `kong:"name='region',env='AWS_REGION'"`
	Function      string `kong:"name='lambda-function',env='AWS_LAMBDA_FUNCTION'"`
	FunctionFile  string `kong:"name='lambda-function-file',env='AWS_LAMBDA_FUNCTION_FILE'"`
	CloudfrontID  string `kong:"name='cloudfront-id',env='AWS_CLOUDFRONT_ID'"`
	RedirectsJSON string `kong:"name='redirects-json',env='REDIRECTS_JSON'"`
}

func (s *AwsCloudfrontUpdateCmd) Run() error {
	var err error
	ver := time.Now().UTC().Format(time.RFC3339)

	zipdt, err := getLambdaFunctionZip(s.FunctionFile, s.RedirectsJSON)
	if err != nil {
		return fmt.Errorf("cannot create lambda function zip: %w", err)
	}

	svc := lambda.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})), &aws.Config{
		Credentials: awsCredentials(),
		Region:      aws.String(s.Region),
	})

	function, err := svc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(s.Function),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() != lambda.ErrCodeResourceNotFoundException {
			return fmt.Errorf("cannot find lambda function %q: %w", s.Function, err)
		}
		_, err = svc.CreateFunction(&lambda.CreateFunctionInput{
			FunctionName: aws.String(s.Function),
		})
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() != lambda.ErrCodeResourceConflictException {
			return err
		}
	}
	codeSha256 := *function.Configuration.CodeSha256
	log.Printf("INFO: updating lambda function %q\n", s.Function)

	updateConfig, err := svc.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(s.Function),
		ZipFile:      zipdt,
	})
	if err != nil {
		return fmt.Errorf("failed to update lambda function code: %s", err)
	}
	log.Printf("INFO: lambda function updated successfully (%s)\n", *updateConfig.FunctionArn)

	if codeSha256 == *updateConfig.CodeSha256 {
		log.Printf("INFO: lambda function code has not changed. skipping publication...")
		return nil
	}

	log.Printf("INFO: waiting for lambda function to be processed\n")
	// the lambda function code image is never ready right away, AWS has to
	// process it, so we wait 3 seconds before trying to publish the version.
	time.Sleep(3 * time.Second)

	publishConfig, err := svc.PublishVersion(&lambda.PublishVersionInput{
		FunctionName: aws.String(s.Function),
		CodeSha256:   aws.String(*updateConfig.CodeSha256),
		Description:  aws.String(ver),
	})
	if err != nil {
		return fmt.Errorf("failed to publish lambda function version %q for %q: %w", ver, s.Function, err)
	}
	log.Printf("INFO: lambda function version %q published successfully (%s)\n", ver, *publishConfig.FunctionArn)

	sess, err := session.NewSession(&aws.Config{
		Credentials: awsCredentials(),
		Region:      aws.String(s.Region)},
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	cfrt := cloudfront.New(sess)
	cfrtDistrib, err := cfrt.GetDistribution(&cloudfront.GetDistributionInput{
		Id: aws.String(s.CloudfrontID),
	})
	if err != nil {
		return fmt.Errorf("cannot find cloudfront distribution %q: %w", s.CloudfrontID, err)
	}
	log.Printf("INFO: cloudfront distribution %q loaded\n", *cfrtDistrib.Distribution.Id)

	cfrtDistribConfig, err := cfrt.GetDistributionConfig(&cloudfront.GetDistributionConfigInput{
		Id: aws.String(s.CloudfrontID),
	})
	if err != nil {
		return fmt.Errorf("cannot load cloudfront distribution config: %w", err)
	}
	log.Printf("INFO: cloudfront distribution configuration loaded\n")

	distribConfig := cfrtDistribConfig.DistributionConfig
	if distribConfig.DefaultCacheBehavior == nil {
		log.Printf("INFO: cloudfront distribution default cache behavior not found. skipping...")
		return nil
	}

	for _, funcAssoc := range distribConfig.DefaultCacheBehavior.LambdaFunctionAssociations.Items {
		if *funcAssoc.EventType != cloudfront.EventTypeViewerRequest {
			continue
		}
		log.Printf("INFO: cloudfront distribution viewer request function ARN found: %q\n", *funcAssoc.LambdaFunctionARN)
	}

	log.Printf("INFO: updating cloudfront config with viewer request function ARN %q", *publishConfig.FunctionArn)
	distribConfig.DefaultCacheBehavior.LambdaFunctionAssociations = &cloudfront.LambdaFunctionAssociations{
		Quantity: aws.Int64(1),
		Items: []*cloudfront.LambdaFunctionAssociation{
			{
				EventType:         aws.String(cloudfront.EventTypeViewerRequest),
				IncludeBody:       aws.Bool(false),
				LambdaFunctionARN: publishConfig.FunctionArn,
			},
		},
	}
	_, err = cfrt.UpdateDistribution(&cloudfront.UpdateDistributionInput{
		Id:                 aws.String(s.CloudfrontID),
		IfMatch:            cfrtDistrib.ETag,
		DistributionConfig: distribConfig,
	})
	if err != nil {
		return err
	}

	log.Printf("INFO: cloudfront config updated successfully\n")
	return nil
}

func getLambdaFunctionZip(funcFilename string, redirectsJSON string) ([]byte, error) {
	funcdt, err := os.ReadFile(funcFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to read lambda function file %q: %w", err)
	}

	var funcbuf bytes.Buffer
	functpl := template.Must(template.New("").Parse(string(funcdt)))
	if err = functpl.Execute(&funcbuf, struct {
		RedirectsJSON string
	}{
		redirectsJSON,
	}); err != nil {
		return nil, err
	}

	tmpdir, err := os.MkdirTemp("", "lambda-zip")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpdir)

	zipfile, err := os.Create(path.Join(tmpdir, "lambda-function.zip"))
	if err != nil {
		return nil, err
	}
	defer zipfile.Close()

	zipwrite := zip.NewWriter(zipfile)
	zipindex, err := zipwrite.Create("index.js")
	if err != nil {
		return nil, err
	}
	if _, err = zipindex.Write(funcbuf.Bytes()); err != nil {
		return nil, err
	}
	if err = zipwrite.Close(); err != nil {
		return nil, err
	}

	zipdt, err := os.ReadFile(zipfile.Name())
	if err != nil {
		return nil, err
	}

	return zipdt, nil
}

func awsCredentials() *credentials.Credentials {
	return credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.StaticProvider{
				Value: credentials.Value{
					AccessKeyID:     getEnvOrSecret("AWS_ACCESS_KEY_ID"),
					SecretAccessKey: getEnvOrSecret("AWS_SECRET_ACCESS_KEY"),
					SessionToken:    getEnvOrSecret("AWS_SESSION_TOKEN"),
				},
			},
		},
	)
}
