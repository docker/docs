
```
import { Modal, Button, Card } from 'common';

openModal = () => {
  this.setState({ isOpen: true });
}

closeModal = () => {
  this.setState({ isOpen: false });
}

<Button onClick={this.openModal}>Open Modal</Button>

<Modal
  isOpen={this.state.isOpen}
  onRequestClose={this.closeModal}
>
  <Card
    title="Hello World"
    ghost
  >
    <Button
      outlined
      variant="panic"
      onClick={this.onCancel}
    >Cancel</Button>
  </Card>
</Modal>

```

### Props

| name           | propType                            | default | required | description                                                                  |
|----------------|-------------------------------------|---------|----------|------------------------------------------------------------------------------|
| isOpen         | bool                                |  FALSE  |    YES   | Determines whether the modal is open or not                                  |
| onRequestclose | func                                |   -     |    YES   | Closes the modal                                                             |

---
