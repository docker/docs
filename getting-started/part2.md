---
title: "Getting Started, Part 2: Creating and Building Your App"
---

# Getting Started, Part 2: Creating and Building Your App

In [Getting Started, Part 1: Orientation and Setup](index.md), you heard an
overview of what containers are, what the
Docker platform does, and what we'll be covering in this multi-part tutorial.
You also got Docker installed on your machine.

In this section, you will write, build, run, and share an app, the Docker way.

## Your development environment

In the past, if you were to start writing a Python app, your first
order of business was to install a Python runtime onto your machine. But,
that creates a situation where the environment on your machine has to be just
so in order for your app to run as expected; ditto for the server that runs
your app.

With Docker, you can just grab a portable Python runtime as an image, no
installation necessary. Then, your build can include the base Python image
right alongside your app code, ensuring that your app, its dependencies, and the
runtime, all travel together.

These builds are configured with something called a `Dockerfile`.

## Your first Dockerfile

Create an empty directory and put this file in it, with the name `Dockerfile`.
`Dockerfile` will define what goes on in the environment inside your
container. Access to resources like networking interfaces and disk drives is
virtualized inside this environment, which is isolated from the rest of your
system, so you have to map ports to the outside world, and
be specific about what files you want to "copy in" to that environment. However,
after doing that, you can expect that the build of your app defined in this
`Dockerfile` will behave exactly the same wherever it runs.

{% gist johndmulhausen/c31813e076827178216b74e6a6f4a087 %}

This `Dockerfile` refers to a couple of things we haven't created yet, namely
`app.py` and `requirements.txt`. We'll get there. But here's what this
`Dockerfile` is saying:

- Download the official image of the Python 2.7 runtime and include it here.
- Create `/app` and set it as the current working directory inside the container
- Copy the contents of the current directory on my machine into `/app` inside the container
- Install any Python packages that I list inside `requirements.txt`
- Ensure that port 80 is exposed to the world outside this container
- Set an environment variable within this container named `NAME` to be the string `World`
- Finally, execute `python` and pass in `app.py` as the "entry point" command,
  the default command that is executed at runtime.

### The app itself

Grab these two files and place them in the same folder as `Dockerfile`.

{% gist johndmulhausen/074cc7f4c26a9a8f9164b20b22602ad7 %}
{% gist johndmulhausen/8728902faede400c057f3205392bb9a8 %}

Now we see that the `Dockerfile` command `pip install requirements.txt` installs
the Flask and Redis libraries for Python. We can also see that app itself
prints the environment variable of `NAME`, which we set as `World`, as well as
the output of a call to `socket.gethostname()`, which the Docker runtime is
going to answer with the container ID, which is sort of like the process ID for
an executable. Finally, because Redis isn't running
(as we've only installed the Python library, and not Redis itself), we should
expect that the attempt to use it here will fail and produce the error message.

## Build the App

That's it! You don't need to have installed Python or anything in
`requirements.txt` on your system, nor will running this app install them in
your system. It doesn't seem like you've really set up an environment with
Python and Flask, but you have. Let's build and run your app and prove it.

7Here's what `ls` should show:

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

Sign up a Docker account at [hub.docker.com](https://hub.docker.com/).
Make note of your username. We're going to use it in a couple commands.

Docker Hub is a public registry. A registry is a collection of accounts and
their various repositories. A repository is a collection of tagged images like a
GitHub repository, except the code is already built.

Log in your local machine to Docker Hub.

```shell
docker login
```

Now, let's publish your image. First, specify the repository you'd like to use
in a tag. The notation for associating a local image with a repository on a
registry, is `username/repository:tag`. The `:tag` is optional, but recommended;
it's the mechnism that registries use to give Docker images a version. So,
putting all that together:

```shell
docker tag friendlyhello YOURUSERNAME/YOURREPO:ARBITRARYTAG
```

From now on, you can use `docker run` on this machine with the fully qualified
tag. But that won't work on other machines until you upload this image, like so:

```shell
docker push YOURUSERNAME/YOURREPO:ARBITRARYTAG
```

Once complete, the results of this upload are [publicly available
on Docker Hub](https://hub.docker.com/).

Now, remembering whatever you specified as your target repo, and whatever you
used as a tag, go on another machine. Any machine where you can install Docker,
and run this command:

```shell
docker run YOURUSERNAME/YOURREPO:ARBITRARYTAG
```

> Note: If you don't specify the `:ARBITRARYTAG` portion of these commands,
  the tag of `:latest` will be assumed, both when you build and when you run
  images.

You'll see this stranger of a machine pull your image, along with Python and all
the dependencies from `requirements.txt`, and run your code. It all travels
together in a neat little package, and the new machine didn't have to install
anything but Docker to run it.

## Recap and cheat sheet for images and containers

To recap: After calling `docker run`, you created and ran a container, based on
the image created when you called `docker build`. Images are defined in a
`Dockerfile`. A container is an instance of an image, and it has any package
installations, file writes, etc that happen after you call `docker run` and run
the app. And lastly, images are shared via a registry.

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
docker login #Log in this CLI session using your Docker credentials (to Docker Hub by default)
docker tag <image> username/repository:tag #Tag <image> on your local machine for upload
docker push username/repository:tag #Upload tagged image to registry (Docker Hub by default)
docker run username/repository:tag #Run image from a registry (Docker Hub by default)
```

[On to "Getting Started, Part 3: Stateful, Multi-container Applications" >>](part3.md){: class="button darkblue-btn"}
