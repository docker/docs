# Profiles

With profiles you can define a set of active profiles so your Compose application model is adjusted for various usages and environments.
The exact mechanism is implementation specific and may include command line flags, environment variables, etc.

The [services](05-services.md) top-level element supports a `profiles` attribute to define a list of named profiles. 
Services without a `profiles` attribute are always enabled. 

A service is ignored by Compose when none of the listed `profiles` match the active ones, unless the service is
explicitly targeted by a command. In that case its profile is added to the set of active profiles.

> **Note**
>
> All other top-level elements are not affected by `profiles` and are always active.

References to other services (by `links`, `extends` or shared resource syntax `service:xxx`) do not
automatically enable a component that would otherwise have been ignored by active profiles. Instead
Compose returns an error.

## Illustrative example

```yaml
services:
  foo:
    image: foo
  bar:
    image: bar
    profiles:
      - test
  baz:
    image: baz
    depends_on:
      - bar
    profiles:
      - test
  zot:
    image: zot
    depends_on:
      - bar
    profiles:
      - debug
```

In the above example:

- If the Compose application model is parsed with no profile enabled, it only contains the `foo` service.
- If the profile `test` is enabled, the model contains the services `bar` and `baz`, and service `foo`, which is always enabled.
- If the profile `debug` is enabled, the model contains both `foo` and `zot` services, but not `bar` and `baz`,
  and as such the model is invalid regarding the `depends_on` constraint of `zot`.
- If the profiles `debug` and `test` are enabled, the model contains all services; `foo`, `bar`, `baz` and `zot`.
- If Compose is executed with `bar` as the explicit service to run, `bar` and the `test` profile
  are active even if `test` profile is not enabled.
- If Compose is executed with `baz` as the explicit service to run, the service `baz` and the
  profile `test` are active and `bar` is pulled in by the `depends_on` constraint.
- If Compose is executed with `zot` as the explicit service to run, again the model is
  invalid regarding the `depends_on` constraint of `zot`, since `zot` and `bar` have no common `profiles`
  listed.
- If Compose is executed with `zot` as the explicit service to run and profile `test` is enabled,
  profile `debug` is automatically enabled and service `bar` is pulled in as a dependency starting both
  services `zot` and `bar`.

See how you can use `profiles` in [Docker Compose](https://docs.docker.com/compose/profiles/).
