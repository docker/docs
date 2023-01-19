---
title: Invoke host binaries
description: Add invocations to host binaries from the frontend with the extension SDK.
keywords: Docker, extensions, sdk, build
---

In some cases, your extension needs to invoke some command from the host (the computer of your users). For example, you
might want to invoke the CLI of your cloud provider to create a new resource, or the CLI of a tool your extension
provides, or even a shell script that you want to run on the host. You could do that executing the CLI from a container with the extension SDK. But this CLI needs to access the host's
filesystem, which isn't easy nor fast if it runs in a container.
Host binaries allow exactly this: to invoke from the extension executables (as binaries, shell scripts)
shipped as part of your extension and deployed to the host. As extensions can run on multiple platforms, this
means that you need to ship the executables for all the platforms you want to support.

Learn more about extensions [architecture](../architecture/index.md).

> **Note**
>
> Only executables shipped as part of the extension can be invoked with the SDK. 

In this example, this CLI will be a simple `Hello world` script that must be invoked with a parameter and will return a 
string.

## Add the executables to the extension

Create a `bash` script for macOS and Linux, in the file `binaries/unix/hello.sh` with the following content:

```bash
#!/bin/sh
echo "Hello, $1!"
```

Create a `batch script` for Windows in another file `binaries/windows/hello.cmd` with the following content:

```bash
@echo off
echo "Hello, %1!"
```

Then update the `Dockerfile` to copy the `binaries` folder into the extension's container filesystem and make the
files executable.

```dockerfile
# Copy the binaries into the right folder
COPY --chmod=0755 binaries/windows/hello.cmd /windows/hello.cmd
COPY --chmod=0755 binaries/unix/hello.sh /linux/hello.sh
COPY --chmod=0755 binaries/unix/hello.sh /darwin/hello.sh
```

## Invoke the executable from the UI

In your extension, use the Docker Desktop Client object to [invoke the shell script](../dev/api/backend.md#invoke-an-extension-binary-on-the-host)
provided by the extension with the `ddClient.extension.host.cli.exec()` function.
In this example, the binary returns a string as result, obtained by `result?.stdout`, as soon as the extension view is rendered.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#react-app" data-group="react">For React</a></li>
  <li><a data-toggle="tab" data-target="#vue-app" data-group="vue">For Vue</a></li>
  <li><a data-toggle="tab" data-target="#angular-app" data-group="angular">For Angular</a></li>
  <li><a data-toggle="tab" data-target="#svelte-app" data-group="svelte">For Svelte</a></li>
</ul>

<div class="tab-content">
  <div id="react-app" class="tab-pane fade in active" markdown="1">

```typescript
export function App() {
  const ddClient = createDockerDesktopClient();
  const [hello, setHello] = useState("");

  useEffect(() => {
    const run = async () => {
      let binary = "hello.sh";
      if (ddClient.host.platform === 'win32') {
        binary = "hello.cmd";
      }

      const result = await ddClient.extension.host?.cli.exec(binary, ["world"]);
      setHello(result?.stdout);

    };
    run();
  }, [ddClient]);
    
  return (
    <div>
      {hello}
    </div>
  );
}
```


  </div>
  <div id="vue-app" class="tab-pane fade" markdown="1">

<br/>

> **Important**
> We don't have an example for Vue yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Vue){: target="_blank" rel="noopener" class="_"}
> and let us know you'd like a sample with Vue.
{: .important }

  </div>
  <div id="angular-app" class="tab-pane fade" markdown="1">

<br/>

> **Important**
> We don't have an example for Angular yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Angular){: target="_blank" rel="noopener" class="_"}
> and let us know you'd like a sample with Angular.
{: .important }

  </div>
  <div id="svelte-app" class="tab-pane fade" markdown="1">

<br/>

> **Important**
> We don't have an example for Svelte yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Svelte){: target="_blank" rel="noopener" class="_"}
> and let us know you'd like a sample with Svelte.
{: .important }

  </div>
</div>

## Configure the metadata file

The host binaries must be specified in the `metadata.json` so that Docker Desktop copies them on the host when installing
the extension. Once the extension is uninstalled, the binaries that were copied will be removed as well.

```json
{
  "vm": {
    ...
  },
  "ui": {
    ...
  },
  "host": {
    "binaries": [
      {
        "darwin": [
          {
            "path": "/darwin/hello.sh"
          }
        ],
        "linux": [
          {
            "path": "/linux/hello.sh"
          }
        ],
        "windows": [
          {
            "path": "/windows/hello.cmd"
          }
        ]
      }
    ]
  }
}
```

The `path` must reference the path of the binary inside the container.
