#!/bin/sh

# Fail hard and fast
set -eo pipefail

# If this fails, docker will restart the container. Yay, docker.
confd -node https://dtr-etcd-${DTR_REPLICA_ID}.dtr-br:2379 -node https://dtr-etcd-${DTR_REPLICA_ID}.dtr-br:4001 -onetime -config-file /etc/confd/confd.toml

# Run confd watcher in the background to watch the upstream servers
confd -node https://dtr-etcd-${DTR_REPLICA_ID}.dtr-br:2379 -node https://dtr-etcd-${DTR_REPLICA_ID}.dtr-br:4001 -config-file /etc/confd/confd.toml &
echo "[nginx] confd is listening for changes on etcd..."

# Start dnsfix to restart nginx whenever dns names change
echo "[nginx] starting dnsfix service..."
/dnsfix &

# Start nginx
echo "[nginx] starting nginx service..."
exec nginx -c /config/nginx.conf -g "daemon off;"
