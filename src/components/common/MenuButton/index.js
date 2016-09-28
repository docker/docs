import React, { Component, PropTypes } from 'react';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import classnames from 'classnames';
import styles from './styles.css';

import { Button, ChevronIcon } from 'common';

const {
  any,
  array,
  func,
  string,
  shape,
  arrayOf,
} = PropTypes;

const itemShape = shape({
  value: string.isRequired,
  label: string.isRequired,
});

export default class MenuButton extends Component {
  static propTypes = {
    className: string,
    children: any,
    onClick: func.isRequired,
    items: arrayOf(itemShape).isRequired,
    onSelect: func.isRequired,
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

  toggle = () => {
    if (this.state.opened) {
      this.close();
    } else {
      this.open();
    }
  }

  selectItem = (value) => () => {
    this.props.onSelect(value);
    this.close();
  }

  renderItems = () => {
    const { items } = this.props;
    return items.map(({ value, label }, idx) => {
      return (
        <li
          key={idx}
          className={styles.item}
          onClick={this.selectItem(value) }
        >
          <div className={styles.dropText}>
            {label}
          </div>
        </li>
      );
    });
  }

  renderMenu = () => {
    const { offset = [0, 0] } = this.props;
    const [offsetX, offsetY] = offset;

    return (
      <ul
        className={styles.menu}
        style={{ marginLeft: offsetX, marginTop: offsetY }}
      >
        {this.renderItems()}
      </ul>
    );
  }

  render() {
    const {
      children,
      className,
      onClick,
    } = this.props;
    const { opened } = this.state;

    const buttonStyles = classnames({
      [className]: className,
      [styles.splitWrap]: true,
    });

    let menu;
    let open;
    let backdrop;

    if (opened) {
      menu = this.renderMenu();
      open = styles.open;
      backdrop = (
        <div
          onClick={() => {this.close();}}
          className={styles.backdrop}
        ></div>
      );
    }

    return (
      <div className={buttonStyles}>
        <div className={styles.splitRow}>
          <Button
            className={`${styles.btnSplit} ${open}`}
            onClick={onClick}
          >
            { children }
          </Button>
          <Button
            className={`${styles.menuSplit} ${open}`}
            onClick={this.toggle}
          >
            <ChevronIcon />
          </Button>
          { backdrop }
        </div>
        <div className={styles.splitRow}>
          <ReactCSSTransitionGroup
            transitionName="menu"
            transitionEnterTimeout={1}
            transitionLeaveTimeout={1}
            className={styles.menuWrap}
          >
            {menu}
          </ReactCSSTransitionGroup>
        </div>
      </div>
    );
  }
}
