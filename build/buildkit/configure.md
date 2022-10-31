---
title: Configure BuildKit
description: Learn how to configure BuildKit for your builder.
keywords: build, buildkit, configuration, buildx, network, cni, registry
---

If you [create a `docker-container` or `kubernetes` builder](../building/drivers/index.md)
with Buildx, you can set a custom [BuildKit configuration](toml-configuration.md)
by passing the [`--config` flag](../../engine/reference/commandline/buildx_create.md#config)
to the [`docker buildx create` command](../../engine/reference/commandline/buildx_create.md):

## Registry mirror

You can define a registry mirror to use for your builds:

```toml
# /etc/buildkitd.toml
debug = true
[registry."docker.io"]
  mirrors = ["mirror.gcr.io"]
```

> **Note**
>
> `debug = true` has been added to be able to debug requests
> in the BuildKit daemon and see if the mirror is effectively used.

Then [create a `docker-container` builder](../building/drivers/docker-container.md)
that will use this [BuildKit configuration](toml-configuration.md):

```console
$ docker buildx create --use --bootstrap \
  --name mybuilder \
  --driver docker-container \
  --config /etc/buildkitd.toml
```

Build an image:

```console
$ docker buildx build --load . -f-<<EOF
FROM alpine
RUN echo "hello world"
EOF
```

Now let's check the BuildKit logs in the builder container:

```console
$ docker logs buildx_buildkit_mybuilder0
```
```text
...
time="2022-02-06T17:47:48Z" level=debug msg="do request" request.header.accept="application/vnd.docker.container.image.v1+json, */*" request.header.user-agent=containerd/1.5.8+unknown request.method=GET spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="fetch response received" response.header.accept-ranges=bytes response.header.age=1356 response.header.alt-svc="h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000,h3-Q050=\":443\"; ma=2592000,h3-Q046=\":443\"; ma=2592000,h3-Q043=\":443\"; ma=2592000,quic=\":443\"; ma=2592000; v=\"46,43\"" response.header.cache-control="public, max-age=3600" response.header.content-length=1469 response.header.content-type=application/octet-stream response.header.date="Sun, 06 Feb 2022 17:25:17 GMT" response.header.etag="\"774380abda8f4eae9a149e5d5d3efc83\"" response.header.expires="Sun, 06 Feb 2022 18:25:17 GMT" response.header.last-modified="Wed, 24 Nov 2021 21:07:57 GMT" response.header.server=UploadServer response.header.x-goog-generation=1637788077652182 response.header.x-goog-hash="crc32c=V3DSrg==" response.header.x-goog-hash.1="md5=d0OAq9qPTq6aFJ5dXT78gw==" response.header.x-goog-metageneration=1 response.header.x-goog-storage-class=STANDARD response.header.x-goog-stored-content-encoding=identity response.header.x-goog-stored-content-length=1469 response.header.x-guploader-uploadid=ADPycduqQipVAXc3tzXmTzKQ2gTT6CV736B2J628smtD1iDytEyiYCgvvdD8zz9BT1J1sASUq9pW_ctUyC4B-v2jvhIxnZTlKg response.status="200 OK" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="fetch response received" response.header.accept-ranges=bytes response.header.age=760 response.header.alt-svc="h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000,h3-Q050=\":443\"; ma=2592000,h3-Q046=\":443\"; ma=2592000,h3-Q043=\":443\"; ma=2592000,quic=\":443\"; ma=2592000; v=\"46,43\"" response.header.cache-control="public, max-age=3600" response.header.content-length=1471 response.header.content-type=application/octet-stream response.header.date="Sun, 06 Feb 2022 17:35:13 GMT" response.header.etag="\"35d688bd15327daafcdb4d4395e616a8\"" response.header.expires="Sun, 06 Feb 2022 18:35:13 GMT" response.header.last-modified="Wed, 24 Nov 2021 21:07:12 GMT" response.header.server=UploadServer response.header.x-goog-generation=1637788032100793 response.header.x-goog-hash="crc32c=aWgRjA==" response.header.x-goog-hash.1="md5=NdaIvRUyfar8201DleYWqA==" response.header.x-goog-metageneration=1 response.header.x-goog-storage-class=STANDARD response.header.x-goog-stored-content-encoding=identity response.header.x-goog-stored-content-length=1471 response.header.x-guploader-uploadid=ADPycdtR-gJYwC7yHquIkJWFFG8FovDySvtmRnZBqlO3yVDanBXh_VqKYt400yhuf0XbQ3ZMB9IZV2vlcyHezn_Pu3a1SMMtiw response.status="200 OK" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="do request" request.header.accept="application/vnd.docker.image.rootfs.diff.tar.gzip, */*" request.header.user-agent=containerd/1.5.8+unknown request.method=GET spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="fetch response received" response.header.accept-ranges=bytes response.header.age=1356 response.header.alt-svc="h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000,h3-Q050=\":443\"; ma=2592000,h3-Q046=\":443\"; ma=2592000,h3-Q043=\":443\"; ma=2592000,quic=\":443\"; ma=2592000; v=\"46,43\"" response.header.cache-control="public, max-age=3600" response.header.content-length=2818413 response.header.content-type=application/octet-stream response.header.date="Sun, 06 Feb 2022 17:25:17 GMT" response.header.etag="\"1d55e7be5a77c4a908ad11bc33ebea1c\"" response.header.expires="Sun, 06 Feb 2022 18:25:17 GMT" response.header.last-modified="Wed, 24 Nov 2021 21:07:06 GMT" response.header.server=UploadServer response.header.x-goog-generation=1637788026431708 response.header.x-goog-hash="crc32c=ZojF+g==" response.header.x-goog-hash.1="md5=HVXnvlp3xKkIrRG8M+vqHA==" response.header.x-goog-metageneration=1 response.header.x-goog-storage-class=STANDARD response.header.x-goog-stored-content-encoding=identity response.header.x-goog-stored-content-length=2818413 response.header.x-guploader-uploadid=ADPycdsebqxiTBJqZ0bv9zBigjFxgQydD2ESZSkKchpE0ILlN9Ibko3C5r4fJTJ4UR9ddp-UBd-2v_4eRpZ8Yo2llW_j4k8WhQ response.status="200 OK" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
...
```

As you can see, requests come from the GCR registry mirror (`response.header.x-goog*`).

## Setting registry certificates

If you specify certificates for registries in the [BuildKit configuration](toml-configuration.md),
the files will be copied into the container under `/etc/buildkit/certs` and
configuration will be updated to reflect that.

Take the following configuration that will be used for pushing an image to
this registry using self-signed certificates:

```toml
# /etc/buildkitd.toml
debug = true
[registry."myregistry.com"]
  ca=["/etc/certs/myregistry.pem"]
  [[registry."myregistry.com".keypair]]
    key="/etc/certs/myregistry_key.pem"
    cert="/etc/certs/myregistry_cert.pem"
```

Here we have configured a self-signed certificate for `myregistry.com` registry.

Now [create a `docker-container` builder](../building/drivers/docker-container.md)
that will use this BuildKit configuration:

```console
$ docker buildx create --use --bootstrap \
  --name mybuilder \
  --driver docker-container \
  --config /etc/buildkitd.toml
```

Inspecting the builder container, you can see that buildkitd configuration
has changed:

```console
$ docker exec -it buildx_buildkit_mybuilder0 cat /etc/buildkit/buildkitd.toml
```
```toml
debug = true

[registry]

  [registry."myregistry.com"]
    ca = ["/etc/buildkit/certs/myregistry.com/myregistry.pem"]

    [[registry."myregistry.com".keypair]]
      cert = "/etc/buildkit/certs/myregistry.com/myregistry_cert.pem"
      key = "/etc/buildkit/certs/myregistry.com/myregistry_key.pem"
```

And certificates copied inside the container:

```console
$ docker exec -it buildx_buildkit_mybuilder0 ls /etc/buildkit/certs/myregistry.com/
myregistry.pem    myregistry_cert.pem   myregistry_key.pem
```

Now you should be able to push to the registry with this builder:

```console
$ docker buildx build --push --tag myregistry.com/myimage:latest .
```

## CNI networking

It can be useful to use a bridge network for your builder if for example you
encounter a network port contention during multiple builds. If you're using
the BuildKit image, CNI is not [(yet)](https://github.com/moby/buildkit/issues/28){:target="_blank" rel="noopener" class="_"}.
available in it.

But you can create your own BuildKit image with CNI support:

```dockerfile
ARG BUILDKIT_VERSION=v{{ site.buildkit_version }}
ARG CNI_VERSION=v1.0.1

FROM --platform=$BUILDPLATFORM alpine AS cni-plugins
RUN apk add --no-cache curl
ARG CNI_VERSION
ARG TARGETOS
ARG TARGETARCH
WORKDIR /opt/cni/bin
RUN curl -Ls https://github.com/containernetworking/plugins/releases/download/$CNI_VERSION/cni-plugins-$TARGETOS-$TARGETARCH-$CNI_VERSION.tgz | tar xzv

FROM moby/buildkit:${BUILDKIT_VERSION}
ARG BUILDKIT_VERSION
RUN apk add --no-cache iptables
COPY --from=cni-plugins /opt/cni/bin /opt/cni/bin
ADD https://raw.githubusercontent.com/moby/buildkit/${BUILDKIT_VERSION}/hack/fixtures/cni.json /etc/buildkit/cni.json
```

> **Note**
>
> Here we use the [CNI config for integration tests in BuildKit](https://github.com/moby/buildkit/blob/master//hack/fixtures/cni.json){:target="_blank" rel="noopener" class="_"},
> but feel free to use your own config.

Now build this image:

```console
$ docker buildx build --tag buildkit-cni:local --load .
```

Then [create a `docker-container` builder](../building/drivers/docker-container.md)
that will use this image:

```console
$ docker buildx create --use --bootstrap \
  --name mybuilder \
  --driver docker-container \
  --driver-opt "image=buildkit-cni:local" \
  --buildkitd-flags "--oci-worker-net=cni"
```
