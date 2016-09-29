#!/bin/sh

echo '# time speed-mbit/sec'
echo "$(date +%s) $(grep '0.0-30.0 sec' | grep 'Mbits/sec' | awk '{print $7}')"
