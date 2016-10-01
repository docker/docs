#!/bin/bash
middleman build
aws s3 sync /app/build s3://${AWS_S3_BUCKET}/apidocs/docker-cloud/ --delete --acl public-read --region us-east-1
