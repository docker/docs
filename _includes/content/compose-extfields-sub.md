It is possible to re-use configuration fragments using extension fields. Those
special fields can be of any format as long as they are located at the root of
your Compose file and their name start with the `x-` character sequence.

```none
version: '2.1'
x-custom:
  items:
    - a
    - b
  options:
    max-size: '12m'
  name: "custom"
```

The contents of those fields are ignored by Compose, but they can be
inserted in your resource definitions using [YAML anchors](http://www.yaml.org/spec/1.2/spec.html#id2765878).
For example, if you want several of your services to use the same logging
configuration:

```none
logging:
  options:
    max-size: '12m'
    max-file: '5'
  driver: json-file
```

You may write your Compose file as follows:

```none
version: '2.1'
x-logging:
  &default-logging
  options:
    max-size: '12m'
    max-file: '5'
  driver: json-file

services:
  web:
    image: myapp/web:latest
    logging: *default-logging
  db:
    image: mysql:latest
    logging: *default-logging
```

It is also possible to partially override values in extension fields using
the [YAML merge type](http://yaml.org/type/merge.html). For example:

```none
version: '2.1'
x-volumes:
  &default-volume
  driver: foobar-storage

services:
  web:
    image: myapp/web:latest
    volumes: ["vol1", "vol2", "vol3"]
volumes:
  vol1: *default-volume
  vol2:
    << : *default-volume
    name: volume02
  vol3:
    << : *default-volume
    driver: default
    name: volume-local
```
