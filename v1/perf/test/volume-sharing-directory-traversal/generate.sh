#!/bin/sh -e

#./cat.py -d 141 -w 1 > testcases/d_141_w_1

tails=900
for spine in $(seq 1 16); do
    ./cat.py -d $spine -w $tails \
             > testcases/$(printf "d=%03d,w=%d" $spine $tails)
done
