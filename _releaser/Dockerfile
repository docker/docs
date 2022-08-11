# syntax=docker/dockerfile:1

ARG GO_VERSION=1.18

FROM scratch AS sitedir

FROM golang:${GO_VERSION}-alpine AS base
RUN apk add --no-cache jq openssl
ENV CGO_ENABLED=0
WORKDIR /src

FROM base AS vendor
COPY go.mod go.sum *.go ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendor AS releaser
COPY go.mod go.sum *.go ./
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o /out/releaser .

FROM base AS netlify-remove
ARG NETLIFY_SITE_NAME
RUN --mount=type=bind,from=releaser,source=/out/releaser,target=/usr/bin/releaser \
  --mount=type=secret,id=NETLIFY_AUTH_TOKEN \
  releaser netlify remove

FROM base AS netlify-deploy
ARG NETLIFY_SITE_NAME
RUN --mount=type=bind,from=sitedir,target=/site \
  --mount=type=bind,from=releaser,source=/out/releaser,target=/usr/bin/releaser \
  --mount=type=secret,id=NETLIFY_AUTH_TOKEN \
  --mount=type=secret,id=NETLIFY_ACCOUNT_SLUG \
  NETLIFY_DIR=/site releaser netlify deploy

FROM base AS aws-s3-update-config
ARG AWS_REGION
ARG AWS_S3_BUCKET
ARG AWS_S3_CONFIG
RUN --mount=type=bind,target=. \
  --mount=type=bind,from=releaser,source=/out/releaser,target=/usr/bin/releaser \
  --mount=type=secret,id=AWS_ACCESS_KEY_ID \
  --mount=type=secret,id=AWS_SECRET_ACCESS_KEY \
  releaser aws s3-update-config

FROM base AS aws-lambda-invoke
ARG AWS_REGION
ARG AWS_LAMBDA_FUNCTION
RUN --mount=type=bind,from=releaser,source=/out/releaser,target=/usr/bin/releaser \
  --mount=type=secret,id=AWS_ACCESS_KEY_ID \
  --mount=type=secret,id=AWS_SECRET_ACCESS_KEY \
  releaser aws lambda-invoke
