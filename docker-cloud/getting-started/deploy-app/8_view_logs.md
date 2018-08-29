---
description: View service logs
keywords: View, logs, Python
redirect_from:
- /docker-cloud/getting-started/python/8_view_logs/
- /docker-cloud/getting-started/golang/8_view_logs/
title: View service logs
notoc: true
---

Docker Cloud grants you access to the logs your application writes to `stdout`.
An internal service multiplexes all the logs from all the containers of a
service into a single stream. To see a service's logs run the `docker-cloud
service logs` command with the name of the service.

If you run `docker-cloud service logs web`, you see logs for both *web-1* and
*web-2*, like the example below.

```none
$ docker-cloud service logs web
[web-1] 2015-01-13T22:45:37.250431077Z  * Running on http://0.0.0.0:80/
[web-1] 2015-01-07T17:20:19.076174813Z 83.50.33.64 - - [07/Jan/2015 17:20:19] "GET / HTTP/1.1" 200 -
[web-1] 2015-01-07T17:20:34.209098162Z 83.50.33.64 - - [07/Jan/2015 17:20:34] "GET / HTTP/1.1" 200 -
[web-1] 2015-01-07T18:46:07.116759956Z 83.50.33.64 - - [07/Jan/2015 18:46:07] "GET / HTTP/1.1" 200 -
[web-2] 2015-01-07T18:48:24.550419508Z  * Running on http://0.0.0.0:5000/
[web-2] 2015-01-07T18:48:37.116759956Z 83.50.33.64 - - [07/Jan/2015 18:48:37] "GET / HTTP/1.1" 200 -
```

To see a specific container's logs, use the `container logs` and the
specific container's name. To learn more about service and container
hostnames, see [Service Discovery](../../apps/service-links.md#using-service-and-container-names-as-hostnames).

```none
$ docker-cloud container logs web-1
2015-01-07T17:18:24.550419508Z  * Running on http://0.0.0.0:80/
2015-01-07T17:20:19.076174813Z 83.50.33.64 - - [07/Jan/2015 17:20:19] "GET / HTTP/1.1" 200 -
2015-01-07T17:20:34.209098162Z 83.50.33.64 - - [07/Jan/2015 17:20:34] "GET / HTTP/1.1" 200 -
2015-01-07T18:46:07.116759956Z 83.50.33.64 - - [07/Jan/2015 18:46:07] "GET / HTTP/1.1" 200 -
```

Visit your application using curl or your browser again. Run the `service logs
web` command again to see another log message for your visit.

## What's Next?

Now, let's explore how to
[Load balance the service](9_load-balance_the_service.md).
