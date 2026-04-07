Volumes are persistent data stores implemented by the container engine. Compose offers a neutral way for services to mount volumes, and configuration parameters to allocate them to infrastructure. The top-level `volumes` declaration lets you configure named volumes that can be reused across multiple services.

If you want a named bind mount, use the `local` driver with `driver_opts`. This pattern gives a Compose volume a stable name while mapping it to a specific host path:

```yaml
volumes:
  app-data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /srv/app-data
```

The `type`, `o`, and `device` keys are passed through to the local driver. For a one-off host-path mount on a single service, see [bind mounts](/manuals/engine/storage/bind-mounts.md).
