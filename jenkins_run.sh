#!/usr/bin/env bash
rm -f config.properties

export HOME="$WORKSPACE"
ls -laR ~/.docker

docker-machine env $DOCKER_MACHINE_NAME || exit 1

eval "$( docker-machine env $DOCKER_MACHINE_NAME )"
# eval "$( cat $HOME/.docker/.dockerenv )"

function diagnostics {
    docker ps

    mkdir -p $HOME/integration/results/
    for machine in $(docker ps -aq)
    do
        docker logs $machine 2> $HOME/integration/results/$machine.err  > $HOME/integration/results/$machine.log
    done
	docker-machine --debug ssh $DOCKER_MACHINE_NAME -- "sudo chmod 777 /var/log/upstart/docker.log"
    docker-machine --debug scp -r $DOCKER_MACHINE_NAME:/var/log/upstart/docker.log $HOME/integration/results/
    
    echo "DHE_HOST=$DTR_HOST" > config.properties
    echo "DTR_HOST=$DTR_HOST" >> config.properties
}

# debug info
docker version
docker info
echo UCP_AS_ENZI = $UCP_AS_ENZI


# we set +x so we don't print the credentials to the logs
set +x
source "$DOCKER_CREDENTIALS_FILE"
docker login -u "$DOCKER_HUB_USERNAME" -p "$DOCKER_HUB_PASSWORD" -e "a@a.a"
set -x


# find the ip of the machine
DTR_HOST="$( docker-machine ip $DOCKER_MACHINE_NAME )"
if [[ -z "$DTR_HOST" ]]
then
	diagnostics
    exit 1
fi


# install UCP

case "$UCP_IMAGE" in
    'dockerorcadev/'*)
    	for i in $(docker run --rm $UCP_IMAGE images --list --image-version dev: ) ; do docker pull $i; done
        docker run -e UCP_ADMIN_USER=admin -e UCP_ADMIN_PASSWORD=password --name ucp \
            -v /var/run/docker.sock:/var/run/docker.sock \
            $UCP_IMAGE install --image-version dev: --swarm-port 3376 --controller-port 444 --san $DTR_HOST --host-address $DTR_HOST
            #  --host-address 172.17.0.1
        ;;
    *)
    	docker run -e UCP_ADMIN_USER=admin -e UCP_ADMIN_PASSWORD=password --name ucp \
            -v /var/run/docker.sock:/var/run/docker.sock \
            $UCP_IMAGE install --swarm-port 3376 --controller-port 444 --san $DTR_HOST --host-address $DTR_HOST
            # --host-address 172.17.0.1
        ;;
esac


docker rm -f ucp

# wire up overlay networking, necessary only before ucp 1.1
#if [ "$UCP_AS_ENZI" == "false" ]
#then
#	docker run -e UCP_ADMIN_USER=admin -e UCP_ADMIN_PASSWORD=password --name ucp \
#    	-v /var/run/docker.sock:/var/run/docker.sock \
#    	docker/ucp:$UCP_TAG engine-discovery --controller 172.17.0.1 --host-address 172.17.0.1 --update
#fi

# restart docker to enable overlay networking

#docker-machine ssh $DOCKER_MACHINE_NAME /etc/init.d/docker restart
docker-machine ssh $DOCKER_MACHINE_NAME sudo service docker restart

# wait for the restart to finish
until docker ps -a;
do
	echo Waiting for daemon to come back with overlay networking enabled...
    sleep 1
done


if [ "$UCP_AS_ENZI" == "false" ]
then  
  # install enzi if we are not using ucp as enzi
  
  # set up certs if they don't exist
  docker volume ls | grep enzi-tls || (docker volume create --name enzi-tls && docker run -i --rm -v enzi-tls:/tls alpine:3.3 sh << 'EOF')
  set -ex
  apk update
  apk upgrade libcrypto1.0 libssl1.0
  apk add openssl
  openssl genrsa -out /tls/key.pem 2048
  openssl req -new -x509 -key /tls/key.pem -out /tls/cert.pem -days 3650 -subj '/CN=*.enzi'
  cp /tls/cert.pem /tls/ca.pem
EOF
  
  # create the network only if necessary
  docker network create enzi || true
  
  # start rethink if needed
  NUM=01
  docker ps | grep enzi-db-$NUM || (docker run -d \
      -v enzi-db-$NUM-data:/var/data `# create a volume for this node's data, mounted at /var/data` \
      -v enzi-tls:/tls               `# mount the tls key pair at /tls` \
      --net enzi                     `# add this container to the enzi network` \
      --name enzi-db-$NUM            `# the container will have the domain name "enzi-db-$NUM.enzi" on the network` \
      --net-alias rethinkdb          `# the domain name 'rethinkdb.enzi' will refer to the longest-running container with this alias` \
      jlhawn/rethinkdb-tls           `# this is the container image to run. It has the entrypoint '/bin/rethinkdb'` \
          --bind all                 `# bind to the network interface` \
          --no-http-admin            `# the admin web console is insecure and useless` \
          --server-name enzi_db_$NUM \
          --canonical-address enzi-db-$NUM.enzi \
          --directory /var/data/rethinkdb \
          --join rethinkdb.enzi      `# The node should ignore joining itself` \
          --driver-tls \
              --driver-tls-key /tls/key.pem \
              --driver-tls-cert /tls/cert.pem \
              --driver-tls-ca /tls/ca.pem \
          --cluster-tls \
              --cluster-tls-key /tls/key.pem \
              --cluster-tls-cert /tls/cert.pem \
              --cluster-tls-ca /tls/ca.pem && sleep 30)
  
  # debug
  docker ps
  
  # set up the db
  docker run --rm \
      --net enzi \
      -v enzi-tls:/tls \
      dockerhubenterprise/dtr-integration-enzi \
          --db-addr=rethinkdb.enzi \
          sync-db || true
  
  # initialize the admin user (if they exist, ignore the failure)
  docker run --rm \
      --net enzi \
      -v enzi-tls:/tls \
      -e USERNAME=admin \
      -e PASSWORD=password \
      dockerhubenterprise/dtr-integration-enzi \
          --db-addr=rethinkdb.enzi \
          create-admin || true
  
  
  # start the api server if needed
  NUM=01
  docker ps | grep enzi-api-$NUM || docker run -d \
      --name enzi-api-$NUM \
      --net enzi \
      --net-alias api \
      --restart always \
      -v enzi-tls:/tls \
      -p 4443:4443 \
      dockerhubenterprise/dtr-integration-enzi \
          --db-addr=rethinkdb.enzi \
          api || true
  
      
  # start the worker if needed
  NUM=01
  docker ps | grep enzi-worker-$NUM || docker run -d \
      --name enzi-worker-$NUM \
      --net enzi \
      --net-alias worker \
      -v enzi-worker-$NUM-data:/work \
      -v enzi-tls:/tls \
      dockerhubenterprise/dtr-integration-enzi \
          --debug \
          --db-addr=rethinkdb.enzi \
          worker \
              --addr=enzi-worker-$NUM.enzi || true

fi


RETRIES=0
# Wait for ucp to stabalize before trying to install dtr
CONTAINERS=0
until [ "$CONTAINERS" != "0" ]
do
	RESP=$(curl -fksSLo - -X POST -d '{"username":"admin","password":"password"}' "https://$DTR_HOST:444/auth/login" || true)
	TOKEN=$(echo "$RESP" | jq ".auth_token" -r || true)
	CONTAINERS=$(curl -k "https://$DTR_HOST:444/api/nodes" -H "Authorization: Bearer "$TOKEN | jq ".[0].containers" -r || true)
    if [ "$CONTAINERS" == "null" ]
    then
    	CONTAINERS=0
    fi
    if [ "$CONTAINERS" == "" ]
    then
    	CONTAINERS=0
    fi
    if [ "$RETRIES" == "600" ]
    then
        echo "Gave up waiting for UCP to stabalize"
        diagnostics
        exit 1
    fi
    RETRIES=`expr $RETRIES + 1`
    sleep 1
    docker ps
done

IMAGE="$DOCKER_IMAGE"
docker pull "$IMAGE" ||:


ENZI_HOST=""
if [ "$UCP_AS_ENZI" == "false" ]
then
	ENZI_HOST="--enzi-host $DTR_HOST:4443"
fi


UNSAFE="--unsafe"
if [ "$IMAGE" == "docker/dtr:2.0.0" ]
then
	UNSAFE=""
fi

LB_PARAM="--dtr-load-balancer"
if [ "$NEW_LB_PARAM" == "true" ]
then
	LB_PARAM="--dtr-external-url"
fi

docker run --rm \
--name dtr \
$IMAGE \
install \
--ucp-url $DTR_HOST:444 \
$LB_PARAM $DTR_HOST \
--ucp-username admin \
--ucp-password password \
--ucp-insecure-tls \
$UNSAFE \
$ENZI_HOST \
--hub-username "$DOCKER_HUB_USERNAME" \
--hub-password "$DOCKER_HUB_PASSWORD" \
|| (diagnostics && exit 1)


# wait for dtr to come up
RETRIES=0
until curl -fksSLo - -u admin:password "https://$DTR_HOST" ; do
	echo "waiting for DTR to come up"
    docker ps
    if [ "$RETRIES" == "60" ]
    then
        echo "Gave up waiting for DTR to come up"
        diagnostics
        exit 1
    fi
    RETRIES=`expr $RETRIES + 1`
    sleep 5
done

diagnostics
