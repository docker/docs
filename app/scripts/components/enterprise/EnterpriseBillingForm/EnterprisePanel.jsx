'use strict';

import React, {
  PropTypes,
  Component
  } from 'react';

const { string } = PropTypes;

export default class EnterprisePanel extends Component {

  static propTypes = {
    heading: string
  }

  render() {
    return (
      <div className='row'>
        <div className='large-12 columns'>
          <h5>{this.props.heading}</h5><br/>
          <div>{this.props.children}</div>
        </div>
      </div>
    );
  }
}
