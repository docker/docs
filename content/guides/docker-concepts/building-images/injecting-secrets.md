---
title: Injecting secrets during the build
keywords: concepts, build, images, docker desktop
description: Injecting secrets during the build
---
<iframe width="650" height="365" src="https://www.youtube.com/embed/nsWWQ1xoEy0?rel=0" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

## Explanation

In this concept, you will learn the following:
- What are Build Secrets?
- How to inject Secrets during the build 

Secrets are sensitive pieces of information like passwords, API keys, or database credentials that your application needs during the build process, hence you shouldn't store them within the final image. Injecting these secrets securely ensures both functionality and security for your app


### Why not store secrets in the image?

Build secrets like passwords shouldn't be in your Docker image! While build arguments and environment variables are easy, they expose secrets. Instead, use secret mounts or SSH mounts for secure access during the build, ensuring your final image stays safe and sound. Remember, security first!

- **Security Risk**: Embedding secrets directly in the image exposes them to anyone who has access to the image, posing a significant security threat.
- **Lack of Flexibility**: Modifying secrets embedded in the image requires rebuilding the entire image, which can be inefficient.

### Secret mounts and SSH mounts

Instead of embedding secrets, use secret mounts or SSH mounts to provide secure access during the build process. These mounts let you to temporarily access secrets stored outside the build context, ensuring they're not stored inside the final image.



### Try it out

In this hands-on, you will learn how to inject sensitive information like passwords and API keys (secrets) into your Docker image build process without embedding them within the final image. This ensures both functionality and security for your application.

### Create a secret file

Create a file named `mysecretfile` containing your secret (replace bs-dfhfkdsfh with your actual secret):


```console
 cat mysecretfile
 bs-dfhfkdsfh
```

### Create a Dockerfile

Create a plain-text file named `Dockerfile` with the following content:

```diff
FROM node:20-alpine
# Mount the secret during the build
RUN --mount=type=secret,id=mysecret cat /run/secrets/mysecret
```

### Build an image

Use the `docker build` command with the `--secret` flag to map an environment variable (mysecret) to the secret ID (mysecret). 
Additionally, specify the source path of the secret file using the `-t` flag to tag the image:

```console
  docker build --secret id=mysecret,src=$PWD/secretfile -t testsecret .
```

### View the image history

Use the `docker history` command to see the build layers of your image. Notice that the secret file content isn't present in the final image layer:

```console
 docker history testsecret:latest
IMAGE          CREATED          CREATED BY                                      SIZE      COMMENT
b0c498a6f907   31 seconds ago   RUN /bin/sh -c cat /run/secrets/mysecret # b…   4.1kB     buildkit.dockerfile.v0
<missing>      11 days ago      /bin/sh -c #(nop)  CMD ["node"]                 0B
<missing>      11 days ago      /bin/sh -c #(nop)  ENTRYPOINT ["docker-entry…   0B
<missing>      11 days ago      /bin/sh -c #(nop) COPY file:4d192565a7220e13…   20.5kB
<missing>      11 days ago      /bin/sh -c apk add --no-cache --virtual .bui…   7.92MB
<missing>      11 days ago      /bin/sh -c #(nop)  ENV YARN_VERSION=1.22.19     0B
<missing>      11 days ago      /bin/sh -c addgroup -g 1000 node     && addu…   129MB
<missing>      11 days ago      /bin/sh -c #(nop)  ENV NODE_VERSION=20.11.1     0B
<missing>      4 weeks ago      /bin/sh -c #(nop)  CMD ["/bin/sh"]              0B
<missing>      4 weeks ago      /bin/sh -c #(nop) ADD file:d0764a717d1e9d0af…   8.42MB
```

- The `RUN` instruction with the `--mount` flag tells Docker to mount the secret during the build process.
- `type=secret,id=mysecret` defines the mount type as secret and assigns the ID mysecret.
- `cat /run/secrets/mysecret` reads the content of the secret from the mounted location.

This approach utilizes BuildKit to securely inject secrets during the build, ensuring your final image remains free of sensitive information. Remember, security is paramount! Always prioritize secure methods for handling secrets.

Note: This guide demonstrates a basic example using a local secret file. In production environments, consider using a secure secret management system like Vault.


## Additional resources

- [Build Secrets](https://docs.docker.com/build/building/secrets/)

{{< button text="Troubleshooting failing builds" url="troubleshooting-failing-builds" >}}
