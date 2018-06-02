---
description: Docker's use of Open Source
keywords: docker, opensource
title: Open source components and licensing
notoc: true
---

Docker Desktop Editions are built using open source software.
For details on the licensing, choose
![whale menu](images/whale-x.png){: .inline} -->
**About Docker** from within the application, then click **Acknowledgements**.

Docker Desktop Editions distribute some components that are licensed under the
GNU General Public License. You can download the source for these components
[here](https://download.docker.com/opensource/License.tar.gz).

The sources for `qemu-img` can be obtained
[here](http://wiki.qemu-project.org/download/qemu-2.4.1.tar.bz2). The sources
for the `gettext` and `glib` libraries that `qemu-img` requires were obtained
from [Homebrew](https://brew.sh) and may be retrieved using `brew install
--build-from-source gettext glib`.
