---
title: Parent image
id: parent image
short_description: >
 The image designated in the `FROM` directive in the image's Dockerfile
---

An image's parent image is the image designated in the `FROM` directive
in the image's Dockerfile. All subsequent commands are based on this parent
image. A Dockerfile with the `FROM scratch` directive uses no parent image, and creates
a base image.