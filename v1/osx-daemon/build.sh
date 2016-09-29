#!/bin/sh
set -ex

eval `opam config env`
oasis setup
make
ocamlfind remove osx-daemon || true
make install
