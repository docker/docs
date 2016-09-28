import React, { PropTypes } from 'react';
import css from './styles.css';

const { node } = PropTypes;

export default class Admin extends React.Component {
  static propTypes = {
    children: node.isRequired,
  }

  render() {
    return (
      <div className={css.admin}>
        {this.props.children}
      </div>
    );
  }
}
