
### Diagnosing from the terminal

On occasions it is useful to run the diagnostics yourself, for instance if
Docker Desktop for Windows cannot start.

First locate the `com.docker.diagnose`, that should be in `C:\Program
Files\Docker\Docker\resources\com.docker.diagnose.exe`.

To create *and upload*  diagnostics in Powershell, run:

```powershell
  PS C:\> & "C:\Program Files\Docker\Docker\resources\com.docker.diagnose.exe" gather -upload
```

After the diagnostics have finished, you should have the following output,
containing your diagnostic ID:

```sh
Diagnostics Bundle: C:\Users\User\AppData\Local\Temp\CD6CF862-9CBD-4007-9C2F-5FBE0572BBC2\20180720152545.zip
Diagnostics ID:     CD6CF862-9CBD-4007-9C2F-5FBE0572BBC2/20180720152545 (uploaded)
```

If you have a paid Docker subscription, open the [Docker Desktop support](https://hub.docker.com/support/desktop/){:target="_blank" rel="noopener" class="_"} form. Fill in the information required and add the ID to the Diagnostics ID field. Click **Submit** to request Docker Desktop support.

### Self-diagnose tool

Docker Desktop contains a self-diagnose tool which helps you to identify some common
problems. Before you run the self-diagnose tool, locate `com.docker.diagnose.exe`. This is usually installed in `C:\Program
Files\Docker\Docker\resources\com.docker.diagnose.exe`.

To run the self-diagnose tool in Powershell:

```powershell
PS C:\> & "C:\Program Files\Docker\Docker\resources\com.docker.diagnose.exe" check
```

The tool runs a suite of checks and displays **PASS** or **FAIL** next to each check. If there are any failures, it highlights the most relevant at the end.

> **Feedback**
>
> Let us know your feedback on the self-diagnose tool by creating an issue in the [for-win](https://github.com/docker/for-win/issues) GitHub repository.



