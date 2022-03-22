<!-- This section is included in topics that contain instructions on how to configure registry.json file to enforce users to sign into Docker Desktop-->

## Create a registry.json file

After you’ve successfully installed Docker Desktop, create a `registry.json`
file. Before you create a `registry.json` file, ensure that the developer is a
member of at least one organization in Docker Hub. If the `registry.json` file
matches at least one organization the developer is a member of, they can sign
into Docker Desktop, and then access all their organizations.

### Windows

On Windows, you must create a file at
`C:\ProgramData\DockerDesktop\registry.json` with file permissions that ensure
that the developer using Docker Desktop cannot remove or edit the file (that is,
only the system administrator can write to the file). The file must be of type
`JSON` and contain the name of the organization in the `allowedOrgs` key.

To create your `registry.json` file on Windows:

1. Open Windows PowerShell and select Run as Administrator.
2. Type the following command `cd /ProgramData/DockerDesktop/`
3. Type `notepad registry.json` and enter the name of the Docker Hub
   organization that the developer belongs to in the `allowedOrgs` key and click
   **Save**. For example:

    ```json
    {
        "allowedOrgs": ["myorg"]
    }
    ```

### Mac

On macOS, you must create a file at `/Library/Application Support/com.docker.docker/registry.json` with file permissions that ensure that
the developer using Docker Desktop cannot remove or edit the file (that is, only
the system administrator can write to the file). The file must be of type `JSON`
and contain the name of the Docker Hub organization names in the `allowedOrgs`
key.

To create your `registry.json` file on macOS:

1. Navigate to VS Code or any text editor of your choice.
2. Enter the name of the Docker Hub organization that the developer belongs to in the  `allowedOrgs` key and save it in your Documents. For example:

    ```json
    {
        "allowedOrgs": ["myorg"]
    }
    ```

3. Open a new terminal and type the following command:

    ```console
    sudo mkdir -p /Library/Application\ Support/com.docker.docker
    ```

    If prompted, type your password associated with your local computer.

4. Type the following command:

     ```console
    sudo cp Documents/registry.json /Library/Application\ Support/com.docker.docker/registry.json
    ```
