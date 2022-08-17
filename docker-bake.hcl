variable "JEKYLL_ENV" {
  default = "development"
}
variable "DOCS_URL" {
  default = "http://localhost:4000"
}
variable "DOCS_SITE_DIR" {
  default = "_site"
}

target "_common" {
  args = {
    JEKYLL_ENV = JEKYLL_ENV
    DOCS_URL = DOCS_URL
  }
  no-cache-filter = ["generate"]
}

group "default" {
  targets = ["release"]
}

target "release" {
  inherits = ["_common"]
  target = "release"
  output = [DOCS_SITE_DIR]
}

target "vendor" {
  target = "vendor"
  output = ["."]
}

group "validate" {
  targets = ["htmlproofer", "mdl"]
}

target "htmlproofer" {
  inherits = ["_common"]
  target = "htmlproofer"
  output = ["type=cacheonly"]
}

target "htmlproofer-output" {
  inherits = ["_common"]
  target = "htmlproofer-output"
  output = ["./lint"]
}

target "mdl" {
  inherits = ["_common"]
  target = "mdl"
  output = ["type=cacheonly"]
}

target "mdl-output" {
  inherits = ["_common"]
  target = "mdl-output"
  output = ["./lint"]
  args = {
    MDL_JSON = 1
  }
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
