---
title: MCP
description: Learn how to use MCP servers with Gordon
keywords: ai, mcp, gordon
---

## What is MCP?

Anthropic recently announced the [Model Context Protocol](https://www.anthropic.com/news/model-context-protocol) (MCP) specification, an open protocol that standardises how applications provide context to large language models. MCP functions as a client-server protocol, where the client (e.g., an application like Gordon) sends requests, and the server processes those requests to deliver the necessary context to the AI.

Gordon, along with other MCP clients like Claude Desktop, can interact with MCP servers running as containers. Docker has partnered with Anthropic to build container images for the [reference implementations](https://github.com/modelcontextprotocol/servers/) of MCP servers, available on Docker Hub under [the mcp namespace](https://hub.docker.com/u/mcp).

## Simple MCP server usage with Gordon

When you run the `docker ai` command in your terminal to ask a question, Gordon looks for a `gordon-mcp.yml` file in your working directory for a list of MCP servers that should be used when in that context. The `gordon-mcp.yml` file is a Docker Compose file that configures MCP servers as Compose services for Gordon to access.

The following minimal example shows how you can use the [mcp-time server](https://hub.docker.com/r/mcp/time) to provide temporal capabilities to Gordon. For more information, you can check out the [source code and documentation](https://github.com/modelcontextprotocol/servers/tree/main/src/time).

1. Create the `gordon-mcp.yml` file and add the time server:
    
    ```yaml
    services:
      time:
        image: mcp/time
    ```
    
2. With this file you can now ask Gordon to tell you the time in another timezone:
    
    ```bash
    $ docker ai 'what time is it now in kiribati?'
    
        • Calling get_current_time
    
      The current time in Kiribati (Tarawa) is 9:38 PM on January 7, 2025.
    
    ```
    

As you can see, Gordon found the MCP time server and called its tool when needed.

## Advanced usage

Some MCP servers need access to your filesystem or system environment variables. Docker Compose can help with this. Since `gordon-mcp.yml` is a Compose file you can add bind mounts using the regular Docker Compose syntax, which makes your filesystem resources available to the container:

```yaml
services:
  fs:
    image: mcp/filesystem
    command:
      - /rootfs
    volumes:
      - .:/rootfs
```

The `gordon-mcp.yml` file adds filesystem access capabilities to Gordon and since everything runs inside a container Gordon only has access to the directories you specify.

Gordon can handle any number of MCP servers. For example, if you give Gordon access to the internet with the `mcp/fetch` server:

```yaml
services:
  fetch:
    image: mcp/fetch
  fs:
    image: mcp/filesystem
    command:
      - /rootfs
    volumes:
      - .:/rootfs
```

You can now ask things like:

```bash
$ docker ai can you fetch rumpl.dev and write the summary to a file test.txt 

    • Calling fetch ✔️
    • Calling write_file ✔️
  
  The summary of the website rumpl.dev has been successfully written to the file test.txt in the allowed directory. Let me know if you need further assistance!
  
  
$ cat test.txt 
The website rumpl.dev features a variety of blog posts and articles authored by the site owner. Here's a summary of the content:

1. **Wasmio 2023 (March 25, 2023)**: A recap of the WasmIO 2023 conference held in Barcelona. The author shares their experience as a speaker and praises the organizers for a successful event.

2. **Writing a Window Manager in Rust - Part 2 (January 3, 2023)**: The second part of a series on creating a window manager in Rust. This installment focuses on enhancing the functionality to manage windows effectively.

3. **2022 in Review (December 29, 2022)**: A personal and professional recap of the year 2022. The author reflects on the highs and lows of the year, emphasizing professional achievements.

4. **Writing a Window Manager in Rust - Part 1 (December 28, 2022)**: The first part of the series on building a window manager in Rust. The author discusses setting up a Linux machine and the challenges of working with X11 and Rust.

5. **Add docker/docker to your dependencies (May 10, 2020)**: A guide for Go developers on how to use the Docker client library in their projects. The post includes a code snippet demonstrating the integration.

6. **First (October 11, 2019)**: The inaugural post on the blog, featuring a simple "Hello World" program in Go.%   

```

## What’s next?

Now that you’ve learned how to use MCP servers with Gordon, here are a few ways you can get started:

- Experiment: Try integrating one or more of the tested MCP servers into your `gordon-mcp.yml` file and explore their capabilities.
1. Explore the ecosystem: Check out the [reference implementations on GitHub](https://github.com/modelcontextprotocol/servers/) or browse the [Docker Hub MCP namespace](https://hub.docker.com/u/mcp) for additional servers that might suit your needs.
2. Build your own: If none of the existing servers meet your needs, or you’re curious about exploring how they work in more detail, consider developing a custom MCP server. Use the [MCP specification](https://www.anthropic.com/news/model-context-protocol) as a guide.
3. Share your feedback: If you discover new servers that work well with Gordon or encounter issues with existing ones, [share your findings to help improve the ecosystem.](https://docker.qualtrics.com/jfe/form/SV_9tT3kdgXfAa6cWa)

With MCP support, Gordon offers powerful extensibility and flexibility to meet your specific use cases whether you’re adding temporal awareness, file management, or internet access.

### List of known working MCP Servers

These are the MCP servers that have been tested successfully with Gordon:

- `mcp/time`
- `mcp/fetch`
- `mcp/filesystem`
- `mcp/postgres`
- `mcp/git`
- `mcp/sqlite`
- `mcp/github`

### List of untested MCP servers

These are the MCP servers that were not tested but should work if given the appropriate API tokens:

- `mcp/brave-search`
- `mcp/gdrive`
- `mcp/slack`
- `mcp/google-maps`
- `mcp/gitlab`
- `mcp/everything`
- `mcp/aws-kb-retrieval-server`
- `mcp/sentry`

### List of MCP servers that don’t work with Gordon

These are the MCP servers that are currently unsupported:

- `mcp/sequentialthinking` - The tool description is too long
- `mcp/puppeteer` - Puppeteer sends back images and Gordon doesn’t know how to handle them, it only handles text responses from tools
- `mcp/everart` - Everart sends back images and Gordon doesn’t know how to handle them, it only handles text responses from tools
- `mcp/memory` - There is no way to configure the server to use a custom path for its knowledge base