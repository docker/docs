import React, { Component, PropTypes } from 'react';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import mergeClasses from 'classnames';
import styles from './styles.css';

const {
  array,
  node,
  func,
  string,
  shape,
  bool,
  arrayOf,
  oneOfType,
} = PropTypes;

const itemShape = shape({
  value: string.isRequired,
  label: oneOfType([string, node]).isRequired,
  disabled: bool,
});

export default class Menu extends Component {
  static propTypes = {
    items: arrayOf(itemShape).isRequired,
    trigger: node.isRequired,
    onSelect: func.isRequired,
    className: string,
    offset: array,
  }

  state = {
    opened: false,
  }

  open = () => {
    this.setState({ opened: true });
  }

  close = () => {
    this.setState({ opened: false });
  }

  toggle() {
    if (this.state.opened) {
      this.close();
    } else {
      this.open();
    }
  }

  selectItem(value) {
    this.props.onSelect(value);
    this.close();
  }

  renderItem({ value, label, disabled }) {
    const classNames = mergeClasses(
      styles.item,
      {
        [styles.disabled]: disabled,
      }
    );
    return (
      <li
        key={value}
        className={classNames}
        onClick={() => this.selectItem(value) }
      >
        {label}
      </li>
    );
  }

  render() {
    const { trigger, items, className, offset = [0, 0] } = this.props;
    const [offsetX, offsetY] = offset;
    const { opened } = this.state;
    const triggerClone = React.cloneElement(trigger, {
      onClick: () => this.toggle(),
    });
    const classNames = mergeClasses(
      'dmenu',
      styles.container,
      className,
      { [styles.opened]: opened }
    );

    const menu = opened ? React.createElement('ul',
      {
        className: styles.menu,
        key: 'dmenu',
        style: {
          marginLeft: offsetX,
          marginTop: offsetY,
        },
      },
      items.map(this.renderItem, this)) : null;

    let backdrop;

    if (opened) {
      backdrop = (
        <div
          className={styles.backdrop}
          onClick={this.close}
          style={{ height: '100vh', width: '100vw' }}
        >
        </div>
      );
    }

    return (
      <div className={classNames}>
        {backdrop}
        {triggerClone}
        <ReactCSSTransitionGroup
          transitionName="menu"
          transitionEnterTimeout={1}
          transitionLeaveTimeout={1}
        >
          {menu}
        </ReactCSSTransitionGroup>
      </div>
    );
  }
}
