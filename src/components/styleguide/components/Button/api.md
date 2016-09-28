
```
<Button>Default</Button>

<Button variant="secondary" outlined icon="left">
  <DockerIcon />
  Dockerhub
</Button>

<Button inverted disabled> Dockerhub</Button>

<Button icon>
  <TutumIcon />
</Button>
```

### Props

| name     | propType                         | default | required | description                                            |
|----------|----------------------------------|---------|----------|--------------------------------------------------------|
| variant  | oneOf [<VARIANTS>]               | primary | NO       | -                                                      |
| outlined | boolean                          | false   | NO       | -                                                      |
| disabled | boolean                          | false   | NO       | -                                                      |
| inverted | boolean                          | false   | NO       | -                                                      |
| text     | boolean                          | false   | NO       | -                                                      |
| icon     | boolean OR "left" OR "right"     | false   | NO       | -                                                      |
| element  | OneOfType(element,string)        | null    | NO       | The element to be used instead of `<button>`           |

### Events

---
