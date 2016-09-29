#!/bin/sh

. ../../lib/functions

echo '# block-size-bytes time/sec'
for file in $(ls logs/*.time | sort -n -t '/' -k 2); do
    block_size=$(echo $file | sed 's!logs/\(.*\)\.time!\1!')
    echo $block_size '\t' $(cat $file)
 done

