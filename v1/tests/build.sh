#!/bin/sh
set -ex

eval `opam config env`
make depends
make
