import React, { Children, Component, PropTypes } from 'react';
import classnames from 'classnames';
import { uniqueId, isString } from 'lodash';
import css from './styles.css';

const {
  node,
  string,
  number,
  func,
  oneOfType,
  bool,
} = PropTypes;

export default class Tabs extends Component {
  static propTypes = {
    children: node.isRequired,
    onSelect: func.isRequired,
    className: string,
    selected: oneOfType([string, number]),
    isVertical: bool,
    icons: bool,
  }

  cloneTab(element, index) {
    const {
      onSelect,
      selected: selectedIndex,
      icons: icon,
      isVertical,
     } = this.props;
    const {
      className: tabClassName,
    } = element.props;
    const val = element.props.value || index;
    const key = element.key || uniqueId();
    const onClick = (ev) => onSelect(ev, val, index);
    const selected = selectedIndex === (isString(selectedIndex) ? val : index);
    const className = classnames({
      [css.selected]: selected,
      // Pass in the Tab's className if it exists
      [tabClassName]: !!tabClassName,
    });
    return React.cloneElement(element, {
      key,
      onClick,
      selected,
      className,
      icon,
      isVertical,
    });
  }

  render() {
    const { className = '', children, isVertical } = this.props;
    const clonedChildren = Children.map(children, this.cloneTab, this);

    const tabsClass = isVertical ? 'dvtabs' : 'dtabs';

    return (
      <div className={`${tabsClass} ${className}`}>
        {clonedChildren}
      </div>
    );
  }
}
