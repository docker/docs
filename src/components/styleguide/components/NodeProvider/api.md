
```
import {NodeProvider} from 'common';

// All these should create the same component
<Nodeprovider providerName="digitalocean"/>
<Nodeprovider providerName="azure"/>
<Nodeprovider providerName="aws"/>
<Nodeprovider providerName="packet"/>

// Standalone (only icon)
<Nodeprovider providerName="digitalocean" standalone/>

// Size
<Nodeprovider providerName="digitalocean" standalone size="large"/>
```

---

### Props

| name       | propType                                     | default | required | description                                                 |
|------------|----------------------------------------------|---------|----------|-------------------------------------------------------------|
| providerName  | string                                       |    -    |   YES    | Any Docker Cloud compatible provider name |
| size       | oneOf(SIZE.REGULAR, SIZE.LARGE, SIZE.XLARGE) | REGULAR |    NO    |                                                             |
| standalone | bool                                         | false   |    NO    | Display only the image icon and not image name              |

---
