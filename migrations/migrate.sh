#!/bin/bash

# When run in the docker containers, the working directory
# is the root of the repo.

iter=0

case $SERVICE_NAME in
	notary_server)
		# have to poll for DB to come up
		until migrate -path=migrations/server/mysql -url="mysql://server@tcp(mysql:3306)/notaryserver" version > /dev/null
		do
			((iter++))
			if (( iter > 30 )); then
				echo "notaryserver database failed to come up within 30 seconds"
				exit 1;
			fi
			echo "waiting for notarymysql to come up."
			sleep 1
		done
		pre=$(migrate -path=migrations/server/mysql -url="mysql://server@tcp(mysql:3306)/notaryserver" version)
		if migrate -path=migrations/server/mysql -url="mysql://server@tcp(mysql:3306)/notaryserver" up ; then
			post=$(migrate -path=migrations/server/mysql -url="mysql://server@tcp(mysql:3306)/notaryserver" version)
			if [ "$pre" != "$post" ]; then
				echo "notaryserver database migrated to latest version"
			else
				echo "notaryserver database already at latest version"
			fi
		else
			echo "notaryserver database migration failed"
			exit 1
		fi
		;;
	notary_signer)
		# have to poll for DB to come up
		until migrate -path=migrations/signer/mysql -url="mysql://signer@tcp(mysql:3306)/notarysigner" version > /dev/null
		do
			((iter++))
			if (( iter > 30 )); then
				echo "notarysigner database failed to come up within 30 seconds"
				exit 1;
			fi
			echo "waiting for notarymysql to come up."
			sleep 1
		done
		pre=$(migrate -path=migrations/signer/mysql -url="mysql://signer@tcp(mysql:3306)/notarysigner" version)
		if migrate -path=migrations/signer/mysql -url="mysql://signer@tcp(mysql:3306)/notarysigner" up ; then 
			post=$(migrate -path=migrations/signer/mysql -url="mysql://signer@tcp(mysql:3306)/notarysigner" version)
			if [ "$pre" != "$post" ]; then
				echo "notarysigner database migrated to latest version"
			else
				echo "notarysigner database already at latest version"
			fi
		else
			echo "notarysigner database migration failed"
			exit 1
		fi
		;;
esac
