#!/bin/sh

DOCKER_CMD=docker

alias AWS_HUB_PROD='aws ec2 describe-instances --filters "Name=tag:aws:cloudformation:stack-name,Values=us-east-1*" "Name=tag:secondary-role,Values=hub" "Name=instance-state-name,Values=running" --output=json'
alias AWS_HUB_STAGE='aws ec2 describe-instances --filters "Name=tag:aws:cloudformation:stack-name,Values=stage-us-east-1*" "Name=tag:secondary-role,Values=hub" "Name=instance-state-name,Values=running" --output=json'
alias AWS_IP="jq -r '.Reservations[].Instances[].PrivateIpAddress'"

HUB_GATEWAY="https://hub.docker.com"
HUB_SERVICE_NAME="hub-web-v2"

DEFAULT_IMAGE_PROD="bagel/hub-prod"
DEFAULT_IMAGE_STAGE="bagel/hub-stage"

NEW_RELIC_APP_NAME="hub.docker.com(aws-node)"
NEW_RELIC_LICENSE_KEY="582e3891446a63a3f99b4d32f9585ec74af1d8d7"

NO_COLOR="\033[0m"
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"

MESSAGE_MISSING_OR_INVALID_ARGS="${RED}Missing or invalid arguments${NO_COLOR}"

# $1: prod or stage
getAWSHosts() {
  if [ $1 == "prod" ]; then
    echo $(AWS_HUB_PROD | ( AWS_IP ; echo ) | sed -e ':a' -e 'N' -e '$!ba' -e 's/\n/:2376 /g')
  elif [ $1 == "stage" ]; then
    echo $(AWS_HUB_STAGE | ( AWS_IP ; echo ) | sed -e ':a' -e 'N' -e '$!ba' -e 's/\n/:2376 /g')
  fi
}

# $1: Exit code
printUsageAndExit() {
  echo
  echo "Usage: deploy.sh [prod|stage|-h <host>] [IMAGE]"
  echo
  echo "  prod        A predefined list of hosts for production"
  echo "  stage       A predefined list of hosts for staging"
  echo "  -h <host>   A single host address"
  echo
  exit $1
}

# $1: Image argument
parseImageArg() {
  if [ -z "$1" ]; then
    echo $MESSAGE_MISSING_OR_INVALID_ARGS
    printUsageAndExit 1
  fi
  IMAGE=$1
}

parseArgs() {
  if [ $1 == "-h" ]; then
    parseImageArg $3
    HOSTS=$2
  else
    if [ $1 == "prod" ]; then
      if [ -z "$2" ]; then
        IMAGE=$DEFAULT_IMAGE_PROD
      else
        parseImageArg $2
      fi
      HOSTS=$( getAWSHosts "prod" )
    elif [ $1 == "stage" ]; then
      if [ -z "$2" ]; then
        IMAGE=$DEFAULT_IMAGE_STAGE
      else
        parseImageArg $2
      fi
      HOSTS=$( getAWSHosts "stage" )
    else
      echo
    	echo $MESSAGE_MISSING_OR_INVALID_ARGS
      printUsageAndExit 1
    fi
  fi
}

# $1: Host IP
# $2: Image
# $3: Container name
# $4: Container port
runContainer() {
  $DOCKER_CMD --tlsverify=false -H tcp://$1 run \
    -de ENV=production \
    -e HUB_API_BASE_URL=$HUB_GATEWAY \
    -e REGISTRY_API_BASE_URL=$HUB_GATEWAY \
    -e SERVICE_NAME=$HUB_SERVICE_NAME \
    -e SERVICE_80_NAME=$HUB_SERVICE_NAME \
    -e NEW_RELIC_LICENSE_KEY=$NEW_RELIC_LICENSE_KEY \
    -e NEW_RELIC_APP_NAME=$NEW_RELIC_APP_NAME \
    -e PORT=80 \
    -p $4:80 \
    --restart "unless-stopped" \
    --name $3 \
    $2
}

# $1: Host IP
# $2: Container name
removeContainer() {
  $DOCKER_CMD --tlsverify=false -H tcp://$1 stop $2
  $DOCKER_CMD --tlsverify=false -H tcp://$1 rm $2
}

# $1: Host IP
# $2: Image name
pullImage() {
  $DOCKER_CMD --tlsverify=false -H tcp://$1 pull $2
}

# $1: Host IP
# $2: Image
deployHost() {
  echo
  echo "Starting to deploy ${YELLOW}$IMAGE${NO_COLOR} to ${YELLOW}$1${NO_COLOR}"

  pullImage $1 $2

  removeContainer $1 "hub_2_0"
  runContainer $1 $2 "hub_2_0" 6600

  removeContainer $1 "hub_2_1"
  runContainer $1 $2 "hub_2_1" 6601

  removeContainer $1 "hub_2_2"
  runContainer $1 $2 "hub_2_2" 6602
}

# Prerequisites:
# 1- AWS
type aws >/dev/null 2>&1 || { echo >&2 "AWS client is required. Make sure 'aws' command is available:\nhttp://docs.aws.amazon.com/cli/latest/userguide/installing.html"; exit 1; }
# 2- JQ
type jq >/dev/null 2>&1 || { echo >&2 "jq JSON processor is required. Make sure 'jq' command is available:\nbrew install jq"; exit 1; }

# Case for no paremeters specified
if [ -z "$1" ]
  then
    echo
    echo $MESSAGE_MISSING_OR_INVALID_ARGS
    printUsageAndExit 1
fi

parseArgs "$@"

echo
echo "Image: ${YELLOW}$IMAGE ${NO_COLOR}"
echo "Hosts: ${YELLOW}$HOSTS${NO_COLOR}"
echo
read -p "Do you want to proceed? [Y/n]" -s -n 1 KEY
echo
if [[ ! $KEY =~ ^[Yy]$ ]]; then
  exit 1
fi

# Run deployment for each host
for HUB_HOST in $HOSTS
do
  deployHost $HUB_HOST $IMAGE
  echo
  echo "Sleeping for 10 seconds to let the containers boot up..."
  echo
  sleep 10
done

echo
echo "${GREEN}All done!${NO_COLOR}"
echo
