---
title: Run your Java tests
linkTitle: Run your tests
weight: 30
keywords: Java, build, test
description: How to build and run your Java tests
aliases:
  - /language/java/run-tests/
  - /guides/language/java/run-tests/
---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Java application](containerize.md).

## Overview

Testing is an essential part of modern software development. Testing can mean a lot of things to different development teams. There are unit tests, integration tests and end-to-end testing. In this guide you'll take a look at running your unit tests in Docker.

### Multi-stage Dockerfile for testing

In the following example, you'll pull the testing commands into your Dockerfile.
Replace the contents of your Dockerfile with the following.

```dockerfile {hl_lines="3-19"}
# syntax=docker/dockerfile:1

FROM eclipse-temurin:21-jdk-jammy as base
WORKDIR /build
COPY --chmod=0755 mvnw mvnw
COPY .mvn/ .mvn/

FROM base as test
WORKDIR /build
COPY ./src src/
RUN --mount=type=bind,source=pom.xml,target=pom.xml \
    --mount=type=cache,target=/root/.m2 \
    ./mvnw test

FROM base as deps
WORKDIR /build
RUN --mount=type=bind,source=pom.xml,target=pom.xml \
    --mount=type=cache,target=/root/.m2 \
    ./mvnw dependency:go-offline -DskipTests

FROM deps as package
WORKDIR /build
COPY ./src src/
RUN --mount=type=bind,source=pom.xml,target=pom.xml \
    --mount=type=cache,target=/root/.m2 \
    ./mvnw package -DskipTests && \
    mv target/$(./mvnw help:evaluate -Dexpression=project.artifactId -q -DforceStdout)-$(./mvnw help:evaluate -Dexpression=project.version -q -DforceStdout).jar target/app.jar

FROM package as extract
WORKDIR /build
RUN java -Djarmode=layertools -jar target/app.jar extract --destination target/extracted

FROM extract as development
WORKDIR /build
RUN cp -r /build/target/extracted/dependencies/. ./
RUN cp -r /build/target/extracted/spring-boot-loader/. ./
RUN cp -r /build/target/extracted/snapshot-dependencies/. ./
RUN cp -r /build/target/extracted/application/. ./
ENV JAVA_TOOL_OPTIONS="-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:8000"
CMD [ "java", "-Dspring.profiles.active=postgres", "org.springframework.boot.loader.launch.JarLauncher" ]

FROM eclipse-temurin:21-jre-jammy AS final
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser
COPY --from=extract build/target/extracted/dependencies/ ./
COPY --from=extract build/target/extracted/spring-boot-loader/ ./
COPY --from=extract build/target/extracted/snapshot-dependencies/ ./
COPY --from=extract build/target/extracted/application/ ./
EXPOSE 8080
ENTRYPOINT [ "java", "-Dspring.profiles.active=postgres", "org.springframework.boot.loader.launch.JarLauncher" ]
```

First, you added a new base stage. In the base stage, you added common instructions that both the test and deps stage will need.

Next, you added a new test stage labeled `test` based on the base stage. In this
stage you copied in the necessary source files and then specified `RUN` to run
`./mvnw test`. Instead of using `CMD`, you used `RUN` to run the tests. The
reason is that the `CMD` instruction runs when the container runs, and the `RUN`
instruction runs when the image is being built. When using `RUN`, the build will
fail if the tests fail.

Finally, you updated the deps stage to be based on the base stage and removed
the instructions that are now in the base stage.

Run the following command to build a new image using the test stage as the target and view the test results. Include `--progress=plain` to view the build output, `--no-cache` to ensure the tests always run, and `--target test` to target the test stage.

Now, build your image and run your tests. You'll run the `docker build` command and add the `--target test` flag so that you specifically run the test build stage.

```console
$ docker build -t java-docker-image-test --progress=plain --no-cache --target=test .
```

You should see output containing the following

```console
...

#15 101.3 [WARNING] Tests run: 45, Failures: 0, Errors: 0, Skipped: 2
#15 101.3 [INFO]
#15 101.3 [INFO] ------------------------------------------------------------------------
#15 101.3 [INFO] BUILD SUCCESS
#15 101.3 [INFO] ------------------------------------------------------------------------
#15 101.3 [INFO] Total time:  01:39 min
#15 101.3 [INFO] Finished at: 2024-02-01T23:24:48Z
#15 101.3 [INFO] ------------------------------------------------------------------------
#15 DONE 101.4s
```

## Next steps

In the next section, you’ll take a look at how to set up a CI/CD pipeline using
GitHub Actions.
