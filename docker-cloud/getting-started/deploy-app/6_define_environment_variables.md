---
description: Define environment variables
keywords: Python, service, environment, service
redirect_from:
- /docker-cloud/getting-started/python/6_define_environment_variables/
- /docker-cloud/getting-started/golang/6_define_environment_variables/
title: Define environment variables
---

Docker lets you store data such as configuration settings, encryption keys, and external resource addresses in environment variables. Docker Cloud makes it easy to define, share, and update the environment variables for your services.

At runtime, environment variables are exposed to the application inside the container.

## Look inside your deployed app

Let's look inside the app you just deployed.

### Python quickstart

Open the file in `quickstart-python/app.py`, and look at the return statement in the method *hello()*. The code uses **os.getenv('NAME', "world")** to get the environment variable
**NAME**.

```python
return html.format(name=os.getenv('NAME', "world"), hostname=socket.gethostname(), visits=visits)
```

### Go quickstart

Open the file in `quickstart-go/main.go`, and look at the *fmt.Fprintf* call in the *indexHandler* method. The code uses **os.Getenv("NAME")** to get the environment variable **NAME**.

```go
fmt.Fprintf(w, "<h1>hello, %s</h1>\n<b>Hostname: </b>%s<br><b>MongoDB Status: </b>%s", os.Getenv("NAME"), hostname, mongostatus)
```

## Edit an environment variable

If you modify the environment variable, the message the app shows when you curl or visit the service webpage changes accordingly. Let's try it!

Run the following command to change the **NAME** variable, and then redeploy the `web` service.

```bash
$ docker-cloud service set --env NAME="Friendly Users" --redeploy web
```

## Check endpoint status

Execute `docker-cloud container ps` again to see the container's new endpoint.
You should now see two `web-1` containers, one with a status of **terminated**
(that's the original container) and another one either **starting** or already
**running**.

```none
$ docker-cloud container ps
NAME                         UUID      STATUS        IMAGE                                          RUN COMMAND      EXIT CODE  DEPLOYED        PORTS
web-1                        a2ff2247  ✘ Terminated  my-username/quickstart-python:latest           python app.py               40 minutes ago  web-1.my-username.cont.dockerapp.io:49165->80/tcp
web-1                        ae20d960  ▶ Running     my-username/quickstart-python:latest           python app.py               20 seconds ago  web-1.my-username.cont.dockerapp.io:49166->80/tcp
```

Now curl the new endpoint to see the updated greeting.

> **Note**: If `docker-cloud container ps` doesn't show an endpoint for the container yet, wait until the container status changes to **running**.

```none
$ curl web-1.$DOCKER_ID_USER.cont.dockerapp.io:49162
Hello Friendly Users!</br>Hostname: e360d05cdb81</br>Counter: Redis Cache not found, counter disabled.%
```

Your service now returns `Hello Friendly Users!`. Great! You've modified your service using environment variables!

### Environment Variables and the Dockerfile

Environment variables can also be set in the Dockerfile, and modified at runtime
(like you just did).

Wondering where the default value for the **NAME** environment variable is set?
Look in the quickstart's Dockerfile.

```none
# Environment Variables
ENV NAME World
```

## What's Next?

Next, we [scale the service](7_scale_the_service.md).
