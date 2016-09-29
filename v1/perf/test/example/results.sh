#!/bin/sh

# Post-process the results and produce a graphable datapoint
# in gnuplot format
echo '# "time since epoch in seconds" "number of figlet invocations"'
echo "$(date +%s) 1"
