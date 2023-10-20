---
title: Run your Java tests
keywords: Java, build, test
description: How to build and run your Java tests
---

## Prerequisites

Work through the steps to build an image and run it as a containerized application in [Use your container for development](develop.md).

## Introduction

Testing is an essential part of modern software development. Testing can mean a lot of things to different development teams. There are unit tests, integration tests and end-to-end testing. In this guide you'll take a look at running your unit tests in Docker.

## Refactor Dockerfile to run tests

The **Spring Pet Clinic** source code has already tests defined in the test directory `src/test/java/org/springframework/samples/petclinic`. You can use the following Docker command to start the container and run tests:

```console
$ docker run -it --rm --name springboot-test java-docker ./mvnw test
...
[INFO] Results:
[INFO]
[WARNING] Tests run: 40, Failures: 0, Errors: 0, Skipped: 1
[INFO]
[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ------------------------------------------------------------------------
[INFO] Total time:  01:49 min
```

### Multi-stage Dockerfile for testing

Now, you'll learn how to pull the testing commands into your Dockerfile. The following is your updated multi-stage Dockerfile that you'll use to build your test image. Replace the contents of your Dockerfile with the following.

```dockerfile
# syntax=docker/dockerfile:1

FROM eclipse-temurin:17-jdk-jammy as base
WORKDIR /app
COPY .mvn/ .mvn
COPY mvnw pom.xml ./
RUN ./mvnw dependency:resolve
COPY src ./src

FROM base as test
CMD ["./mvnw", "test"]

FROM base as development
CMD ["./mvnw", "spring-boot:run", "-Dspring-boot.run.profiles=mysql", "-Dspring-boot.run.jvmArguments='-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:8000'"]

FROM base as build
RUN ./mvnw package


FROM eclipse-temurin:17-jre-jammy as production
EXPOSE 8080
COPY --from=build /app/target/spring-petclinic-*.jar /spring-petclinic.jar
CMD ["java", "-Djava.security.egd=file:/dev/./urandom", "-jar", "/spring-petclinic.jar"]
```

You added a new build stage labeled `test`. You'll use this stage for running your tests.

Now, rebuild your image and run your tests. You'll run the `docker build` command like you did previously, but this time you'll add the `--target test` flag so that you specifically run the test build stage.

```console
$ docker build -t java-docker --target test .
[+] Building 0.7s (6/6) FINISHED
...
 => => writing image sha256:967ac80cb7799a5d12a4bdfc67c37b5a6533c6e418c903907d3e86b7d4ebf89a
 => => naming to docker.io/library/java-docker
```

Now that your test image is built, you can run it as a container and see if your tests pass.

```console
$ docker run -it --rm --name springboot-test java-docker
[INFO] Scanning for projects...
[INFO]
[INFO] ------------< org.springframework.samples:spring-petclinic >------------
[INFO] Building petclinic 2.4.2
...

[INFO] Results:
[INFO]
[WARNING] Tests run: 40, Failures: 0, Errors: 0, Skipped: 1
[INFO]
[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ------------------------------------------------------------------------
[INFO] Total time:  01:22 min
```

The build output is truncated, but you can see that the Maven test runner was successful and all your tests passed.

This is great. However, you'll have to run two Docker commands to build and run your tests. You can improve this slightly by using a `RUN` statement instead of the `CMD` statement in the test stage. The `CMD` statement isn't executed during the building of the image, but is executed when you run the image in a container. When using the `RUN` statement, your tests run when building the image, and stop the build when they fail.

Update your Dockerfile with the following.

```dockerfile
# syntax=docker/dockerfile:1

FROM eclipse-temurin:17-jdk-jammy as base
WORKDIR /app
COPY .mvn/ .mvn
COPY mvnw pom.xml ./
RUN ./mvnw dependency:resolve
COPY src ./src

FROM base as test
RUN ["./mvnw", "test"]

FROM base as development
CMD ["./mvnw", "spring-boot:run", "-Dspring-boot.run.profiles=mysql", "-Dspring-boot.run.jvmArguments='-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:8000'"]

FROM base as build
RUN ./mvnw package

FROM eclipse-temurin:17-jre-jammy as production
EXPOSE 8080
COPY --from=build /app/target/spring-petclinic-*.jar /spring-petclinic.jar
CMD ["java", "-Djava.security.egd=file:/dev/./urandom", "-jar", "/spring-petclinic.jar"]
```

Now, to run your tests, you just need to run the `docker build` command.

```console
$ docker build -t java-docker --target test .
[+] Building 27.6s (11/12)
 => CACHED [base 3/6] COPY .mvn/ .mvn
 => CACHED [base 4/6] COPY mvnw pom.xml ./
 => CACHED [base 5/6] RUN ./mvnw dependency:resolve
 => CACHED [base 6/6] COPY src ./src
 => [test 1/1] RUN ["./mvnw", "test"]
 => exporting to image
 => => exporting layers
=> => writing image sha256:10cb585a7f289a04539e95d583ae97bcf8725959a6bd32c2f5632d0e7c1d16a0
=> => naming to docker.io/library/java-docker
```

The build output is truncated for simplicity, but you can see that your tests ran successfully and passed. Look at one of the tests and observe the output when your tests fail.

Open the `src/test/java/org/springframework/samples/petclinic/model/ValidatorTests.java` file and change the following assertion.

```diff
- assertThat(violation.getMessage()).isEqualTo("must not be empty");
+ assertThat(violation.getMessage()).isEqualTo("must be empty");
```

Now, run the `docker build` command and observe that the build fails and the failing testing information is printed to the console.

```console
$ docker build -t java-docker --target test .
 => [base 6/6] COPY src ./src
 => ERROR [test 1/1] RUN ["./mvnw", "test"]
...
------
executor failed running [./mvnw test]: exit code: 1
```

## Next steps

In this module, you took a look at running tests as part of your Docker image build process.

In the next module, you’ll take a look at how to set up a CI/CD pipeline using GitHub Actions.

{{< button text="Configure CI/CD" url="configure-ci-cd.md" >}}
