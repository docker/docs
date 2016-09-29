#!/bin/sh

R=0

rm -rf /tmp/dentry_test
mkdir -p /tmp/dentry_test

mkdir /tmp/dentry_test/d0
rmdir /tmp/dentry_test/d0
mkdir /tmp/dentry_test/d0
if [ $? != 0 ]; then
    echo "mkdir 1 FAIL"
    R=1
fi

mv /tmp/dentry_test/d0 /tmp/dentry_test/d1
mkdir /tmp/dentry_test/d0
if [ $? != 0 ]; then
    echo "mkdir 2 FAIL"
    R=1
fi

touch /tmp/dentry_test/d0/f
mv /tmp/dentry_test/d0 /tmp/dentry_test/d1
rm -rf /tmp/dentry_test/d1
if [ $? != 0 ]; then
    echo "rm 1 FAIL"
    R=1
fi

touch /tmp/dentry_test/f0
rm /tmp/dentry_test/f0
touch /tmp/dentry_test/f0
if [ $? != 0 ]; then
    echo "touch 1 FAIL"
    R=1
fi

touch /tmp/dentry_test/f1
rm /tmp/dentry_test/f0
ln /tmp/dentry_test/f1 /tmp/dentry_test/f0
if [ $? != 0 ]; then
    echo "ln 1 FAIL"
    R=1
fi

rm /tmp/dentry_test/f0
ln -s /tmp/dentry_test/f1 /tmp/dentry_test/f0
if [ $? != 0 ]; then
    echo "ln 2 FAIL"
    R=1
fi

rm /tmp/dentry_test/f0
touch /tmp/dentry_test/f0
mv /tmp/dentry_test/f0 /tmp/dentry_test/f1
touch /tmp/dentry_test/f0
if [ $? != 0 ]; then
    echo "touch 2 FAIL"
    R=1
fi

mv /tmp/dentry_test/f0 /tmp/dentry_test/f1
ln /tmp/dentry_test/f1 /tmp/dentry_test/f0
if [ $? != 0 ]; then
    echo "ln 3 FAIL" && false
    R=1
fi

rm /tmp/dentry_test/f1
mv /tmp/dentry_test/f0 /tmp/dentry_test/f1
ln -s /tmp/dentry_test/f1 /tmp/dentry_test/f0
if [ $? != 0 ]; then
    echo "ln 4 FAIL" && false
    R=1
fi

exit $R
