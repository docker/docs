---
description: Create a Docker Compose application using ASP.NET Core and SQL Server on Linux in Docker. 
keywords: docker, dockerize, dockerizing, dotnet, .NET, Core, article, example, platform, installation, containers, images, image, dockerfile, build, ASP.NET Core, SQL Server, mssql
title: "Quickstart: Compose and ASP.NET Core with SQL Server"
---

This quick-start guide demonstrates how to use Docker Compose to set up and run the sample ASP.NET Core application using the [ASP.NET Core Build image](https://hub.docker.com/r/microsoft/aspnetcore-build/) with the [SQL Server on Linux image](https://hub.docker.com/r/microsoft/mssql-server-linux/). You just need to have [Docker Engine](https://docs.docker.com/engine/installation/) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your platform of choice.

For this sample, we will use [Yeoman](http://yeoman.io/) to generate a sample `dotnet core web application`. After that, we will configure this app to use our `SQL Server database` and then create a `docker-compose.yml` that will define the behavior of all of these components.

# Pre-requisite: Install yeoman and the asp-net generator.

1. Install the [node package manager (npm)](http://blog.npmjs.org/post/85484771375/how-to-install-npm) on your platform of choice.

1. Install Yeoman using `npm`.

    ```bash
    npm install -g yo
    ```

1. Install the `aspnet generator` using `npm`.

    ```bash
    npm install -g generator-aspnet
    ```

# Create the Docker Compose application

1. Create an empty directory and navigate into it.

    This directory will be the context of your docker-compose project. It should only contain the resources necessary to build it.

1. Within your directory, use `Yeoman` to generate a sample web application:

    The parameters are `yo [projecttype [applicationname] [uiframework]]`
    
    In the following command, Yeoman will generate a `web` application called `webapp` that uses the `bootstrap` framework for the UI (This last parameter won't affect this tutorial)

    ```bash
    yo aspnet web webapp bootstrap
    ``` 

1. Navigate into the `webapp` directory that the previous command just generated. 

    Voilá! There's already a `Dockerfile` that was automatically generated for you. In the next few steps we will adapt this project to connect to the SQL Server database.

1. Create a `docker-compose.yml` file. Write the following in the file, and make sure to replace the password in the `SA_PASSWORD` environment variable under `db` below. This file will define the way the images will interact as micro-services. 

    >**Note**: SQL Server requires a secure password (Minimum length 8 characters, including uppercase and lowercase letters, base 10 digits and/or non-alphanumeric symbols).

    ```
    version: '2'

    services:
        web:
            build: .
            ports: 
                - "8000:80"
            depends_on:
                - db
        db:
            image: "microsoft/mssql-server-linux"
            ports: 
                - "1433:1433"
            environment:
                SA_PASSWORD: "your_password"
                ACCEPT_EULA: "Y"
    ```

    This file defines the `web` and `db` micro-services, their relationship, the ports they are using, and their specific environment variables.

1. Go to `Startup.cs` and locate the function called `ConfigureServices` (Hint: it should be under line 42). Replace the entire function to use the following code (watch out for the brackets!).

    >**Note**: Make sure to update the `Password` field in the `connection` variable below to the one you defined in the `docker-compose.yml` file.

    ```csharp
    [...]
    public void ConfigureServices(IServiceCollection services)
    {
        // Database connection string. 
        // Make sure to update the Password value below from "your_password" to your actual password.
        var connection = @"Server=db;Database=master;User=sa;Password=your_password;";
        
        // This line uses 'UseSqlServer' in the 'options' parameter
        // with the connection string defined above.
        services.AddDbContext<ApplicationDbContext>(
            options => options.UseSqlServer(connection));

        services.AddIdentity<ApplicationUser, IdentityRole>()
            .AddEntityFrameworkStores<ApplicationDbContext>()
            .AddDefaultTokenProviders();

        services.AddMvc();

        // Add application services.
        services.AddTransient<IEmailSender, AuthMessageSender>();
        services.AddTransient<ISmsSender, AuthMessageSender>();
    }
    [...]
    ```

1. Now replace the entire `Dockerfile` to the following content:

    ```
    FROM microsoft/aspnetcore-build:latest

    COPY . /app

    WORKDIR /app

    RUN ["dotnet", "restore"]

    RUN ["dotnet", "build"]

    EXPOSE 80/tcp

    RUN chmod +x ./entrypoint.sh

    CMD /bin/bash ./entrypoint.sh
    ```

    This file defines how to build the web app image. It will use the [microsoft/aspnetcore-build](https://hub.docker.com/r/microsoft/aspnetcore-build/) base image, copy the generated code, restore the dependencies, build the project and expose port 80. After that, it will call an `entrypoint script` that we will create in the next step. 

1. The previous `Dockerfile` makes use of an entrypoint to your webapp Docker image. Create this script in a file called `entrypoint.sh` and paste the contents below.

    >**Note**: Make sure to use UNIX line delimiters. The script won't work if you use Windows-based delimiters (Carriage return and line feed).

    ```bash
    #!/bin/bash

    set -e
    run_cmd="dotnet run --server.urls http://*:5000"

    until dotnet ef database update; do
    >&2 echo "SQL Server is starting up"
    sleep 1
    done

    >&2 echo "SQL Server is up - executing command"
    exec $run_cmd
    ```

    This script will restore the database after it starts up, and then will run the application. This allows some time for the SQL Server database image to start up.

1. Ready! You can now run the `docker-compose build` command.

    ```bash
    docker-compose build
    ```

1. Run the `docker-compose up` command. After a few seconds, you should be able to open [http://localhost:8000](http://localhost:8000) and see the landing ASP.NET core sample website. The application is listening on port 80 by default, but we mapped it to port 8000 in the `Dockerfile`.

    ```bash
    docker-compose up
    ```

    Go ahead and try out the website! This entire sample will use the SQL Server database image that is also running in Docker.

Ready! You now have a ASP.NET Core application running against SQL Server in Docker Compose!

## Further reading

- [Build your app using SQL Server](https://www.microsoft.com/en-us/sql-server/developer-get-started/?utm_medium=Referral&utm_source=docs.docker.com)
- [SQL Server on DockerHub](https://hub.docker.com/r/microsoft/mssql-server-linux/)
- [ASP.NET Core](https://www.asp.net/core)
- [ASP.NET Core Docker image](https://hub.docker.com/r/microsoft/aspnetcore/) on DockerHub