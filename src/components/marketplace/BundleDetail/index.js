import React, { Component, PropTypes, Children } from 'react';
import { connect } from 'react-redux';

const { node } = PropTypes;

@connect()
export default class BundleDetail extends Component {
  static propTypes = {
    children: node.isRequired,
  }

  render() {
    const { children } = this.props;
    const childrenWithProps = Children.map(children, (child) => {
      return React.cloneElement(child, this.props);
    });
    return (
      <div id="module-images-detail">
        {childrenWithProps}
      </div>
    );
  }
}
