#!/bin/bash
set -e

cd $(dirname $0)
ROOT=$(dirname $0)/../..
HW=$(ls -1 ${ROOT}/results)
mkdir -p output
rm -f output/*
for i in ${HW}; do
  find ${ROOT}/results/${i}/iperf-to-container/native.{xhyve,hyperkit} -name results.dat | xargs cat >> output/${i}.native.dat
  find ${ROOT}/results/${i}/iperf-to-container/slirp.{xhyve,hyperkit} -name results.dat | xargs cat >> output/${i}.slirp.dat
done

cat > output/plot.gp <<EOT
set terminal png
set output "output/plot.png"
set title "iperf throughput comparison"
set ylabel 'Throughput in Mbit/s'
set xdata time
set timefmt "%s"
set xlabel 'Time of experiment'
set key left top
EOT

echo -n "plot " >> output/plot.gp
COMMA=0
for i in ${HW}; do
  if [ -e ${ROOT}/labels/hardware/${i} ]; then
    label=$(cat ${ROOT}/labels/hardware/${i})
  else
    label=${i}
  fi
  if [ "$COMMA" -eq 1 ]; then
    echo "\\" >> output/plot.gp
    echo -n "," >> output/plot.gp
  fi
  echo "\"output/${i}.native.dat\" using 1:2 title \"${label} via vmnet\" with points\\" >> output/plot.gp
  echo -n ",\"output/${i}.slirp.dat\" using 1:2 title \"${label} via slirp\" with points " >> output/plot.gp
  COMMA=1
done
