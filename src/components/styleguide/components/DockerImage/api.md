
```
import {DockerImage} from 'common';

// All these should create the same component
<DockerImage imageName="redis"/>
<DockerImage imageName="redis:latest"/>
<DockerImage imageName="library/redis:latest"/>
<DockerImage imageName="tutum/redis:latest"/>

// Standalone (only icon)
<DockerImage imageName="tutum" standalone/>

// Size
<DockerImage imageName="tutum" standalone size="large"/>
```

---

### Props

| name       | propType                                     | default | required | description                                                 |
|------------|----------------------------------------------|---------|----------|-------------------------------------------------------------|
| imageName  | string                                       |    -    |   YES    | Any docker image format: namespace/name:tag, name:tag, name |
| size       | oneOf(SIZE.REGULAR, SIZE.LARGE, SIZE.XLARGE) | REGULAR |    NO    |                                                             |
| standalone | bool                                         | false   |    NO    | Display only the image icon and not image name              |

---
