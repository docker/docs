
```
import { Uptime } from 'common';

// Date Object
<Uptime since={new Date()} />

// As String
<Uptime since={(new Date()).toString()} />

// As Timestamp
<Uptime since={+new Date()} />

```

---

### Props

| name     | propType                       | default | required | description                 |
|----------|--------------------------------|---------|----------|-----------------------------|
| since    | oneOfType([String,Number,Date])|         |   YES    | -                           |
| interval | number (milliseconds)          |  60000  |   NO     | How often should we update? |
| prefix   | String                         |         |   NO     | Prefix time with value      |

---
