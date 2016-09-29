'use strict';

import React, { createClass } from 'react';

export default createClass({
  displayName: 'Row',
  render() {
    return (
      <div className='row'>{this.props.children}</div>
    );
  }
});
