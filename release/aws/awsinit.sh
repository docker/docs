#!/bin/sh
PASSWORD=$(curl http://169.254.169.254/latest/meta-data/instance-id/)
while true; do
    curl https://localhost -k
    if [ $? -eq 0 ]; then
        # The loadbalancer is up... wait a bit more and then try
        sleep 30
        curl --data '{"method":"managed","managed":{"users":[{"username":"admin","password":"'$PASSWORD'","isReadOnly":false,"isNew":true,"isAdmin":true,"isReadWrite":false,"teamsChanged":true}]}}' https://localhost/api/v0/admin/settings/auth -k -X PUT -H 'Content-Type: application/json; charset=utf-8'
        break
    else
        sleep 1
    fi
done
