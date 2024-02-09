---
title: Use WSL
description: How to develop with Docker and WSL 2 and understand GPU support for WSL
keywords: wsl, wsl 2, develop, docker desktop, windows
---

## Develop with Docker and WSL 2

The following section describes how to start developing your applications using Docker and WSL 2. We recommend that you have your code in your default Linux distribution for the best development experience using Docker and WSL 2. After you have turned on the WSL 2 feature on Docker Desktop, you can start working with your code inside the Linux distro and ideally with your IDE still in Windows. This workflow is straightforward if you are using [VS Code](https://code.visualstudio.com/download).

1. Open VS Code and install the [Remote - WSL](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-wsl) extension. This extension lets you to work with a remote server in the Linux distro and your IDE client still on Windows.
2. Start working in VS Code remotely. To do this, open your terminal and type:

    ```console
    $ wsl
    ```

    ```console
    $ code .
    ```

    This opens a new VS Code window connected remotely to your default Linux distro which you can check in the bottom corner of the screen.

    Alternatively, you can type the name of your default Linux distro in your Start menu, open it, and then run `code` .
3. When you are in VS Code, you can use the terminal in VS Code to pull your code and start working natively from your Windows machine.
