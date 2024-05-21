The `volumes` attribute define mount host paths or named volumes that are accessible by service containers. You can use `volumes` to define multiple types of mounts; `volume`, `bind`, `tmpfs`, or `npipe`. 

If the mount is a host path and is only used by a single service, it can be declared as part of the service
definition. To reuse a volume across multiple services, a named
volume must be declared in the [top-level `volumes` key](07-volumes.md).