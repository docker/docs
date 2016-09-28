
```
import { Pagination } from 'common';

goToPage = (page) => {
  // do something with the page number
  // this.context.router.push(...)
  // fetchData(...)
}

render() {
  const currentPage = this.props.location.query.page || 1;
  const maxVisible = 5;
  // determine last page from API response, which varies between hub and cloud
  const lastPage = total / page_size;
  return (
    <Pagination
      currentPage={currentPage}
      lastPage={lastPage}
      maxVisible={maxVisible}
      onChangePage={this.goToPage}
    />
  );

}

```

---

### Props

| name     | propType                          | default | required | description                 |
|----------|-----------------------------------|---------|----------|-----------------------------|
| `currentPage`    | oneOfType([number, string]) |    -    |    YES   | Current page number that is highlighted |
| `lastPage`    | oneOfType([number, string]) |    -    |    YES   | The last page number |
| `maxVisible`    | oneOfType([number, string]) |    10   |    NO   | The number of pages that are displayed at a given time |
| `onChangePage`    | func |    -    |    YES   | A function that handles changing the page in the router, and dispatching any actions. It will be called `onClick` of a page number with the page number as an argument |

NOTE: Route `onEnter` hooks will not be fired on a query change (such as changing `page_size` in the URL), so any fetching of data must be handled in the `onChangePage` function.

---
