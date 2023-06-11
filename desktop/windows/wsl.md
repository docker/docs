---
description: Turn on the Docker WSL 2 backend and get to work using best practices, GPU support, and more in this thorough guide. 
keywords: wsl, wsl2, installing wsl2, wsl installation, docker wsl2, wsl docker, wsl2 tech preview, wsl install docker, install docker wsl, how to install docker in wsl
redirect_from:
- /docker-for-windows/wsl/
- /docker-for-windows/wsl-tech-preview/
title: Docker Desktop WSL 2 backend on Windows
---

Windows Subsystem for Linux (WSL) 2 is a full Linux kernel built by Microsoft, which allows Linux distributions to run without managing virtual machines. With Docker Desktop running on WSL 2, users can leverage Linux workspaces and avoid maintaining both Linux and Windows build scripts. In addition, WSL 2 provides improvements to file system sharing and boot time.

Docker Desktop uses the dynamic memory allocation feature in WSL 2 to improve the resource consumption. This means, Docker Desktop only uses the required amount of CPU and memory resources it needs, while enabling CPU and memory-intensive tasks such as building a container, to run much faster.

Additionally, with WSL 2, the time required to start a Docker daemon after a cold start is significantly faster. It takes less than 10 seconds to start the Docker daemon compared to almost a minute in the previous version of Docker Desktop.

## Prerequisites

Before you turn on the Docker Desktop WSL 2, ensure you have:

- WSL version 1.1.3.0 or above.
- Windows 10, version 21H2 or higher, or Windows 11, version 21H2 or higher. For more information, see [System requirements](https://docs.docker.com/desktop/install/windows-install/#system-requirements).
- Enabled WSL 2 feature on Windows. For detailed instructions, refer to the [Microsoft documentation](https://docs.microsoft.com/en-us/windows/wsl/install-win10){:target="_blank" rel="noopener" class="_"}.
- Downloaded and installed the [Linux kernel update package](https://docs.microsoft.com/windows/wsl/wsl2-kernel){:target="_blank" rel="noopener" class="_"}.

## Turn on Docker Desktop WSL 2

1. Download [Docker Desktop for Windows](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe).
2. Follow the usual installation instructions to install Docker Desktop. If you are running a supported system, Docker Desktop prompts you to enable WSL 2 during installation. Read the information displayed on the screen and enable WSL 2 to continue.
3. Start Docker Desktop from the **Windows Start** menu.
4. From the Docker menu, select **Settings** and then **General**.
5. Select the **Use WSL 2 based engine** check box.

    If you have installed Docker Desktop on a system that supports WSL 2, this option is enabled by default.
6. Select **Apply & Restart**.

Now `docker` commands work from Windows using the new WSL 2 engine.

## Enabling Docker support in WSL 2 distros

WSL 2 adds support for "Linux distros" to Windows, where each distro behaves like a VM except they all run on top of a single shared Linux kernel.

Docker Desktop does not require any particular Linux distros to be installed. The `docker` CLI and UI all work fine from Windows without any additional Linux distros. However for the best developer experience, we recommend installing at least one additional distro and enabling Docker support by:

1. Ensure the distribution runs in WSL 2 mode. WSL can run distributions in both v1 or v2 mode.

    To check the WSL mode, run:

     ```console
     $ wsl.exe -l -v
     ```

    To upgrade your existing Linux distro to v2, run:

    ```console
    $ wsl.exe --set-version (distro name) 2
    ```

    To set v2 as the default version for future installations, run:

    ```console
    $ wsl.exe --set-default-version 2
    ```

2. When Docker Desktop starts, go to **Settings** > **Resources** > **WSL Integration**.

    The Docker-WSL integration is enabled on your default WSL distribution. To change your default WSL distro, run `wsl --set-default <distro name>`

    For example, to set Ubuntu as your default WSL distro, run:
    
    ```console
    $ wsl --set-default ubuntu
    ```

    Optionally, select any additional distributions you would like to enable the Docker-WSL integration on.

3. Select **Apply & Restart**.

> **Note**
>
> Docker Desktop installs two special-purpose internal Linux distros `docker-desktop` and `docker-desktop-data`. The first (`docker-desktop`) is used to run the Docker engine (`dockerd`) while the second (`docker-desktop-data`) stores containers and images. Neither can be used for general development.

## Best practices

- To get the best out of the file system performance when bind-mounting files, we recommend storing source code and other data that is bind-mounted into Linux containers, for instance with `docker run -v <host-path>:<container-path>`, in the Linux file system, rather than the Windows file system. You can also refer to the [recommendation](https://docs.microsoft.com/en-us/windows/wsl/compare-versions){:target="_blank" rel="noopener" class="_"} from Microsoft.

  - Linux containers only receive file change events, "inotify events", if the
      original files are stored in the Linux filesystem. For example, some web development workflows rely on inotify events for automatic reloading when files have changed.
  - Performance is much higher when files are bind-mounted from the Linux
      filesystem, rather than remoted from the Windows host. Therefore avoid
      `docker run -v /mnt/c/users:/users`, where `/mnt/c` is mounted from Windows.
  - Instead, from a Linux shell use a command like `docker run -v ~/my-project:/sources <my-image>`
      where `~` is expanded by the Linux shell to `$HOME`.
- If you have concerns about the size of the docker-desktop-data VHDX, or need to change it, take a look at the [WSL tooling built into Windows](https://docs.microsoft.com/en-us/windows/wsl/vhd-size){:target="_blank" rel="noopener" class="_"}.
- If you have concerns about CPU or memory usage, you can configure limits on the memory, CPU, and swap size allocated to the [WSL 2 utility VM](https://docs.microsoft.com/en-us/windows/wsl/wsl-config#global-configuration-options-with-wslconfig){:target="_blank" rel="noopener" class="_"}.
- To avoid any potential conflicts with using WSL 2 on Docker Desktop, you must [uninstall any previous versions of Docker Engine](../../engine/install/ubuntu.md#uninstall-docker-engine) and CLI installed directly through Linux distributions before installing Docker Desktop.

## Develop with Docker and WSL 2

The following section describes how to start developing your applications using Docker and WSL 2. We recommend that you have your code in your default Linux distribution for the best development experience using Docker and WSL 2. After you have enabled WSL 2 on Docker Desktop, you can start working with your code inside the Linux distro and ideally with your IDE still in Windows. This workflow is straightforward if you are using [VSCode](https://code.visualstudio.com/download){:target="_blank" rel="noopener" class="_"}.

1. Open VS Code and install the [Remote - WSL](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-wsl){:target="_blank" rel="noopener" class="_"} extension. This extension allows you to work with a remote server in the Linux distro and your IDE client still on Windows.
2. Now, you can start working in VS Code remotely. To do this, open your terminal and type:

    ```console
    $ wsl
    ```

    ```console
    $ code .
    ```

    This opens a new VS Code connected remotely to your default Linux distro which you can check in the bottom corner of the screen.

    Alternatively, you can type the name of your default Linux distro in your Start menu, open it, and then run `code` .
3. When you are in VS Code, you can use the terminal in VS Code to pull your code and start working natively from your Windows machine.

## GPU support

> **Note**
>
> GPU support is only available in Docker Desktop for Windows with the WSL2 backend.

Starting with Docker Desktop 3.1.0, Docker Desktop supports WSL 2 GPU Paravirtualization (GPU-PV) on NVIDIA GPUs. To enable WSL 2 GPU Paravirtualization, you need:

- A machine with an NVIDIA GPU
- The latest Windows Insider version from the Dev Preview ring
- [Beta drivers](https://developer.nvidia.com/cuda/wsl){:target="_blank" rel="noopener" class="_"} from NVIDIA supporting WSL 2 GPU Paravirtualization
- Update WSL 2 Linux kernel to the latest version using `wsl --update` from an elevated command prompt
- Make sure the WSL 2 backend is enabled in Docker Desktop

To validate that everything works as expected, run the following command to run a short benchmark on your GPU:

```console
$ docker run --rm -it --gpus=all nvcr.io/nvidia/k8s/cuda-sample:nbody nbody -gpu -benchmark
```
The following displays:

```console
Run "nbody -benchmark [-numbodies=<numBodies>]" to measure performance.
        -fullscreen       (run n-body simulation in fullscreen mode)
        -fp64             (use double precision floating point values for simulation)
        -hostmem          (stores simulation data in host memory)
        -benchmark        (run benchmark to measure performance)
        -numbodies=<N>    (number of bodies (>= 1) to run in simulation)
        -device=<d>       (where d=0,1,2.... for the CUDA device to use)
        -numdevices=<i>   (where i=(number of CUDA devices > 0) to use for simulation)
        -compare          (compares simulation results running once on the default GPU and once on the CPU)
        -cpu              (run n-body simulation on the CPU)
        -tipsy=<file.bin> (load a tipsy model file for simulation)

> NOTE: The CUDA Samples are not meant for performance measurements. Results may vary when GPU Boost is enabled.

> Windowed mode
> Simulation data stored in video memory
> Single precision floating point simulation
> 1 Devices used for simulation
MapSMtoCores for SM 7.5 is undefined.  Default to use 64 Cores/SM
GPU Device 0: "GeForce RTX 2060 with Max-Q Design" with compute capability 7.5

> Compute 7.5 CUDA device: [GeForce RTX 2060 with Max-Q Design]
30720 bodies, total time for 10 iterations: 69.280 ms
= 136.219 billion interactions per second
= 2724.379 single-precision GFLOP/s at 20 flops per interaction
```

## Feedback

To provide feedback, create an issue in the [Docker Desktop for Windows GitHub](https://github.com/docker/for-win/issues){:target="_blank" rel="noopener" class="_"} repository and add the **WSL 2** label.
