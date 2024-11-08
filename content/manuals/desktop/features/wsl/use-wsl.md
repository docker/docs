---
title: Use WSL
description: How to develop with Docker and WSL 2 and understand GPU support for WSL
keywords: wsl, wsl 2, develop, docker desktop, windows
aliases:
- /desktop/wsl/use-wsl/
---

The following section describes how to start developing your applications using Docker and WSL 2. We recommend that you have your code in your default Linux distribution for the best development experience using Docker and WSL 2. After you have turned on the WSL 2 feature on Docker Desktop, you can start working with your code inside the Linux distribution and ideally with your IDE still in Windows. This workflow is straightforward if you are using [VS Code](https://code.visualstudio.com/download).

## Develop with Docker and WSL 2

1. Open VS Code and install the [Remote - WSL](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-wsl) extension. This extension lets you work with a remote server in the Linux distribution and your IDE client still on Windows.
2. Open your terminal and type:

    ```console
    $ wsl
    ```
3. Navigate to your project directory and then type:

    ```console
    $ code .
    ```

    This opens a new VS Code window connected remotely to your default Linux distribution which you can check in the bottom corner of the screen.


Alternatively, you can open your default Linux distribution from the **Start** menu, navigate to your project directory, and then run `code .`

