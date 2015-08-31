Docker Toolbox
==================================

[![docker toolbox logo](https://cloud.githubusercontent.com/assets/251292/9585188/2f31d668-4fca-11e5-86c9-826d18cf45fd.png)](https://www.docker.com/toolbox)

The Docker Toolbox installs everything you need to get started with
Docker on Mac OS X and Windows, including:

|                        | Mac    | Windows     | Desktop Linux  (Help Wanted)   |
|------------------------|--------|-------------|--------------------------------|
| Docker Client / Engine | Client | Client      | Engine                         |
| Docker Machine         | Yes    | Yes         | Yes                            |
| Docker Compose         | Yes    | Coming Soon | Yes                            |
| Docker Kitematic       | Yes    | Yes         | Coming Soon                    |
| VirtualBox 5.0         | Yes    | Yes         | No                             |
| Delivery Format        | .pkg   | .exe        | script* (cURL)                 |


## Installation and documentation

Documentation for Mac [is available
here](https://docs.docker.com/mac/started/).

Documentation for Windows [is available here](https://docs.docker.com/windows/started/).

Toolbox is currently unavailable for Linux; To get started with Docker on Linux, please follow the Linux [Getting Started Guide](https://docs.docker.com/linux/started/).

## Building the Docker Toolbox

Toolbox installers are built using Docker, so you'll need a Docker host set up. For example, using [Docker Machine](https://github.com/docker/machine):

```
$ docker-machine create -d virtualbox toolbox
$ eval "$(docker-machine env toolbox)"
```

Then, to build the Toolbox for both platforms:

```
make
```

Build for a specific platform:

```
make osx
```

or

```
make windows
```

The resulting installers will be in the `dist` directory.

## Frequently Asked Questions

**Do I have to install VirtualBox?**

No, you can deselect VirtualBox during installation. It is bundled in case you want to have a working environment for free.
