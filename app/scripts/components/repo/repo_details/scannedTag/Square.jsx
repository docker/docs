'use strict';

import React, { Component, PropTypes } from 'react';
const { object, bool, func, shape, oneOf, number } = PropTypes;
import classnames from 'classnames';
import styles from './Square.css';
import capitalize from 'lodash/string/capitalize';
import Tooltip from 'rc-tooltip';
import { mkComponentId } from '../nautilusUtils.js';
const severities = ['critical', 'major', 'minor', 'secure'];

export default class Square extends Component {
  static propTypes = {
    component: object.isRequired,
    isSelected: bool,
    isLayerSelected: bool,
    numVulnsBySeverity: shape({
      critical: number,
      major: number,
      minor: number
    }),
    onClick: func.isRequired,
    severity: oneOf(severities)
  }

  static defaultProps = {
    isSelected: false,
    isLayerSelected: false,
    severity: 'secure'
  }

  onClick = (e) => {
    const { isSelected, component } = this.props;
    if (isSelected) {
      //unselect component
      this.props.onClick('');
    } else {
      // Note:  this is the ID that normalizr uses to refer to components
      this.props.onClick(mkComponentId(component));
    }
  }

  mkTooltipContent = () => {
    const { component, version } = this.props.component;
    const { critical, major, minor } = this.props.numVulnsBySeverity;
    let vulnList = (
      <div>
        <div> { this.mkTooltipVulnLine('critical', critical) } </div>
        <div> { this.mkTooltipVulnLine('major', major) } </div>
        <div> { this.mkTooltipVulnLine('minor', minor) } </div>
      </div>
    );
    if (!critical && !major && !minor) {
      vulnList = <div> { this.mkTooltipVulnLine('secure') } </div>;
    }
    return (
      <div>
        <div className={ styles.component }>COMPONENT</div>
        <div>{ `${component} ${version}` }</div>
        { vulnList }
      </div>
    );
  }

  // Create line with square for inside tooltip, ex: `[] 4 Critical Vulnerabilities`
  mkTooltipVulnLine = (severity, num) => {
    if (!num && severity !== 'secure') {
      return null;
    }
    const v = num === 1 ? `Vulnerability` : `Vulnerabilities`;
    const squareClasses = classnames({
      [styles[severity]]: true,
      [styles.tooltipSquare]: true
    });
    const text = severity === 'secure' ? `No Known Vulnerabilities` : `${num} ${capitalize(severity)} ${v}`;
    return (
      <div className={ styles.vulnLineWrapper }>
        <span className={ squareClasses }>&nbsp;</span>
        <span>{ text }</span>
      </div>
    );
  }

  render() {
    const { component, version } = this.props.component;
    const tooltipId = `square-${component}-${version}`;
    const tooltipContent = this.mkTooltipContent();
    const { isLayerSelected, isSelected, severity } = this.props;
    const squareClasses = classnames({
      [styles.square]: true,
      //layer is not selected or this component is selected
      [styles.solid]: !isLayerSelected || isSelected,
      //another component is selected
      [styles.faded]: isLayerSelected && !isSelected,
      [styles.outline]: isSelected,
      [styles[severity]]: true
    });
    return (
        <Tooltip overlay={ tooltipContent }
          overlayClassName={ styles.tooltip }
          placement='top'
          align={ { overflow: { adjustY: 0 } } }
          mouseEnterDelay={0.1}>
          <div className={ squareClasses } onClick={ this.onClick }>
            &nbsp;
          </div>
        </Tooltip>
    );
  }
}
