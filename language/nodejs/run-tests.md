---
title: "Run your Tests using Node.js and Mocha frameworks"
keywords: Node.js, build, Mocha, test
description: How to Build and Run your Tests using Node.js and Mocha frameworks
---

{% include_relative nav.html selected="4" %}

Testing is an essential part of modern software development. Testing can mean a lot of things to different development teams. There are unit tests, integration tests and end-to-end testing. In this guide we take a look at running your unit tests in Docker. Let's assume we have defined Mocha tests in a `./test` folder within our application.

### Running locally and testing the application

Now that we have our sample application locally, let’s build our Docker image and make sure everything is running properly. Run the following commands to build and then run your Docker image in a container.

```shell
$ docker build -t node-docker .
$ docker run -it --rm --name app -p 8080:80 node-docker
```

Now let’s test our application by POSTing a JSON payload and then make an HTTP GET request to make sure our JSON was saved correctly.

```shell
$ curl --request POST \
  --url http://localhost:8080/services/test \
  --header 'content-type: application/json' \
  --data '{
	"msg": "testing"
}'
```

Now, perform a GET request to the same endpoint to make sure our JSON payload was saved and retrieved correctly. The “id” and “createDate” will be different for you.

```shell
$ curl http://localhost:8080/services/test

{"code":"success","payload":[{"msg":"testing","id":"e88acedb-203d-4a7d-8269-1df6c1377512","createDate":"2020-10-11T23:21:16.378Z"}]}
```

## Install Mocha

Run the following command to install Mocha and add it to the developer dependencies:

```shell
$ npm install --save-dev mocha
```

## Refactor Dockerfile to run tests

Okay, now that we know our application is running properly, let’s try and run our tests inside of the container. We’ll use the same docker run command we used above but this time, we’ll override the CMD that is inside of our container with npm run test. This will invoke the command that is in the package.json file under the “script” section. See below.

```shell
{
...
  "scripts": {
    "test": "mocha ./**/*.js",
    "start": "nodemon server.js"
  },
...
}
```

Below is the Docker command to start the container and run tests:

```shell
$ docker run -it --rm --name app -p 8080:80 node-docker npm run test
> node-docker@0.1.0 test /code
> mocha ./**/*.js

sh: 1: mocha: not found
```

As you can see, we received an error. This error is telling us that the Mocha executable could not be found. Let’s take a look at the Dockerfile.

```dockerfile
FROM node:14.15.4

WORKDIR /code

COPY package.json package.json
COPY package-lock.json package-lock.json

RUN npm ci --production
COPY . .

CMD [ "node", "server.js" ]
```

The error is occurring because we are passing the `--production` flag to the npm ci command when it installs our dependencies. This tells npm to not install packages that are located under the "devDependencies" section in the package.json file. Therefore, Mocha will not be installed inside the image and will not be found when we try to run it.

Since we want to follow best practices and not include anything inside the container that we do not need to run our application we can’t just remove the `--production` flag. We have a couple of options to fix this. One option is to create a second Dockerfile that would only be used to run tests. This has a couple of problems. The primary being keeping two Dockerfiles up-to-date. The second option is to use multi-stage builds. We can create a stage for production and one for testing. This is our preferred solution.

### Multi-stage Dockerfile for testing

Below is a multi-stage Dockerfile tha we will use to build our production image and our test image. Add the highlighted lines to your Dockerfile.

```dockerfile
FROM node:14.15.4 as base

WORKDIR /code

COPY package.json package.json
COPY package-lock.json package-lock.json

FROM base as test
RUN npm ci
COPY . .
CMD [ "npm", "run", "test" ]

FROM base as prod
RUN npm ci --production
COPY . .
CMD [ "node", "server.js" ]
```

We first add a label to the `FROM node:14.15.4` statement. This allows us to refer to this build stage in other build stages. Next we add a new build stage labeled test. We will use this stage for running our tests.

Now let’s rebuild our image and run our tests. We will run the same docker build command as above but this time we will add the `--target test` flag so that we specifically run the test build stage.

```shell
$ docker build -t node-docker --target test .
[+] Building 0.7s (11/11) FINISHED
...
 => => writing image sha256:049b37303e3355f...9b8a954f
 => => naming to docker.io/library/node-docker
```

Now that our test image is built, we can run it in a container and see if our tests pass.

```shell
$ docker run -it --rm --name app -p 8080:80 node-docker
> node-docker@0.1.0 test /code
> mocha ./**/*.js

  Calculator
    adding
      ✓ should return 4 when adding 2 + 2
      ✓ should return 0 when adding zeros
    subtracting
      ✓ should return 4 when subtracting 4 from 8
      ✓ should return 0 when subtracting 8 from 8

  4 passing (11ms)
```

I’ve truncated the build output but you can see that the Mocha test runner completed and all our tests passed.

This is great but at the moment we have to run two docker commands to build and run our tests. We can improve this slightly by using a RUN statement instead of the CMD statement in the test stage. The CMD statement is not executed during the building of the image but is executed when you run the image in a container. While with the RUN statement, our tests will be run during the building of the image and stop the build when they fail.

Update your Dockerfile with the highlighted line below.

```dockerfile
FROM node:14.15.4 as base

WORKDIR /code

COPY package.json package.json
COPY package-lock.json package-lock.json

FROM base as test
RUN npm ci
COPY . .
RUN npm run test

FROM base as prod
RUN npm ci --production
COPY . .
CMD [ "node", "server.js" ]
```

Now to run our tests, we just need to run the docker build command as above.

```dockerfile
$ docker build -t node-docker --target test .
[+] Building 1.2s (13/13) FINISHED
...
 => CACHED [base 2/4] WORKDIR /code
 => CACHED [base 3/4] COPY package.json package.json
 => CACHED [base 4/4] COPY package-lock.json package-lock.json
 => CACHED [test 1/3] RUN npm ci
 => CACHED [test 2/3] COPY . .
 => CACHED [test 3/3] RUN npm run test
 => exporting to image
 => => exporting layers
 => => writing image sha256:bcedeeb7d9dd13d...18ca0a05034ed4dd4
 ```

I’ve truncated the output again for simplicity but you can see that our tests are run and passed. Let’s break one of the tests and observe the output when our tests fail.

Open the util/math.js file and change line 12 to the following.

```shell
11 function subtract( num1, num2 ) {
12   return num2 - num1
13 }
```

Now, run the same docker build command from above and observe that the build fails and the failing testing information is printed to the console.

```shell
$ docker build -t node-docker --target test .
 > [test 3/3] RUN npm run test:
#11 0.509
#11 0.509 > node-docker@0.1.0 test /code
#11 0.509 > mocha ./**/*.js
#11 0.509
#11 0.811
#11 0.813
#11 0.815   Calculator
#11 0.815     adding
#11 0.817       ✓ should return 4 when adding 2 + 2
#11 0.818       ✓ should return 0 when adding zeros
#11 0.818     subtracting
#11 0.822       1) should return 4 when subtracting 4 from 8
#11 0.823       ✓ should return 0 when subtracting 8 from 8
#11 0.823
#11 0.824
#11 0.824   3 passing (14ms)
#11 0.824   1 failing
#11 0.824
#11 0.827   1) Calculator
#11 0.827        subtracting
#11 0.827          should return 4 when subtracting 4 from 8:
#11 0.827
#11 0.827       AssertionError [ERR_ASSERTION]: Expected values to be strictly equal:
#11 0.827
#11 0.827 -4 !== 4
#11 0.827
#11 0.827       + expected - actual
#11 0.827
#11 0.827       --4
#11 0.827       +4
#11 0.827
#11 0.827       at Context.<anonymous> (util/math.test.js:18:14)
#11 0.827       at processImmediate (internal/timers.js:461:21)
...
executor failed running [/bin/sh -c npm run test]: exit code: 1
```

## Next steps

In this module, we took a look at creating a general development image that we can use pretty much like our normal command line. We also set up our Compose file to map our source code into the running container and exposed the debugging port.

In the next module, we’ll take a look at how to set up a CI/CD pipeline using GitHub Actions. See:

[Configure CI/CD](configure-ci-cd.md){: .button .outline-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs](https://github.com/docker/docker.github.io/issues/new?title=[Node.js%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.

<br />
