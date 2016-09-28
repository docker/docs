
```
import { List, ListItem } from 'common';

<List>
  <ListItem>Item1</ListItem>
  <ListItem>Item2</ListItem>
  <ListItem>Item3</ListItem>
</List>

<List selectable hover>
  <ListItem>Item1</ListItem>
  <ListItem selected>Item2</ListItem>
  <ListItem disabled>Item3</ListItem>

</List>
```

---

### List Props

| name       | propType          | default | required | description                                                |
|------------|-------------------|---------|----------|------------------------------------------------------------|
| selectable | bool              |  false  |    NO    | Whether the list items have a checkbox and can be selected |
| hover      | bool              |  false  |    NO    | Whether the list responds to hover events (only css)       |

---

### ListItem Props

| name       | propType              | default | required | description                                                  |
|------------|-----------------------|---------|----------|--------------------------------------------------------------|
| selectable | bool                  |  false  |    NO    | Whether the item itseld has a checkbox and can be selected   |
| hover      | bool                  |  false  |    NO    | Whether the list responds to hover events (only css)         |
| id         | string                |    -    |    NO    | Very much like **key**. This is the value passed to onSelect |
| selected   | bool                  |  false  |    NO    | Whether the row is currently selected                        |
| disabled   | bool                  |  false  |    NO    | Whether the row is currently disabled                        |
| onSelect   | func (ev,selected, id)|         |    NO    | Event handler for the toggling of the checkbox               |

---
