---
title: "Getting Started, Part 2: Creating and Building Your App"
---

In the past, if you were to start writing a Python app, your first order of
business was to install a Python runtime onto your machine. With Docker, your
app, its dependencies, and the runtime, are all built together into an image.

These builds are configured with something called a `Dockerfile`.

## Your first Dockerfile

`Dockerfile` will define what goes on in the environment inside your
container. Access to resources like networking interfaces and disk drives is
virtualized inside this environment, which is isolated from the rest of your
system, so you have to map ports to the outside world, and
be specific about what files you want to "copy in" to that environment. However,
after doing that, you can expect that the build of your app defined in this
`Dockerfile` will behave exactly the same wherever it runs.

### `Dockerfile`

Create an empty directory and put this file in it, with the name `Dockerfile`.

```
# Using official python runtime base image
FROM python:2.7-slim

# Set the working directory to /app
WORKDIR /app

# Copy the host machine directoy's contents to /app in the container
ADD . /app

# Install any needed packages specified in requirements.txt
RUN pip install -r requirements.txt

# Make port 80 available to the world outside this container
EXPOSE 80

# Define environment variable
ENV NAME World

# Run app.py when the container launches
CMD ["python", "app.py"]
```

This `Dockerfile` refers to a couple of things we haven't created yet, namely
`app.py` and `requirements.txt`. Let's get those in place next.

## The app itself

Grab these two files and place them in the same folder as `Dockerfile`.

### `requirements.txt`

```
Flask
Redis
```

### `app.py`

```python
from flask import Flask
from redis import Redis, RedisError
import os
import socket

# Connect to Redis
redis = Redis(host="redis", db=0)

app = Flask(__name__)


@app.route("/")
def hello():
    try:
        visits = redis.incr('counter')
    except RedisError:
        visits = "<i>cannot connect to Redis, counter disabled</i>"

    html = "<h3>Hello {name}!</h3>" \
           "<b>Hostname:</b> {hostname}<br/>" \
           "<b>Visits:</b> {visits}"
    return html.format(name=os.getenv('NAME', "world"), hostname=socket.gethostname(), visits=visits)


if __name__ == "__main__":
	app.run(host='0.0.0.0', port=80)
```

Now we see that `pip install requirements.txt` installs the Flask and Redis
libraries, and the app prints the environment variable of `NAME`, as well as the
output of a call to `socket.gethostname()`. Finally, because Redis isn't running
(as we've only installed the Python library, and not Redis itself), we should
expect that the attempt to use it here will fail and produce the error message.

> *Note*: Accessing the name of the host when inside a container retrieves the
container ID, which is like the process ID for a running executable.

## Build the App

That's it! You don't need Python or anything in `requirements.txt` on your
system, nor will running this app install them. It doesn't seem like you've
really set up an environment with Python and Flask, but you have. Let's build
and run your app and prove it.

Here's what `ls` should show:

```shell
$ ls
Dockerfile		app.py			requirements.txt
```

Now run the build command. This creates a Docker image, which we're going to
tag using `-t` so it has a friendly name.

```shell
docker build -t friendlyhello .
```

In the output spew you can see everything defined in the `Dockerfile` happening.
Where is your built image? It's in your machine's local Docker image registry.
Check it out:

```shell
$ docker images
REPOSITORY            TAG                 IMAGE ID            CREATED             SIZE
friendlyhello         latest              326387cea398        47 seconds ago      192.1 MB
```

## Run the app

Run the app, mapping our machine's port 4000 to the container's exposed port 80
using `-p`:

```shell
docker run -p 4000:80 friendlyhello
```

You should see a notice that Python is serving your app at `http://0.0.0.0:80`.
But that message coming from inside the container, which doesn't know you
actually want to access your app at: `http://localhost:4000`. Go there, and
you'll see the "Hello World" text, the container ID, and the Redis error
message, all printed out in beautiful Times New Roman.

Hit `CTRL+C` in your terminal to quit.

Now let's run the app in the background, in detached mode:

```shell
docker run -d -p 4000:80 friendlyhello
```

You get a hash ID of the container instance and then are kicked back to your
terminal. Your app is running in the background. Let's see it with `docker ps`:

```shell
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS
1fa4ab2cf395        friendlyhello       "python app.py"     28 seconds ago      Up 25 seconds
```

You'll see that `CONTAINER ID` matches what's on `http://localhost:4000`, if you
refresh the browser page. Now use `docker stop` to end the process, using
`CONTAINER ID`, like so:

```shell
docker stop 1fa4ab2cf395
```

## Share your image

Sign up for a Docker account at [cloud.docker.com](https://cloud.docker.com/).
Make note of your username. We're going to use it to upload our build to Docker
Store and make it retrievable from anywhere.

Docker Store is a public registry, and a registry is a collection of
repositories. A repository is a collection of tagged images like a GitHub
repository, except the code is already built.

Log in your local machine to Docker Store.

```shell
docker login
```

Now, let's publish your image. The notation for associating a local image with a
repository on a registry, is `username/repository:tag`. The `:tag` is optional,
but recommended; it's the mechnism that registries use to give Docker images a
version. So, putting all that together:

```shell
docker tag friendlyhello YOURUSERNAME/YOURREPO:ARBITRARYTAG
```

Upload this image:

```shell
docker push YOURUSERNAME/YOURREPO:ARBITRARYTAG
```

Once complete, the results of this upload are [publicly available
on Docker Hub](https://store.docker.com/). From now on, you can use `docker run`
on any machine with the fully qualified tag:

```shell
docker run YOURUSERNAME/YOURREPO:ARBITRARYTAG
```

> Note: If you don't specify the `:ARBITRARYTAG` portion of these commands,
  the tag of `:latest` will be assumed, both when you build and when you run
  images.

No matter where `docker run` executes, it pulls your image, along with Python
and all the dependencies from `requirements.txt`, and runs your code. It all
travels together in a neat little package, and the host machine doesn't have to
install anything but Docker to run it.

That's all for this page. You can continue to the next phase where we link up
Redis to our web page so that counter works. Or, if you want a quick recap
first, scroll down for the [cheat sheet](#cheat-sheet).

[On to "Getting Started, Part 3: Stateful, Multi-container Applications" >>](part3.md){: class="button outline-btn"}

## Cheat sheet

Here's [a terminal recording of everything that we did on this page](https://asciinema.org/a/blkah0l4ds33tbe06y4vkme6g).

<script type="text/javascript" src="https://asciinema.org/a/blkah0l4ds33tbe06y4vkme6g.js" id="asciicast-blkah0l4ds33tbe06y4vkme6g" async></script>

And here is a list of basic commands from this page that deal with our same

```shell
docker build -t friendlyname . #Create image using this directory's Dockerfile
docker run -p 4000:80 friendlyname #Run image "friendlyname" mapping port 4000 to 80
docker run -d -p 4000:80 friendlyname #Same thing, but in detached mode
docker ps #See a list of all running containers
docker stop <hash> #Gracefully stop the specified container
docker ps -a #See a list of all containers on this machine, even the ones not running
docker kill <hash> #Force shutdown of the specified container
docker rm <hash> #Remove the specified container from this machine
docker rm $(docker ps -a -q) #Remove all containers from this machine
docker images -a #Show all images that have been built or downloaded onto this machine
docker rmi <imagename> #Remove the specified image from this machine
docker rmi $(docker images -q) #Remove all images from this machine
docker login #Log in this CLI session using your Docker credentials
docker tag <image> username/repository:tag #Tag <image> on your local machine for upload
docker push username/repository:tag #Upload tagged image to registry
docker run username/repository:tag #Run image from a registry
```
