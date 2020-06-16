package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/sethvargo/go-githubactions"
)

// Upload a bucket's website configuration.
func main() {
	bucketName := githubactions.GetInput("bucketName")
	if bucketName == "" {
		githubactions.Fatalf("missing input 'bucketName'")
	}

	regionName := githubactions.GetInput("regionName")
	if regionName == "" {
		githubactions.Fatalf("missing input 'regionName'")
	}

	websiteConfig := githubactions.GetInput("websiteConfig")
	if websiteConfig == "" {
		githubactions.Fatalf("missing input 'websiteConfig'")
	}

	file, err := ioutil.ReadFile(websiteConfig)
	if err != nil {
		exitErrorf("Error reading json file %s, %v", websiteConfig, err)
	}

	data := s3.WebsiteConfiguration{}
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		exitErrorf("Error parsing JSON from file %s, %v", websiteConfig, err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(regionName)},
	)

	svc := s3.New(sess)

	// Create SetBucketWebsite parameters based on the JSON file input
	params := s3.PutBucketWebsiteInput{
		Bucket:               aws.String(bucketName),
		WebsiteConfiguration: &data,
	}

	// Set the website configuration on the bucket. Replacing any existing
	// configuration.
	_, err = svc.PutBucketWebsite(&params)
	if err != nil {
		exitErrorf("Unable to set bucket %q website configuration, %v",
			bucketName, err)
	}

	fmt.Printf("Successfully set bucket %q website configuration\n", bucketName)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
