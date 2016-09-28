import React, { Component, PropTypes } from 'react';
import RCTooltip from 'rc-tooltip';
import './bootstrap_white.css';
import './styles.css';
import classnames from 'classnames';
const {
  arrayOf,
  node,
  number,
  object,
  oneOf,
  shape,
  string,
} = PropTypes;
const triggers = ['hover', 'click', 'focus'];
const placements = ['left', 'right', 'top', 'bottom',
  'topLeft', 'topRight', 'bottomLeft', 'bottomRight'];
const themes = ['white', 'dark'];

/* Uses react-component/tooltip
 * https://github.com/react-component/tooltip
 * for a longer description of available API and propTypes
 */
export default class Tooltip extends Component {
  static propTypes = {
    // https://github.com/yiminghe/dom-align
    align: shape({
      offset: object,
      targetOffset: object,
    }),
    children: node.isRequired,
    className: string,
    content: node.isRequired,
    mouseEnterDelay: number,
    mouseLeaveDelay: number,
    placement: oneOf(placements),
    theme: oneOf(themes),
    trigger: arrayOf(oneOf(triggers)),
  }

  static defaultProps = {
    mouseEnterDelay: 0.1,
    mouseLeaveDelay: 0.1,
    placement: 'top',
    theme: 'white',
    trigger: ['hover'],
  }

  render() {
    const {
      align,
      className,
      content,
      mouseEnterDelay,
      mouseLeaveDelay,
      placement,
      theme,
      trigger,
    } = this.props;

    const classes = classnames({
      darkTheme: theme === 'dark',
      [className]: !!className,
    });

    return (
      <RCTooltip
        align={align}
        mouseEnterDelay={mouseEnterDelay}
        mouseLeaveDelay={mouseLeaveDelay}
        overlay={content}
        overlayClassName={classes}
        placement={placement}
        prefixCls="dTooltip"
        trigger={trigger}
      >
        { this.props.children }
      </RCTooltip>
    );
  }
}
