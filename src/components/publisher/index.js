import React, { Component, PropTypes } from 'react';


const { node } = PropTypes;

export default class Publisher extends Component {
  static propTypes = {
    children: node.isRequired,
  }

  render() {
    return (
      <div>
        {this.props.children}
      </div>
    );
  }
}
