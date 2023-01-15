variable "DOCS_SITE_DIR" {
  default = "public"
}
variable "DOCS_ENFORCE_GIT_LOG_HISTORY" {
  default = "0"
}

group "default" {
  targets = ["release"]
}

target "release" {
  target = "release"
  output = [DOCS_SITE_DIR]
}

#
# releaser targets are defined in _releaser/Dockerfile
# and are used for Netlify and AWS S3 deployment
#

target "releaser-build" {
  context = "_releaser"
  target = "releaser"
  output = ["type=cacheonly"]
}

variable "NETLIFY_SITE_NAME" {
  default = ""
}

target "_common-netlify" {
  args = {
    NETLIFY_SITE_NAME = NETLIFY_SITE_NAME
  }
  secret = [
    "id=NETLIFY_AUTH_TOKEN,env=NETLIFY_AUTH_TOKEN"
  ]
}

target "netlify-remove" {
  inherits = ["_common-netlify"]
  context = "_releaser"
  target = "netlify-remove"
  no-cache-filter = ["netlify-remove"]
  output = ["type=cacheonly"]
}

target "netlify-deploy" {
  inherits = ["_common-netlify"]
  context = "_releaser"
  target = "netlify-deploy"
  contexts = {
    sitedir = DOCS_SITE_DIR
  }
  no-cache-filter = ["netlify-deploy"]
  output = ["type=cacheonly"]
}

variable "AWS_REGION" {
  default = ""
}
variable "AWS_S3_BUCKET" {
  default = ""
}
variable "AWS_S3_CONFIG" {
  default = ""
}
variable "AWS_CLOUDFRONT_ID" {
  default = ""
}
variable "AWS_LAMBDA_FUNCTION" {
  default = ""
}

target "_common-aws" {
  args = {
    AWS_REGION = AWS_REGION
    AWS_S3_BUCKET = AWS_S3_BUCKET
    AWS_S3_CONFIG = AWS_S3_CONFIG
    AWS_CLOUDFRONT_ID = AWS_CLOUDFRONT_ID
    AWS_LAMBDA_FUNCTION = AWS_LAMBDA_FUNCTION
  }
  secret = [
    "id=AWS_ACCESS_KEY_ID,env=AWS_ACCESS_KEY_ID",
    "id=AWS_SECRET_ACCESS_KEY,env=AWS_SECRET_ACCESS_KEY",
    "id=AWS_SESSION_TOKEN,env=AWS_SESSION_TOKEN"
  ]
}

target "aws-s3-update-config" {
  inherits = ["_common-aws"]
  context = "_releaser"
  target = "aws-s3-update-config"
  no-cache-filter = ["aws-update-config"]
  output = ["type=cacheonly"]
}

target "aws-lambda-invoke" {
  inherits = ["_common-aws"]
  context = "_releaser"
  target = "aws-lambda-invoke"
  no-cache-filter = ["aws-lambda-invoke"]
  output = ["type=cacheonly"]
}

target "aws-cloudfront-update" {
  inherits = ["_common-aws"]
  context = "_releaser"
  target = "aws-cloudfront-update"
  contexts = {
    sitedir = DOCS_SITE_DIR
  }
  no-cache-filter = ["aws-cloudfront-update"]
  output = ["type=cacheonly"]
}
