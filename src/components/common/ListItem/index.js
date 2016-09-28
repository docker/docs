import React, { Component, PropTypes } from 'react';
import Checkbox from '../Checkbox';
import css from './styles.css';

export default class ListItem extends Component {
  static propTypes = {
    children: PropTypes.node.isRequired,
    selectable: PropTypes.bool,
    selected: PropTypes.bool,
    disabled: PropTypes.bool,
    hover: PropTypes.bool,
    onSelect: PropTypes.func,
    id: PropTypes.string,
    className: PropTypes.string,
  }

  render() {
    const {
      className = '',
      children,
      selectable,
      selected,
      disabled,
      hover,
      id,
      onSelect,
    } = this.props;

    const styles = 'dlistItem' +
      ` ${className}` +
      ` ${css.item}` +
      ` ${selectable ? css.selectable : ''}` +
      ` ${hover ? css.hover : ''}` +
      ` ${selected ? css.selected : ''}`;

    let checkbox;
    if (onSelect) {
      checkbox = (
        <Checkbox
          disabled={!!disabled}
          checked={selected}
          onCheck={(ev, checked) => onSelect(ev, checked, id)}
        />
      );
    }

    return (
      <li {...this.props} className={styles}>
        { selectable && (<div className={css.selector}>
          <div className={css.checkbox}>
            {checkbox}
          </div>
        </div>)}
        <div className={css.content}>{children}</div>
      </li>
    );
  }
}
