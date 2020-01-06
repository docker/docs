---
description: Cluster CLI environment variables
keywords: documentation, docs, docker, cluster, infrastructure, automation
title: Cluster CLI environment variables
---

>{% include enterprise_label_shortform.md %}

Use the following environment variables as needed to configure the Docker Cluster command-line behavior.

## AWS\_ACCESS\_KEY\_ID
Represents your AWS Access Key. Overrides the use of `AWS_SHARED_CREDENTIALS_FILE` and `AWS_PROFILE`.

```bash
export AWS_ACCESS_KEY_ID="AKIFAKEAWSACCESSKEYNLQ"
```

## AWS\_SECRET\_ACCESS\_KEY
Represents your AWS Secret Key. Overrides the use of `AWS_SHARED_CREDENTIALS_FILE` and `AWS_PROFILE`.
```bash
export AWS_SECRET_ACCESS_KEY="3SZYfAkeS3cr3TKey+L0ok5/rEalBu71sFak3vmy"
```

## AWS\_DEFAULT\_REGION
Specifies the AWS region to provision resources.
```bash
export AWS_DEFAULT_REGION="us-east-1"
```

## AWS\_PROFILE
Specifies the AWS profile name as set in the shared credentials file.
```bash
export AWS_PROFILE="default"
```
## AWS\_SESSION\_TOKEN
Specifies the session token used for validating temporary credentials. This is typically provided after 
successful identity federation or Multi-Factor Authentication (MFA) login. With MFA login, this is the 
session token provided afterwards, not the 6 digit MFA code used to get temporary credentials.
```bash
export AWS_SESSION_TOKEN=AQoDYXdzEJr...<remainder of security token>
```
## AWS\_SHARED\_CREDENTIALS\_FILE
Specifies the path to the shared credentials file. If this is not set and a profile is specified, `~/.aws/credentials` 
is used.

```bash
export AWS_SHARED_CREDENTIALS_FILE="~/.production/credentials"
```

## CLUSTER\_ORGANIZATION
Specifies the Docker Hub organization to pull the `cluster` container.

```bash
export CLUSTER_ORGANIZATION="docker"
```

## CLUSTER\_TAG
Specifies the tag of the `cluster` container to pull.

```bash
export CLUSTER_TAG="latest"
```

## DOCKER\_PASSWORD
Overrides docker password lookup from `~/.docker/config.json`.

```bash
export DOCKER_PASSWORD="il0v3U3000!"
```
## DOCKER\_USERNAME
Overrides docker username lookup from `~/.docker/config.json`.

```bash
export DOCKER_USERNAME="ironman"
```
