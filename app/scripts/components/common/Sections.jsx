'use strict';

import React, { Component, PropTypes } from 'react';
const { string, node, object, oneOfType, bool } = PropTypes;
import Card, { Block } from '@dux/element-card';
import styles from './Sections.css';

export class SplitSection extends Component {
  static propTypes = {
    title: string,
    subtitle: node,
    /**
     * module is a props passthrough for precision control over how
     * Module displays. An example is setting the `intent` on a module.
     */
    module: oneOfType([object, bool])
  };
  static defaultProps = {
    module: true
  };

  render() {
    var children;
    if (this.props.module) {
      children = (
        <Card>
          <Block>
            {this.props.children}
          </Block>
        </Card>
      );
    } else {
      children = this.props.children;
    }

    return (
      <div className='row'>
        <div className="columns large-4">
          <h5>{this.props.title}</h5>
          <div>
            {this.props.subtitle}
          </div>
        </div>
        <div className="columns large-7 large-offset-1">
          { children }
        </div>
      </div>
    );
  }
}

export class FullSection extends Component {
  static propTypes = {
    title: string
  };

  render() {
    return (
      <div className='row'>
        <div className="columns large-12">
          <h5>{this.props.title}</h5>
        </div>
        {this.props.children}
      </div>
    );
  }
}
