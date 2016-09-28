import React, { Children, Component, PropTypes } from 'react';
import css from './styles.css';

export default class List extends Component {
  static propTypes = {
    children: PropTypes.node.isRequired,
    selectable: PropTypes.bool,
    hover: PropTypes.bool,
    className: PropTypes.string,
  }

  render() {
    const { selectable, hover, className = '' } = this.props;
    const children = Children.map(this.props.children,
      (element) => React.cloneElement(element, { selectable, hover }));

    return (
      <div className={`dlist ${css.container} ${className}`}>
        <ul className={css.list}>{children}</ul>
      </div>
    );
  }
}
