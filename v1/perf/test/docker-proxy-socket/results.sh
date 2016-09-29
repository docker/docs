#!/bin/sh
echo '# wall-clock time as measured by "time"'
echo '# time baseline copy-up copy-down'
echo "$(date +%s) $(cat logs/baseline.time | grep real | cut -f 2 -d 'm' | cut -f 1 -d 's') $(cat logs/test-up.time | grep real | cut -f 2 -d 'm' | cut -f 1 -d 's') $(cat logs/test-down.time | grep real | cut -f 2 -d 'm' | cut -f 1 -d 's')"
