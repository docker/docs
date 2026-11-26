---
title: Time an install into local Kubernetes
linkTitle: Kubernetes install
description: Measure how long it takes to install the OpenTelemetry Demo into a local Kubernetes cluster until a test checkout succeeds.
keywords: docker desktop, performance, benchmark, kubernetes, helm, opentelemetry demo, startup
weight: 80
---

This check measures the wait between installing an application into a local
Kubernetes cluster and the application working. It uses the OpenTelemetry Demo,
a web shop of 19 services that includes Kafka and a cache, installed from its
official Helm chart at a pinned version.

Image downloads stay out of the timing. You install the application once to pull
the images onto the node, then remove it. The timed part starts with an empty
namespace and ends when a purchase succeeds through the running shop: view a
product, add it to the cart, read the cart back, and check out. Pods becoming
ready doesn't end the timer. The checkout does.

The result is wall-clock time on your machine.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/).

## Additional prerequisites

- A single-node Kubernetes cluster from the engine you're measuring, with
  `kubectl` pointing at it. Give the cluster several GB of memory. 19 services,
  some on the JVM, don't fit in a small cluster.
- `helm`, `curl`, and `jq`.

## Set up

1. Download the pinned chart and verify its digest. If the digest doesn't match,
   stop. You don't have the chart the published numbers came from.

   ```console
   $ mkdir -p ~/perfcheck/kubernetes-startup && cd ~/perfcheck/kubernetes-startup

   $ helm repo add open-telemetry \
     https://open-telemetry.github.io/opentelemetry-helm-charts --force-update
   $ helm repo update open-telemetry
   $ helm pull open-telemetry/opentelemetry-demo --version 0.35.1

   $ echo 'cd0650001eead436ad7af3e2a1ff94ebb019039dcf0de76f49f2ce0fd5b65aee  opentelemetry-demo-0.35.1.tgz' | shasum -a 256 --check

   $ tar -xzf opentelemetry-demo-0.35.1.tgz
   ```

   On Linux, use `sha256sum -c` instead of `shasum -a 256 --check`.

2. Write the purchase details for the timed checkout request. The card number is
   the demo's test value, and nothing leaves your cluster.

   ```console
   $ cat > checkout.json <<'EOF'
   {
     "email": "perfcheck@example.com",
     "address": {
       "streetAddress": "1 Benchmark Way",
       "zipCode": "94105",
       "city": "San Francisco",
       "state": "CA",
       "country": "United States"
     },
     "userCurrency": "USD",
     "creditCard": {
       "creditCardNumber": "4432-8015-6152-0454",
       "creditCardExpirationMonth": 1,
       "creditCardExpirationYear": 2039,
       "creditCardCvv": 672
     }
   }
   EOF
   ```

3. Define a helper with the install options, so the untimed and timed installs
   are identical. The options turn off the demo's traffic generator and
   monitoring interfaces, which the purchase doesn't use and which add load.

   ```bash
   otel_helm() {
     helm "$@" otel-demo ./opentelemetry-demo \
       --namespace perfcheck-otel --create-namespace \
       --set components.loadgenerator.enabled=false \
       --set-json 'components.flagd.sidecarContainers=[]' \
       --set grafana.enabled=false \
       --set jaeger.enabled=false \
       --set prometheus.enabled=false \
       --set opensearch.enabled=false \
       --wait --timeout 5m
   }
   ```

   Confirm that the chart describes the right application:

   ```console
   $ otel_helm template | grep -c '^kind: Deployment'
   19
   ```

## Run the check

1. Pull the images onto the node, outside the timing. Install the application,
   remove it, and wait for the namespace to disappear. The first run downloads
   several GB.

   ```console
   $ kubectl delete namespace perfcheck-otel --ignore-not-found=true --wait=true --timeout=2m
   $ otel_helm upgrade --install
   $ helm uninstall otel-demo --namespace perfcheck-otel --wait --timeout 2m
   $ kubectl delete namespace perfcheck-otel --ignore-not-found=true --wait=true --timeout=2m
   ```

   The last command must finish. Installing into a namespace that's still
   terminating invalidates the run.

2. Run the timed command in bash or zsh:

   ```console
   $ time (
     otel_helm upgrade --install || exit 1

     kubectl port-forward --namespace perfcheck-otel \
       service/frontendproxy 18080:8080 > port-forward.log 2>&1 &
     port_forward_pid=$!
     trap 'kill "$port_forward_pid" 2>/dev/null' EXIT

     attempt=0
     until curl -fsS --connect-timeout 1 --max-time 10 \
         http://127.0.0.1:18080/api/products/OLJCESPC7Z > product.json &&
         jq -e '.id == "OLJCESPC7Z"' product.json >/dev/null; do
       attempt=$((attempt + 1)); [ "$attempt" -lt 30 ] || exit 1; sleep 0.2
     done

     curl -fsS --connect-timeout 1 --max-time 10 -H 'Content-Type: application/json' \
       -d '{"item":{"productId":"OLJCESPC7Z","quantity":1},"userId":"perfcheck"}' \
       http://127.0.0.1:18080/api/cart >/dev/null &&
     curl -fsS --connect-timeout 1 --max-time 10 \
       'http://127.0.0.1:18080/api/cart?sessionId=perfcheck&currencyCode=USD' > cart.json &&
     jq -e '.items[0].productId == "OLJCESPC7Z" and .items[0].quantity == 1' cart.json >/dev/null &&
     jq --arg user_id perfcheck '. + {userId: $user_id}' checkout.json > checkout.request.json &&
     curl -fsS --connect-timeout 1 --max-time 30 -H 'Content-Type: application/json' \
       --data-binary @checkout.request.json \
       http://127.0.0.1:18080/api/checkout > checkout.response.json &&
     jq -e '.orderId != null and .shippingTrackingId != null' checkout.response.json >/dev/null
   )
   ```

3. Record the elapsed time. A non-zero exit status means the shop never became
   usable, and the run doesn't count.

4. Check that the run is valid:

   ```console
   $ jq '.orderId, .shippingTrackingId' checkout.response.json
   $ kubectl get deployment --namespace perfcheck-otel \
     --selector app.kubernetes.io/instance=otel-demo --no-headers | wc -l
   19
   ```

   Both values from the first command must be present, not `null`. If the
   connection to the shop dies during a run, `port-forward.log` shows why. That
   run is invalid.

5. Clean up. Do this even if the run failed, or later measurements on the
   machine are invalid.

   ```console
   $ pkill -f 'port-forward.*frontendproxy' 2>/dev/null
   $ helm uninstall otel-demo --namespace perfcheck-otel --ignore-not-found --wait --timeout 2m
   $ kubectl delete namespace perfcheck-otel --ignore-not-found=true --wait=true --timeout=2m
   ```
