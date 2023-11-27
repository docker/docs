# Configs top-level element

Configs allow services to adapt their behaviour without the need to rebuild a Docker image.

Services can only access configs when explicitly granted by a [`configs`](05-services.md#configs) attribute within the `services` top-level element.

As with volumes, configs are mounted as files into a service's container's filesystem. The location of the mount point within the container defaults to `/<config-name>` in Linux containers and `C:\<config-name>` in Windows containers. 

By default, the config:
- Is owned by the user running the container command but can be overridden by service configuration.
- Has world-readable permissions (mode 0444), unless the service is configured to override this.

The top-level `configs` declaration defines or references configuration data that is granted to services in your Compose application. The source of the config is either `file` or `external`.

- `file`: The config is created with the contents of the file at the specified path.
- `environment`: The config content is created with the value of an environment variable.
- `content`: The content is created with the inlined value.
- `external`: If set to true, `external` specifies that this config has already been created. Compose does not
  attempt to create it, and if it does not exist, an error occurs.
- `name`: The name of the config object in the container engine to look up. This field can be used to
  reference configs that contain special characters. The name is used as is
  and will **not** be scoped with the project name.

## Example 1

`<project_name>_http_config` is created when the application is deployed,
by registering the content of the `httpd.conf` as the configuration data.

```yml
configs:
  http_config:
    file: ./httpd.conf
```

Alternatively, `http_config` can be declared as external. Compose looks up `http_config` to expose the configuration data to relevant services.

```yml
configs:
  http_config:
    external: true
```

## Example 2

`<project_name>_app_config` is created when the application is deployed,
by registering the inlined content as the configuration data. This comes with the
benefits Compose will infer variables when creating the config, which allows to
adjust content according to service configuration:

```yml
configs:
  app_config:
    content: |
      debug=${DEBUG}
      spring.application.admin.enabled=${DEBUG}
      spring.application.name=${COMPOSE_PROJECT_NAME}
```

## Example 3

External configs lookup can also use a distinct key by specifying a `name`. 

The following
example modifies the previous one to look up a config using the parameter `HTTP_CONFIG_KEY`. The
the actual lookup key will is set at deployment time by the [interpolation](12-interpolation.md) of
variables, but exposed to containers as hard-coded ID `http_config`.

```yml
configs:
  http_config:
    external: true
    name: "${HTTP_CONFIG_KEY}"
```

If `external` is set to `true`, all other attributes apart from `name` are irrelevant. If Compose detecs any other attribute, it rejects the Compose file as invalid.
