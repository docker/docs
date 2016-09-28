
```
import { MenuButton } from 'common';

const items = [
  {label: 'first', value: 'first_item'},
  {label: 'second', value: 'second_item'},
  {label: 'third', value: 'third_item'},
];


<MenuButton
  items={items}
  onSelect={(item) => console.log(`selected: ${item}`) }
  onClick={() => console.log('MenuButton Clicked')}
>Menu Button</MenuButton>

```

### Props

| name     | propType                                   | default | required | description                                                                  |
|----------|--------------------------------------------|---------|----------|------------------------------------------------------------------------------|
| items    | arrayOf(shape({value, label}))             |   -     |    YES   | Items to be displayed in the menu                                            |
| onClick  | func                                       |   NO    |    YES   | Called on click of button component                                          |
| onSelect | func                                       |   NO    |    YES   | Called with item.value                                                       |

---
