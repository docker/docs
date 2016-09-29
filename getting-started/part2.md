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

Normally if you were to start writing a Python app on your laptop, your first
order of business would be to install a Python runtime onto your machine. But,
that creates a situation where the environment on your machine has to be just so
in order for your app to run as expected.

In Docker, you can just grab an image of Python runtime that is already set up,
and use that as a base for creating your app. Then, your build can include the
base Python image right alongside your app code, ensuring that your app and the
runtime it needs to run all travel together.

It's done with something called a Dockerfile.

## Your first Dockerfile

Create a folder and put this file in it, with the name `Dockerfile` (no
extension). This Dockerfile defines what goes on in the environment inside your
container. Things are virtualized inside this environment, which is isolated
from the rest of your system, so you have to map ports to the outside world, and
be specific about what files you want to "copy in" to that environment. However,
after doing that, you can expect that the build of your app with this
`Dockerfile` will behave exactly the same wherever it runs.

{% gist johndmulhausen/c31813e076827178216b74e6a6f4a087 %}

This `Dockerfile` refers to a couple of things we haven't created yet, namely
`app.py` and `requirements.txt`. We'll get there. But here's what this
`Dockerfile` is saying:

- Go get the base Python 2.7 runtime
- Create `/app` and set it as the current working directory inside the container
- Copy the contents of my current directory (on my machine) into `/app` (in this container image)
- Install any Python packages that I list inside what is now `/app/requirements.txt` inside the container
- Ensure that this container has port 80 open when it runs
- Set an environment variable within this container named `NAME` to be the string `World`
- Finally, when the container runs, execute `python` and pass in what is now `/app/app.py`

This paradigm is how developing with Docker essentially works. Make a
`Dockerfile` that includes the base image, grabs your code, installs
dependencies, initializes variables, and runs the command.

### The app itself

Grab these two files that were referred to in the above `Dockerfile` and place
them together with `Dockerfile`, all in the same folder.

{% gist johndmulhausen/074cc7f4c26a9a8f9164b20b22602ad7 %}
{% gist johndmulhausen/8728902faede400c057f3205392bb9a8 %}

You're probably getting the picture by now. In `Dockerfile` we told the `pip`
package installer to install whatever was in `requirements.txt`, which we
now see is the Flask and Redis libraries for Python. The app itself is going to
print the environment variable of `NAME`, which we set as `World`, as well as
the output of a call to `socket.gethostname()`, which the Docker runtime is
going to answer with the container ID. Finally, because Redis isn't running
(we've only installed the Python library), we should expect that the attempt to
use it here will fail and show the error message.

## Build the App

That's it! You don't need to have installed Python or anything in
`requirements.txt` on your system, nor will running this app install them in
your system. It doesn't seem like you've really set up an environment with
Python and Flask, but you have. Let's build and run your app and prove it.

Make sure you're in the directory where you saved the three files we've shown,
and you've got everything.

```shell
$ ls
Dockerfile		app.py			requirements.txt
```

Now run the build command. This creates a Docker image, which we're going to
tag using `-t` so it has a friendly name, which you can use interchangeable
with the image ID in commands.

```shell
docker build -t "friendlyhello" .
```

In the output spew you can see everything defined in the `Dockerfile` happening,
including the installation of the packages we specified in `requirements.txt`.
Where is your built image? It's in your machine's local Docker image registry.
Check it out:

```shell
$ docker images
REPOSITORY            TAG                 IMAGE ID            CREATED             SIZE
friendlyhello         latest              326387cea398        47 seconds ago      192.1 MB
```

## Run the app

We're going to run the app and route traffic from our machine's port 80 to the
port 80 we exposed

```shell
docker run -p 80:80 friendlyhello
```

You should see a notice that Python is serving your app at `http://0.0.0.0:80`.
You can go there, or just to `http://localhost`, and see your app, "Hello World"
text, the container ID, and the Redis error message, all printed out in
beautiful Times New Roman.

Hit `CTRL+C` and let's run the app in the background, in detached mode.

```shell
docker run -d -p 80:80 friendlyhello
```

You get a hash ID of the container instance and then are kicked back to your
terminal. Your app is running in the background. Let's see it with `docker ps`:

```shell
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS
1fa4ab2cf395        friendlyhello       "python app.py"     28 seconds ago      Up 25 seconds
```

You'll see that `CONTAINER ID` matches what's on `http://localhost`, if you
refresh the browser page. You can't `CTRL+C` now, so let's kill the process this
way. Use the value you see under `CONTAINER ID`:

```shell
docker kill (containerID)
```

## Share the App

Now let's test how portable this app really is.

Sign up for Docker Hub at [https://hub.docker.com/](https://hub.docker.com/).
Make note of your username. We're going to use it in a couple commands.

Docker Hub is a public registry. A registry is a collection of accounts and
their various repositories. A repository is a collection of assets associated
with your account - like a GitHub repository, except the code is already built.


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

You'll see this stranger of a machine pull your image, along with Python and all
the dependencies from `requirements.txt`, and run your code. It all travels
together in a neat little package, and the new machine didn't have to install
anything but Docker to run it.

[On to "Getting Started, Part 3: Stateful, Multi-container Applications" >>](part3.md){: class="button darkblue-btn"}
