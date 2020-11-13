---
title: "Use containers for development"
keywords: get started, NodeJS, local, development
description: Learn how to develop your application locally.
---

{% include_relative nav.html selected="3" %}

## Prerequisites

Work through the steps to build an image and run it as a containerized application in [Run your image as a container](run-containers.md).

## Introduction

In this module, we’ll walk through setting up a local development environment for the application we built in the previous modules. We’ll use Docker to build our images and Docker Compose to make everything a whole lot easier.

## Use Compose to develop locally

The notes-service project uses MongoDB as its data store. If you remember from Part I of this series, we had to start the Mongo container manually and connect it to the same network that our notes-service is running on. We also had to create a couple of volumes so we could persist our data across restarts of our application and MongoDB.

In this section, we’ll create a Compose file to start our node-docker and the MongoDB with one command. We’ll also set up the Compose file to start the node-docker in debug mode so that we can connect a debugger to the running node process.

Open the notes-service in your IDE or text editor and create a new file named `docker-compose.dev.yml`. Copy and paste the below commands into the file.

```yaml
version: '3.8'

services:
 notes:
  build:
   context: .
  ports:
   - 8080:8080
   - 9229:9229
  environment:
   - SERVER_PORT=8080
   - DATABASE_CONNECTIONSTRING=mongodb://mongo:27017/notes
  volumes:
   - ./:/code
  command: npm run debug

 mongo:
  image: mongo:4.2.8
  ports:
   - 27017:27017
  volumes:
   - mongodb:/data/db
   - mongodb_config:/data/configdb
volumes:
 mongodb:
 mongodb_config:
```

This Compose file is super convenient as we do not have to type all the parameters to pass to the `docker run` command. We can declaratively do that in the Compose file.

We are exposing port 9229 so that we can attach a debugger. We are also mapping our local source code into the running container so that we can make changes in our text editor and have those changes picked up in the container.

One other really cool feature of using a Compose file, is that we have service resolution set up to use the service names. So we are now able to use `“mongo”` in our connection string. The reason we use mongo is because that is what we have named our mongo service in the Compose file as.

Let’s start our application and confirm that it is running properly.

```shell
$ docker-compose -f docker-compose.dev.yml up --build
```

We pass the `--build` flag so Docker will compile our image and then starts it.

If all goes will you should see something similar:

  ![node-compile](images/node-compile.png)

Now let’s test our API endpoint. Run the following curl command:

```shell
$ curl --request GET --url http://localhost:8080/services/m/notes
```

You should receive the following response:

```json
{"code":"success","meta":{"total":0,"count":0},"payload":[]}
```

## Connect a debugger

We’ll use the debugger that comes with the Chrome browser. Open Chrome on your machine and then type the following into the address bar.

`about:inspect`

It opens the following screen.

  ![Chrome-inspect](images/chrome-inspect.png)

Click the **Open dedicated DevTools for Node** link. This opens the DevTools that are connected to the running Node.js process inside our container.

Let’s change the source code and then set a breakpoint.

Add the following code to the server.js file on line 19 and save the file.

```node
 server.use( '/foo', (req, res) => {
   return res.json({ "foo": "bar" })
 })
```

If you take a look at the terminal where our Compose application is running, you’ll see that nodemon noticed the changes and reloaded our application.

 ![nodemon](images/nodemon.png)

Navigate back to the Chrome DevTools and set a breakpoint on line 20 and then run the following curl command to trigger the breakpoint.

```shell
$ curl --request GET --url http://localhost:8080/foo
```

You should have seen the code break on line 20 and now you are able to use the debugger just like you would normally. You can inspect and watch variables, set conditional breakpoints, view stack traces, etc.

## Conclusion

In this article, we took a look at creating a general development image that we can use pretty much like our normal command line. We also set up our Compose file to map our source code into the running container and exposed the debugging port.
