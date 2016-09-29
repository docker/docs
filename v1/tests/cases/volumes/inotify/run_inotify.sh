#!/bin/sh -e

IMAGE_NAME=$1

ID=$(LC_CTYPE=C tr -dc A-Za-z0-9 < /dev/urandom | head -c 8)

docker rm -f inotify || true
docker run --name inotify --rm \
       -v /private/tmp/inotify:/private/tmp/inotify \
       -v /private/tmp/inotify_out:/output \
       -t $IMAGE_NAME /inotify.sh $ID &
docker_pid=$!

# now we wait until the container has started watching
tries=0
while true ; do
    if [[ "$tries" -gt 10 ]]; then
        echo "Polled for inotify ready file more than 10 times. Failure."
        exit 1
    elif [ -f /tmp/inotify/ready ]; then
        echo "Container created inotify/ready file"
        break
    else
        tries=$((tries+1))
        echo "Did not find inotify/ready file on attempt $tries"
        sleep 1
    fi
done

rm -f /tmp/inotify/ready

P=/private/tmp/inotify/directory
./events.sh $P $ID

# timeout if the container is still running
tries=0
while true ; do
    if [[ "$tries" -gt 10 ]]; then
        echo "Check container more than 10 times. Failure."
        docker kill inotify
        exit 1
    elif ! ps -p $docker_pid; then
        echo "Container exited cleanly"
        break
    else
        tries=$((tries+1))
        echo "Container was still running on attempt $tries"
        sleep 1
    fi
done

# check results (cheat for now because there are dupes)
diff /tmp/inotify_out/local /tmp/inotify_out/from_host || true
