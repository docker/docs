A Dockerfile is a text document that contains the instructions to assemble a
Docker image. When we tell Docker to build our image by executing the `docker build`
command, Docker reads these instructions, executes them, and creates a Docker
image as a result.

Let’s walk through the process of creating a Dockerfile for our application. In
the root of your project, create a file named `Dockerfile` and open this file in
your text editor.

> **What to name your Dockerfile?**
>
> The default filename to use for a Dockerfile is `Dockerfile` (without a file-
> extension). Using the default name allows you to run the `docker build` command
> without having to specify additional command flags.
>
> Some projects may need distinct Dockerfiles for specific purposes. A common
> convention is to name these `Dockerfile.<something>` or `<something>.Dockerfile`.
> Such Dockerfiles can then be used through the `--file` (or `-f` shorthand)
> option on the `docker build` command. Refer to the
> ["Specify a Dockerfile" section](/engine/reference/commandline/build/#specify-a-dockerfile--f)
> in the `docker build` reference to learn about the `--file` option.
>
> We recommend using the default (`Dockerfile`) for your project's primary
> Dockerfile, which is what we'll use for most examples in this guide.

The first line to add to a Dockerfile is a [`# syntax` parser directive](/engine/reference/builder/#syntax).
While _optional_, this directive instructs the Docker builder what syntax to use
when parsing the Dockerfile, and allows older Docker versions with BuildKit enabled
to upgrade the parser before starting the build. [Parser directives](/engine/reference/builder/#parser-directives)
must appear before any other comment, whitespace, or Dockerfile instruction in
your Dockerfile, and should be the first line in Dockerfiles.

```dockerfile
# syntax=docker/dockerfile:1
```

We recommend using `docker/dockerfile:1`, which always points to the latest release
of the version 1 syntax. BuildKit automatically checks for updates of the syntax
before building, making sure you are using the most current version.