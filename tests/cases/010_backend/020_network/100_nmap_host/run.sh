#!/bin/sh
set -x
# run n parallel, fairly aggressive nmap against the specified IP address

IP=$1
CONC=$2
PORTS=$3


for i in $(seq 0 "$CONC"); do
    echo "$i" of "$CONC"
    nmap -n -T 5 --min-parallelism 50 --max-parallelism 100 --disable-arp-ping "$PORTS" "$IP" &
done

wait

ping -c 1 "$IP"
exit $?
