---
title: Invoke host binaries
description: Add invocations to host binaries from the frontend with the extension SDK.
keywords: Docker, extensions, sdk, build
---

In some cases, your extension needs to invoke some command from the host (the computer of your users). For example, you 
might wand to invoke the CLI of your cloud provider to create a new resource, or the CLI of a tool your extension 
provides. 
You could do that executing the CLI from a container with the extension SDK. But this CLI needs to access the host's 
filesystem, which isn't possible if it runs in a container.
Host binaries allows exactly this: to invoke binaries from the host. As extensions can run on multiple platforms, 
this means that you need to ship the binaries for all the platforms you want to support.

In this example, this CLI will be a simple `Hello world` script that must be invoked with a parameter and will return a 
string. Since the extension will support multiple platforms, the script will be written in `bash` for macOS and Linux, 
and in `powershell` for Windows.

## Add the binaries to the extension
Create a file `binaries/unix/hello.sh` with the following content:

```bash
#!/bin/sh
echo "Hello, $1!"
```

Create another file `binaries/windows/hello.ps1` with the following content:

```powershell
Write-Output "Hello, $args[0]!"
```

Make them both executable:

```bash
chmod +x binaries/*
```

Then update the `Dockerfile` to copy the `binaries` folder into the extension container.

```dockerfile
# Copy the binaries into the right folder
COPY binaries/windows/hello.ps1 /windows/hello.ps1
COPY binaries/unix/hello.sh /linux/hello.sh
COPY binaries/unix/hello.sh /darwin/hello.sh
```

## Invoke the host binary from the UI

In your app, use the Docker Desktop Client object and then invoke the binary provided by the extension with the 
`ddClient.extension.host.cli.exec()` function.
In this example, the binary returns a string as result, obtained by `result?.stdout`, as soon as the app starts.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#react-app" data-group="react">For React</a></li>
  <li><a data-toggle="tab" data-target="#vue-app" data-group="vue">For Vue</a></li>
  <li><a data-toggle="tab" data-target="#angular-app" data-group="angular">For Angular</a></li>
  <li><a data-toggle="tab" data-target="#svelte-app" data-group="svelte">For Svelte</a></li>
</ul>

<div class="tab-content">
  <div id="react-dockerfile" class="tab-pane fade in active" markdown="1">

```typescript

export function App() {
  const ddClient = createDockerDesktopClient();
  const [hello, setHello] = useState("");

  useEffect(() => {
    const run = async () => {
      let binary = "hello.sh";
      if (ddClient.host.platform === 'win32') {
        binary = "hello.ps1";
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
> We don't have an example for Vue yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Vue)
> and let us know you'd like a sample with Vue.
{: .important }

  </div>
  <div id="angular-app" class="tab-pane fade" markdown="1">

<br/>

> **Important**
> We don't have an example for Angular yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Angular)
> and let us know you'd like a sample with Angular.
{: .important }

  </div>
  <div id="svelte-app" class="tab-pane fade" markdown="1">

<br/>

> **Important**
> We don't have an example for Svelte yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Svelte)
> and let us know you'd like a sample with Svelte.
{: .important }

  </div>
</div>

## Configure the metadata file

The host binaries must be specified in the `metadata.json` so that Docker Desktop copies them on the host when installing the extension.

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
            "path": "/windows/hello.ps1"
          }
        ]
      }
    ]
  }
}
```

The `path` must reference the path of the binary inside the container.
