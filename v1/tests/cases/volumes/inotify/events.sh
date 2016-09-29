#!/bin/sh -ex

P=$1
ID=$2
T=0.2

# do our events
touch $P/file_$ID
sleep $T

mkdir $P/dir_$ID
sleep $T

echo "modified" > $P/file_$ID
sleep $T

# time modify file
touch $P/file_$ID
sleep $T

# time modify directory
touch $P/dir_$ID
sleep $T

rmdir $P/dir_$ID
sleep $T

mkdir $P/dir_$ID
sleep $T

echo "modified" > $P/file_$ID

echo "new child" > $P/dir_$ID/child
sleep $T

# perm modify file
chmod 660 $P/file_$ID
sleep $T

# perm modify directory
chmod 770 $P/dir_$ID
sleep $T

ln $P/dir_$ID/child $P/link_$ID
sleep $T

ln -s $P/link_$ID $P/sym_$ID
sleep $T

rm -r $P/dir_$ID
sleep $T

rm $P/file_$ID
sleep $T

ln $P/link_$ID $P/link2_$ID
sleep $T

rm $P/sym_$ID
sleep $T

rm $P/link2_$ID
sleep $T

rm $P/link_$ID
sleep $T

# signal completion
rmdir $P
