#!/bin/sh

set -ex

ID=$1
EVENTS=create,delete,delete_self,modify,attrib
P=/private/tmp/inotify/directory

mkdir -p "$P"
cd "$P" && inotify-events --events=$EVENTS > /output/from_host &
inotify_event_pid=$!
cd /

touch /private/tmp/inotify/ready

wait $inotify_event_pid

# sleep 1 # wait for the inotify event queue to expire

mkdir -p "$P"
cd "$P" && inotify-events --events="$EVENTS" > /output/local &
inotify_event_pid=$!
cd /

./events.sh "$P" "$ID"

wait $inotify_event_pid
