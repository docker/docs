Docker Toolbox
==================================

[![docker toolbox logo](https://cloud.githubusercontent.com/assets/3325447/9101412/6754c6a6-3b9b-11e5-8cc4-215358caee06.png)](https://www.docker.com/toolbox)

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

Documentation for Windows [is available here](https://docs.docker.com/windows/started/). *Note:* Some Windows computers may not have VT-X enabled by default. It is required for VirtualBox. To enable VT-X, please see the guide [here.](http://www.howtogeek.com/213795/how-to-enable-intel-vt-x-in-your-computers-bios-or-uefi-firmware).

Toolbox is currently unavailable for Linux; To get started with Docker on Linux, please follow the Linux [Getting Started Guide](https://docs.docker.com/linux/started/).

## Frequently Asked Questions

**Do I have to install VirtualBox?**

No, you can deselect VirtualBox during installation. It is bundled in case you want to have a working environment for free.
