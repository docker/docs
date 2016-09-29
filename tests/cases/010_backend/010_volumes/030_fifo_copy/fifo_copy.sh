#!/bin/sh

set -ex

cd /tmp
mkfifo fifo
cat fifo_copy_input > fifo &
cat fifo > fifo_copy_output
