+++
title = "images"
description = "description"
[menu.main]
identifier = "ucp_images"
parent = "ucp_ref"
+++

# images

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

Verifies all the required images used by UCP on the current engine. By default,
this command pulls any missing images. Use the `--pull` argument to change
behavior.

## Options

| Option                    | Description                                                                      |
|---------------------------|----------------------------------------------------------------------------------|
| `--debug`, `-D`           | Enable debug.                                                                    |
| `--jsonlog`               | Produce json formatted output for easier parsing.                                |
| `--interactive`, `-i`     | Enable interactive mode. You are prompted to enter all required information. |
| `--image-version "0.7.0"` | Select a specific UCP version.                                                   |
| `--pull "missing"`        | Specify image pull behavior (`always`, when `missing`, or `never`).              |
| `--list`                  | Don`t do anything, just list the images used by UCP                              |
