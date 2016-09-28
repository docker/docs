# Docker AWS Editions Notes

These instructions show you how to run UCP Seattle on top of the Docker AWS Edition

Make sure you've got the latest version of the AWS CLI:

```bash
sudo pip install awscli --upgrade
```

Make sure you've got all your environment variables set and/or have configured the aws CLI so it can authenticate properly.

## Deploy your swarm

```bash
STACK_NAME=dhstack3
# The name of your pre-loaded SSH public key on AWS:
KEY_NAME=dhiltgen
MANAGER_AWS_INSTANCE_TYPE=m3.medium
WORKER_AWS_INSTANCE_TYPE=t2.micro
MANAGER_COUNT=1
WORKER_COUNT=1
LOCAL_PROXY_PORT=2375


aws cloudformation create-stack \
    --stack-name ${STACK_NAME} \
    --template-url https://s3-us-west-2.amazonaws.com/docker-cf-templates/docker_on_aws_brc_1.json \
    --parameters \
        ParameterKey=KeyName,ParameterValue=${KEY_NAME} \
        ParameterKey=ManagerInstanceType,ParameterValue=${MANAGER_AWS_INSTANCE_TYPE} \
        ParameterKey=InstanceType,ParameterValue=${WORKER_AWS_INSTANCE_TYPE} \
        ParameterKey=ManagerSize,ParameterValue=${MANAGER_COUNT} \
        ParameterKey=ClusterSize,ParameterValue=${WORKER_COUNT} \
    --capabilities CAPABILITY_IAM

time aws cloudformation wait stack-create-complete --stack-name ${STACK_NAME}

SSH_ELB_PHYS_ID=$(aws cloudformation describe-stack-resources --stack-name ${STACK_NAME} --logical-resource-id SSHLoadBalancer | jq -r ".StackResources[0].PhysicalResourceId")

SSH_ELB_HOSTNAME=$(aws elb describe-load-balancers --load-balancer-names ${SSH_ELB_PHYS_ID} | jq -r ".LoadBalancerDescriptions[0].DNSName")

# Add port 443 since we'll need it later...
aws elb create-load-balancer-listeners --load-balancer-name ${SSH_ELB_PHYS_ID} --listeners "Protocol=TCP,LoadBalancerPort=443,InstanceProtocol=TCP,InstancePort=443"

# Make sure things look good
ssh root@${SSH_ELB_HOSTNAME} /usr/docker/bin/docker swarm inspect

ssh -NL localhost:${LOCAL_PROXY_PORT}:/var/run/docker.sock root@${SSH_ELB_HOSTNAME} &
TUNNEL_PID=$!

```

## Deploy UCP on top of it

You can either deploy from your local build, or deploy from official
builds pushed up to the private hub org (all Docker employees should
have access)

### Developer mode (assumes you're at the top of the orca tree)

```bash
# Copy binaries
DOCKER=/usr/docker/bin/docker ./script/copy_orca_images root@${SSH_ELB_HOSTNAME}

# install it
docker -H localhost:${LOCAL_PROXY_PORT} run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp:1.2.0-dev \
        install --san ${SSH_ELB_HOSTNAME}
echo "ACTUAL LOGIN: https://${SSH_ELB_HOSTNAME}"

# Kill the tunnel now since we no longer need it
kill ${TUNNEL_PID}
```

### Using latest official private build from hub

```bash
# Load binaries (make sure to `docker login` first)
for i in $(docker run --rm dockerorcadev/ucp:1.2.0-latest images --list --image-version dev: ) ; do
    docker -H localhost:${LOCAL_PROXY_PORT} pull $i
done

# install it
docker -H localhost:${LOCAL_PROXY_PORT} run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        dockerorcadev/ucp:1.2.0-latest \
        install --image-version dev: --san ${SSH_ELB_HOSTNAME}
echo "ACTUAL LOGIN: https://${SSH_ELB_HOSTNAME}"

# Kill the tunnel now since we no longer need it
kill ${TUNNEL_PID}
```


## Clean up when you're done

**PLEASE don't leave zombies around!!!**

```bash
aws cloudformation delete-stack --stack-name ${STACK_NAME}
```
