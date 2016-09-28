#!/bin/sh

service rsyslog start
touch /var/log/syslog
tail -f /var/log/syslog
