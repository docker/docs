---
title: Hello Build
description: Build Hello World
keywords: build, buildx, buildkit, getting started, dockerfile, image layers, build instructions, build context
---

## Hello Build!

It all starts with a Dockerfile.

Dockerfiles are text files containing instructions. Dockerfiles adhere to a specific format and contain a **set of instructions** for which you can find a full reference in the [Dockerfile reference](../../engine/reference/builder).  
Docker builds images by reading the instructions from a Dockerfile.

Docker images consist of **read-only layers**, each resulting from an instruction in the Dockerfile. Layers are stacked sequentially and each one is a delta representing the changes applied to the previous layer.

## Dockerfile basics

A Dockerfile is a text file containing all necessary instructions needed to assemble your application into a Docker container image.

Here are the most common types of instructions:
 
* [**FROM \<image\>**](../../engine/reference/builder/#from) - defines a base for your image.
* [**RUN \<command\>**](../../engine/reference/builder/#run) - executes any commands in a new layer on top of the current image and commits the result.
  RUN also has a shell form for running commands.
* [**WORKDIR \<directory\>**](../../engine/reference/builder/#workdir) - sets the working directory for any RUN, CMD, ENTRYPOINT, COPY, and ADD instructions that follow it in the Dockerfile.
* [**COPY \<src\> \<dest\>**](../../engine/reference/builder/#copy) - copies new files or directories from `<src>` and adds them to the filesystem of the container at the path `<dest>`.
* [**CMD \<command\>**](../../engine/reference/builder/#cmd) - lets you define the default program that is run once you start the container based on this image.
Each Dockerfile only has one CMD, and only the last CMD instance is respected when multiple exist.
  
Dockerfiles are crucial inputs for image builds and can facilitate automated, multi-layer image builds based on your unique configurations. Dockerfiles can start simple and grow with your needs and support images that require complex instructions.  
For all the possible instructions, see the [Dockerfile reference](../../engine/reference/builder/).

## Example
Here’s a simple Dockerfile example to get you started with building images. We’ll take a simple "Hello World" Python Flask application, and bundle it into a Docker image that we can test locally or deploy anywhere!

**Sample A**  
Let’s say we have the following in a `hello.py` file in our local directory:

```python
from flask import Flask
app = Flask(__name__)

@app.route("/")
def hello():
    return "Hello World!"
```

Don’t worry about understanding the full example if you’re not familiar with Python - it’s just a simple web server that will contain a single page that says “Hello World”.

> **Note** 
>
> If you test the example, make sure to copy over the indentation as well! For more information about this sample Flask application, check the [Flask Quickstart](https://flask.palletsprojects.com/en/2.1.x/quickstart/){:target="_blank" rel="noopener" class="_"} page.


**Sample B**  
Here’s a Dockerfile that Docker Build can use to create an image for our application:

```dockerfile
# syntax=docker/dockerfile:1
FROM ubuntu:22.04

# install app dependencies
RUN apt-get update && apt-get install -y python3 python3-pip
RUN pip install flask==2.1.*

# install app
COPY hello.py /

# final configuration
ENV FLASK_APP=hello
EXPOSE 8000
CMD flask run --host 0.0.0.0 --port 8000
```

* `# syntax=docker/dockerfile:1` 

    This is our syntax directive. It pins the exact version of the dockerfile syntax we’re using. As a [best practice](../../develop/dev-best-practices/), this should be the very first line in all our Dockerfiles as it informs Buildkit the right version of the Dockerfile to use. 
      See also [Syntax](../../engine/reference/builder/#syntax).
    
    > **Note** 
    > 
    > Initiated by a `#` like regular comments, this line is treated as a directive when you are using BuildKit (default), otherwise it is ignored.


* `FROM ubuntu:22.04`

    Here the `FROM` instruction sets our base image to the 22.04 release of Ubuntu. All following instructions are executed on this base image, in this case, a Ubuntu environment. 
    The notation `ubuntu:22:04`, follows the `name:tag` standard for naming docker images. 
    When you build your image you use this notation to name your images and use it to specify any existing Docker image.
    There are many public images you can leverage in your projects.  
    Explore [Docker Hub](https://hub.docker.com/search?image_filter=official&q=&type=image){:target="_blank" rel="noopener" class="_"} to find out.

* `# install app dependencies`

    Comments in dockerfiles begin with the `#` symbol. 
    As your Dockerfile evolves, comments can be instrumental to document how your dockerfile works for any future readers and editors of the file.  
    See also the [FROM instruction](../../engine/reference/builder/#from) page in the Dockerfile reference.

* `RUN apt-get update && apt-get install -y python3 python3-pip`

    This `RUN` instruction executes a shell command in the build context. A build's context is the set of files located in the specified PATH or URL. In this example, our context is a full Ubuntu operating system, so we have access to its package manager, apt. The provided commands update our package lists and then, after that succeeds, installs python3 and pip, the package manager for Python.  
    See also the [RUN instruction](../../engine/reference/builder/#run) page in the Dockerfile reference.

* `RUN pip install flask`

    This second `RUN` instruction requires that we’ve installed pip in the layer before. After applying the previous directive, we can use the pip command to install the flask web framework. This is the framework we’ve used to write our basic “Hello World” application from above, so to run it in Docker, we’ll need to make sure it’s installed.  
    See also the [RUN instruction](../../engine/reference/builder/#run) page in the Dockerfile reference.

* `COPY hello.py /`

    This COPY instruction copies our `hello.py` file from the build’s context local directory into the root directory of our image. After this executes, we’ll end up with a file called `/hello.py` inside the image, with all the content of our local copy!  
    See also the [COPY instruction](../../engine/reference/builder/#copy) page in the Dockerfile reference.
    

* `ENV FLASK_APP=hello` 

    This ENV instruction sets a Linux environment variable we’ll need later. This is a flask-specific variable, that configures the command later used to run our `hello.py` application. Without this, flask wouldn’t know where to find our application to be able to run it.  
    See also the [ENV instruction](../../engine/reference/builder/#env) page in the Dockerfile reference.


* `EXPOSE 8000` 

    This EXPOSE instruction marks that our final image has a service listening on port 8000. This isn’t required, but it is a good practice, as users and tools can use this to understand what your image does.  
    See also the [EXPOSE instruction](../../engine/reference/builder/#expose) page in the Dockerfile reference.

* `CMD flask run --host 0.0.0.0 --port 8000`

    This CMD instruction sets the command that is run when the user starts a container based on this image. In this case we’ll start the flask development server listening on all hosts on port 8000.  
    [CMD instruction](../../engine/reference/builder/#cmd) page in the Dockerfile reference.

## Test the example

Go ahead and try this example in your local Docker installation or you can use [Play with Docker](https://labs.play-with-docker.com){:target="_blank" rel="noopener" class="_"} that provides you with a temporary Docker instance on the cloud.

To test this example:
1. Create a file hello.py with the content of sample A.
2. Create a file named Dockerfile without an extension with the contents of sample B.
3. From your Docker instance build it with `docker build -t test:latest .`

    Breaking down the docker build command:  
    * **`-t test:latest`** option specifies the name (required) and tag (optional) of the image we’re building.
    * **`.`** specifies the build context as the current directory. In this example, this is where build expects to find the Dockerfile and the local files the Dockerfile needs to access, in this case your python application.
      So, in accordance with the build command issued and how build context works, your Dockerfile and python app need to be on the same directory.

4. Run your newly built image with `docker run -p 8000:8000 test:latest` 
From your computer, open a browser and navigate to `http://localhost:8000` or, if you’re using [Play with Docker](https://labs.play-with-docker.com){:target="_blank" rel="noopener" class="_"}, click on Open Port. 

## Other resources

If you are interested in examples in other languages, such as GO, check out our [language-specific guides](../../language) in the Guides section.
