'use strict';

import cn from 'classnames';
import React, { Component, PropTypes } from 'react';
import styles from './repositoryName.css';
export default class RepositoryName extends Component {
  static propTypes = {
    className: PropTypes.string,
    children: PropTypes.node.isRequired,
    inlined: PropTypes.bool
  }

  static defaultProps = {
    inlined: false
  }

  render() {
    const {
      className,
      children,
      inlined
    } = this.props;
    if (children.length !== 2) {
      throw new Error('Should be in the format of <RepositoryName>{namespace}{name}</RepositoryName>');
    }
    const [ namespace, name ] = children;
    const classNames = cn({
      [styles.RepositoryName]: true,
      [styles.inlined]: inlined
    }, className);
    return (
      <div className={ classNames }>
        { namespace }
        <div className={ styles.divider }>/</div>
        { name }
      </div>
    );
  }
}
