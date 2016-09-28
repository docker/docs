[Demo](https://css-modules.github.io/webpack-demo/)

The approach is thus: Use Foundation as a "browser reset" stylesheet,
then put everything that isn't a foundation `_settings.scss` variable
in CSSModule sidecar files. This links our javascript modules with our
css and increases the ease with which we can create a module library.

# Implementation

## File Structure

```
app/scripts/
|-- ScopedSelectors.js
|-- ScopedSelectors.css
```

## Usage

```javascript
import styles from './ScopedSelectors.css';

import React, { Component } from 'react';

export default class ScopedSelectors extends Component {

  render() {
    return (
      <div className={ styles.root }>
        <p className={ styles.text }>Scoped Selectors</p>
      </div>
    );
  }

};
```

```css
.root {
  border-width: 2px;
  border-style: solid;
  border-color: #777;
  padding: 0 20px;
  margin: 0 6px;
  max-width: 400px;
}

.text {
  color: #777;
  font-size: 24px;
  font-family: helvetica, arial, sans-serif;
  font-weight: 600;
}
```

# Approach

* modules should be scoped to themselves and not affect children or
  siblings.
* Webpack already has support for css-modules in it's `css-loader`. So
  we'll start with that.
* [css-modules and preprocessors (sass)](https://github.com/css-modules/css-modules#usage-with-preprocessors)
