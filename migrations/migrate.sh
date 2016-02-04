#!/bin/bash

# When run in the docker containers, the working directory
# is the root of the repo.

iter=0

case $SERVICE_NAME in
	notaryserver)
		# have to poll for DB to come up
		until migrate -path=migrations/server/mysql -url="mysql://server@tcp(notarymysql:3306)/notaryserver" up
		do
			((iter++))
			if (( iter > 30 )); then
				exit 1;
			fi
			echo "waiting for notarymysql to come up."
			sleep 1
		done
		echo "notaryserver database migrated to latest version"
		;;
	notarysigner)
		# have to poll for DB to come up
		until migrate -path=migrations/signer/mysql -url="mysql://signer@tcp(notarymysql:3306)/notarysigner" up 
		do
			((iter++))
			if (( iter > 30 )); then
				exit 1;
			fi
			echo "waiting for notarymysql to come up."
			sleep 1
		done
		echo "notarysigner database migrated to latest version"
		;;
esac
