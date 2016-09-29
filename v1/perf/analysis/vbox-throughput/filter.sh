#!/bin/bash
set -e

cd $(dirname $0)
ROOT=$(dirname $0)/../..
HW=$(ls -1 ${ROOT}/results)
mkdir -p output
rm -f output/*
for i in ${HW}; do
  find ${ROOT}/results/${i}/iperf-to-container/native.{xhyve,hyperkit} -name results.dat | xargs cat >> output/${i}.to.vmnet.dat
  find ${ROOT}/results/${i}/iperf-to-container/unknown.virtualbox -name results.dat | xargs cat >> output/${i}.to.vbox.dat
  find ${ROOT}/results/${i}/iperf-from-container/native.{xhyve,hyperkit} -name results.dat | xargs cat >> output/${i}.from.vmnet.dat
  find ${ROOT}/results/${i}/iperf-from-container/unknown.virtualbox -name results.dat | xargs cat >> output/${i}.from.vbox.dat
done
