
```
import { Tab, Tabs } from 'common';

<Tabs selected={1} onSelect={this.handleSelect}>
  <Tab>Tab 1</Tab>
  <Tab>Tab 2</Tab>
  <Tab>Tab 3</Tab>
</Tabs>

<Tabs selected={this.selectedValue()} onSelect={this.handleSelect}>
  <Tab value="val1">Tab 1</Tab>
  <Tab value="val2">Tab 2</Tab>
  <Tab value="val3">Tab 3</Tab>
</Tabs>

<Tab>Im a Button</Tab>
<Tab selected>And I look like a Tab</Tab>

<Tabs icons>
  <Tab><DockerFlatIcon /></Tab>
  <Tab><PrivateIcon /></Tab>
  <Tab><PublicIcon /></Tab>
</Tabs>

```

### Tabs Props

| name     | propType                         | default | required | description                                                 |
|----------|----------------------------------|---------|----------|-------------------------------------------------------------|
| onSelect | func                             |   NO    |    YES   |                                                             |
| selected | oneOf([string, number])          |   0     |    NO    |                                                             |
| icons    | bool                             |  false  |    NO    | Use icons as tabs                                           |

### Tab Props

Extends `<Button>`.

| name     | propType                         | default | required | description                                                 |
|----------|----------------------------------|---------|----------|-------------------------------------------------------------|
| selected | bool                             |  false  |   NO     |                                                             |
| value    | oneOf([string, number])          |    -    |   NO     | Defaults to the current tab index when used withing `<Tabs>`|

---
