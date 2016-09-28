import React, { Component, PropTypes } from 'react';
import mergeClasses from 'classnames';
import Button from '../Button';
import css from './styles.css';

const { oneOfType, string, bool, node, number } = PropTypes;

export default class Tab extends Component {
  static propTypes = {
    value: oneOfType([string, number]),
    element: node,
    children: node,
    className: string,
    selected: bool,
    icon: bool,
    isVertical: bool,
    fullWidth: bool,
    isUnderlined: bool,
  }

  static defaultProps = {
    isUnderlined: false,
  }

  render() {
    const {
      className = '',
      selected,
      children,
      icon,
      isVertical,
      isUnderlined,
    } = this.props;
    const classPrefix = isVertical ? 'dvtab' : 'dtab';
    const classNames = mergeClasses(classPrefix, className, {
      [css.tab]: true,
      [css.addBorder]: true,
      [css.noBorder]: !isUnderlined && !selected,
      [css.icon]: !!icon,
      [css.selected]: selected && !isVertical,
      [css.vselected]: selected && isVertical,
      [className]: !!className,
    });

    const props = {
      ...this.props,
      className: classNames,
      icon: false,
    };

    return (
      <Button {...props}>{children}</Button>
    );
  }
}
