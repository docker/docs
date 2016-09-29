#!/bin/sh
set -ex

eval `opam config env`
oasis setup
make
ocamlfind remove osx-hyperkit || true
make install
