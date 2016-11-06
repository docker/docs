+++
title = "images"
description = "Verify the UCP images on this Docker Engine."
keywords= ["images, ucp, images"]
[menu.main]
parent = "ucp_ref"
identifier = "ucp_ref_images"
+++

# docker/ucp images

Verify the UCP images on this engine.

## Usage

```
docker run --rm -it \
           --name ucp \
           -v /var/run/docker.sock:/var/run/docker.sock \
           docker/ucp \
           images [command options]
```

## Description

This command will verify all the required images used by UCP on the current engine.
By default, this will pull any missing images. Use the `--pull` argument to change
behavior.

## Options

```nohighlight
--debug, -D                Enable debug mode
--jsonlog                  Produce json formatted output for easier parsing
--pull value               Specify image pull behavior ('always', when 'missing', or 'never') (default: "missing")
--registry-username value  Specify the username to pull required images with [$REGISTRY_USERNAME]
--registry-password value  Specify the password to pull required images with [$REGISTRY_PASSWORD]
--list                     Don't do anything, just list the images used by UCP
```
