#!/bin/sh
# SUMMARY: Check if we have sudo access, ask for password if necessary
# LABELS: installer
# AUTHOR: Magnus Skjegstad <magnus.skjegstad@docker.com>

set -e

sudo -n ls || echo "WARNING: sudo requires password, will not work unattended"
sudo ls
