# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23

FROM scratch AS sitedir

FROM golang:${GO_VERSION}-alpine AS base
RUN apk add --no-cache openssl
ENV CGO_ENABLED=0
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

FROM base AS releaser
RUN --mount=type=bind,target=. \
  --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o /out/releaser .

FROM base AS aws-s3-update-config
ARG DRY_RUN=false
ARG AWS_REGION
ARG AWS_S3_BUCKET
ARG AWS_S3_CONFIG
RUN --mount=type=bind,target=. \
  --mount=type=bind,from=releaser,source=/out/releaser,target=/usr/bin/releaser \
  --mount=type=secret,id=AWS_ACCESS_KEY_ID \
  --mount=type=secret,id=AWS_SECRET_ACCESS_KEY \
  --mount=type=secret,id=AWS_SESSION_TOKEN \
  releaser aws s3-update-config

FROM base AS aws-lambda-invoke
ARG DRY_RUN=false
ARG AWS_REGION
ARG AWS_LAMBDA_FUNCTION
RUN --mount=type=bind,from=releaser,source=/out/releaser,target=/usr/bin/releaser \
  --mount=type=secret,id=AWS_ACCESS_KEY_ID \
  --mount=type=secret,id=AWS_SECRET_ACCESS_KEY \
  --mount=type=secret,id=AWS_SESSION_TOKEN \
  releaser aws lambda-invoke

FROM base AS aws-cloudfront-update
ARG DRY_RUN=false
ARG AWS_REGION
ARG AWS_LAMBDA_FUNCTION
ARG AWS_CLOUDFRONT_ID
ARG AWS_LAMBDA_FUNCTION_FILE="cloudfront-lambda-redirects.js"
ARG REDIRECTS_FILE="/site/redirects.json"
ARG REDIRECTS_PREFIXES_FILE="redirects-prefixes.json"
RUN --mount=type=bind,target=. \
  --mount=type=bind,from=sitedir,target=/site \
  --mount=type=bind,from=releaser,source=/out/releaser,target=/usr/bin/releaser \
  --mount=type=secret,id=AWS_ACCESS_KEY_ID \
  --mount=type=secret,id=AWS_SECRET_ACCESS_KEY \
  --mount=type=secret,id=AWS_SESSION_TOKEN \
  releaser aws cloudfront-update
