Your configuration options can contain environment variables. Compose uses the
variable values from the shell environment in which `docker-compose` is run.
For example, suppose the shell contains `EXTERNAL_PORT=8000` and you supply
this configuration:

    web:
      build: .
      ports:
        - "${EXTERNAL_PORT}:5000"

When you run `docker-compose up` with this configuration, Compose looks for
the `EXTERNAL_PORT` environment variable in the shell and substitutes its
value in. In this example, Compose resolves the port mapping to `"8000:5000"`
before creating the `web` container.

If an environment variable is not set, Compose substitutes with an empty
string. In the example above, if `EXTERNAL_PORT` is not set, the value for the
port mapping is `:5000` (which is of course an invalid port mapping, and will
result in an error when attempting to create the container).

You can set default values for environment variables using a
[`.env` file](env-file.md), which Compose will automatically look for. Values
set in the shell environment will override those set in the `.env` file.

    $ unset EXTERNAL_PORT
    $ echo "EXTERNAL_PORT=6000" > .env
    $ docker-compose up          # EXTERNAL_PORT will be 6000
    $ export EXTERNAL_PORT=7000
    $ docker-compose up          # EXTERNAL_PORT will be 7000

Both `$VARIABLE` and `${VARIABLE}` syntax are supported.
Additionally when using the [2.1 file format](compose-versioning.md#version-21), it
is possible to provide inline default values using typical shell syntax:

- `${VARIABLE:-default}` will evaluate to `default` if `VARIABLE` is unset or
  empty in the environment.
- `${VARIABLE-default}` will evaluate to `default` only if `VARIABLE` is unset
  in the environment.

Other extended shell-style features, such as `${VARIABLE/foo/bar}`, are not
supported.

You can use a `$$` (double-dollar sign) when your configuration needs a literal
dollar sign. This also prevents Compose from interpolating a value, so a `$$`
allows you to refer to environment variables that you don't want processed by
Compose.

    web:
      build: .
      command: "$$VAR_NOT_INTERPOLATED_BY_COMPOSE"

If you forget and use a single dollar sign (`$`), Compose interprets the value as an environment variable and will warn you:

  The VAR_NOT_INTERPOLATED_BY_COMPOSE is not set. Substituting an empty string.
