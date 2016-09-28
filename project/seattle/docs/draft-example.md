# Draft example

This is an example draft that you can use to start creating documentation.
We'll take what you've created and transform it into public-facing docs.

A draft has:
* The user story
* Explanation of important concepts the user might not be familiar with
* Steps and commands the user needs to execute the story

## User story

As a developer, I want to push an image to DTR, in order to share new
developments with my team and no one else.

## Configure your Docker Engine

By default Docker Engine uses TLS when pushing and pulling images to an
image registry like Docker Trusted Registry.

If DTR is using the default configurations or was configured to use
self-signed certificates, you need to configure your Docker Engine to trust DTR.
Otherwise, when you try to login or push and pull images to DTR, you'll get an
error:

```bash
$ docker login <dtr-domain-name>

x509: certificate signed by unknown authority
```

See here how to configure your Docker daemon.

## Create a repository

To create a new repository, navigate to the DTR web application, and click
the 'New repository' button.

[We'll probably need an image here to make it obvious]

Add a name and description for the repository, and choose whether your
repository is public or private:

  * Public repositories are visible to all users, but can only be changed by
  users granted with permission to write them.
  * Private repositories can only be seen by users that have been granted
  permissions to that repository.

Click 'Save' to create the repository.

## Tag the image

Before you can push an image to DTR, you need to tag it with the full
repository name.

```
# Pull from Docker Hub the 1.7 tag of the golang image
$ docker pull golang:1.7
# Tag the golang:1.7 image with the full repository name we've created in DTR
$ docker tag golang:1.7 dtr.company.org/dave.lauper/golang:1.7
```

## Push the image

Now that you have tagged the image, you only need to authenticate and push the
image to DTR.

```bash
$ docker login dtr.company.org
$ docker push dtr.company.org/dave.lauper/golang:1.7
```
