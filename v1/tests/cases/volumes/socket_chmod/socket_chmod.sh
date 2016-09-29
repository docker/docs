#!/bin/sh -ex

cd /tmp
touch socket_chmod_output
socat unix-listen:sock open:socket_copy_output \
      > socat_listen.stdout \
      2> socat_list.stderr &
listen_pid=$!

# check on our listener
sleep 1
if ! (ps | grep -v grep | grep " $listen_pid "); then
    wait $listen_pid
    echo "Listener died young with exit code $?"
    echo "stdout:"
    cat socat_listen.stdout
    echo "stderr:"
    cat socat_listen.stderr
    exit 1
fi

# now we wait until the background task succeeds in creation of socket
tries=0
while true ; do
    if [[ "$tries" -gt 10 ]]; then
        echo "Polled for socket file more than 10 times. Failure."
        exit 1
    elif [ -S sock ]; then
        echo "Background socat created socket file"
        break
    else
        tries=$((tries+1))
        echo "Did not find socket file on attempt $tries"
        sleep 1
    fi
done

chmod a-rw sock
chmod a+rw sock
