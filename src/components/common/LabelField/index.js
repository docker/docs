import React, { Component, PropTypes } from 'react';
import css from './styles.css';
const { bool, element, string } = PropTypes;
import mergeClasses from 'classnames';

export default class LabelField extends Component {
  static propTypes = {
    children: element.isRequired,
    className: string,
    label: string.isRequired,
    disabled: bool,
  }

  render() {
    const { className, children, disabled, label } = this.props;

    const classNames = mergeClasses('dlabelField', {
      [css.disabled]: disabled,
    }, css.labelFieldContainer, className);

    return (
      <div className={classNames}>
        <span className={css.label}>{ label }</span>
        <div className={css.field}>{ children }</div>
      </div>
    );
  }
}
