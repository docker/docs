
```
import { Menu } from 'common';

const trigger = <Button icon variant="dull"><EllipsisIcon /></Button>;
const items = [
  {label: 'first', value: 'first_item'},
  {label: 'second', value: 'second_item'},
  {label: 'third', value: 'third_item'},
];


<Menu
  trigger={trigger}
  onSelect={ (item) => console.log(`selected: ${item}`) }
  items={items} />

```

### Props

| name     | propType                                   | default | required | description                                                                  |
|----------|--------------------------------------------|---------|----------|------------------------------------------------------------------------------|
| items    | arrayOf(shape({value, label, disabled}))   |   NO    |    YES   | Items to be displayed in the menu                                            |
| trigger  | node                                       |   NO    |    YES   | Element that will open the menu                                              |
| onSelect | func                                       |   NO    |    YES   | Called with item.value                                                       |
| offset   | array                                      |  [0,0]  |    NO    | X and Y values for tweaking the position of the menu relative to the trigger |

---
