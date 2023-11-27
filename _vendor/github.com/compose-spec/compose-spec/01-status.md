# Status of this document

This document specifies the Compose file format used to define multi-containers applications. Distribution of this document is unlimited.

## Requirements and optional attributes

The Compose specification includes properties designed to target a local [OCI](https://opencontainers.org/) container runtime,
exposing Linux kernel specific configuration options, but also some Windows container specific properties. It is also designed for cloud platform features related to resource placement on a cluster, replicated application distribution, and scalability.

We acknowledge that no Compose implementation is expected to support all attributes, and that support for some properties
is platform dependent and can only be confirmed at runtime. The definition of a versioned schema to control the supported
properties in a Compose file, established by the [docker-compose](https://github.com/docker/compose) tool where the Compose
file format was designed, doesn't offer any guarantee to the end-user that attributes will be actually implemented.

The specification defines the expected configuration syntax and behavior. Unless noted, supporting any of these is optional.

A Compose implementation to parse a Compose file using unsupported attributes should warn users. We recommend the following implementors
to support those running modes:

* Default: warn the user about unsupported attributes, but ignore them
* Strict: warn the user about unsupported attributes and reject the Compose file
* Loose: ignore unsupported attributes AND unknown attributes (that were not defined by the spec by the time implementation was created)

From this point onwards, references made to 'Compose' can be interpreted as 'a Compose implementation'. 
