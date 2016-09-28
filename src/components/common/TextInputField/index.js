import React, { Component, PropTypes } from 'react';
import css from './styles.css';
const { array, element, oneOf, string } = PropTypes;

export default class LabelField extends Component {
  static propTypes = {
    children: oneOf(array, element),
    className: string,
    label: string,
    hintText: string,
  }

  render() {
    const { className, children, label } = this.props;

    return (
      <div className={className || css.labelFieldContainer}>
        <span className={css.label}>{ label }</span>
        <span className={css.field}>{ children }</span>
      </div>
    );
  }
}
