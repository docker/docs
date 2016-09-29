'use strict';

import React, { PropTypes, Component } from 'react';
import { StrippedModule } from 'dux';
import styles from './ListSelector.css';
import _ from 'lodash';

const { string, arrayOf, node } = PropTypes;

export default class ListSelector extends Component {

  static propTypes = {
    header: node,
    items: arrayOf(node)
  }

  render() {
    var header = null;
    if (_.isString(this.props.header)) {
      header = (<h5 className={styles.header}>{this.props.header}</h5>);
    } else if (_.isObject(this.props.header)) {
      header = <div className={styles.header}>{this.props.header}</div>;
    }

    return (
      <StrippedModule>
        {header}
        <ul className={styles.listSelectorItems}>
          {this.props.items}
        </ul>
      </StrippedModule>);
  }
}
