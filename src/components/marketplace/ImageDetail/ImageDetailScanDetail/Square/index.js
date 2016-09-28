import React, { Component, PropTypes } from 'react';
const { bool, func, number, object, oneOf, shape } = PropTypes;
import classnames from 'classnames';
import css from './styles.css';
import capitalize from 'lodash/capitalize';
import { Tooltip } from 'common';
import { mkComponentId } from 'lib/utils/nautilus-utils';
import {
  CRITICAL,
  MAJOR,
  MINOR,
  SECURE,
  severities,
} from 'lib/constants/nautilus';

export default class Square extends Component {
  static propTypes = {
    component: object.isRequired,
    isSelected: bool,
    isLayerSelected: bool,
    numVulnsBySeverity: shape({
      critical: number,
      major: number,
      minor: number,
    }),
    onClick: func.isRequired,
    severity: oneOf(severities),
  }

  static defaultProps = {
    isSelected: false,
    isLayerSelected: false,
    severity: SECURE,
  }

  onClick = () => {
    const { isSelected, component } = this.props;
    if (isSelected) {
      // unselect component
      this.props.onClick(undefined);
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
        { this.mkTooltipVulnLine(CRITICAL, critical) }
        <div> { this.mkTooltipVulnLine(MAJOR, major) } </div>
        <div> { this.mkTooltipVulnLine(MINOR, minor) } </div>
      </div>
    );
    if (!critical && !major && !minor) {
      vulnList = <div> { this.mkTooltipVulnLine(SECURE) } </div>;
    }
    return (
      <div>
        <div className={css.tooltipTitle}>Component</div>
        <div className={css.tooltipSubheader}>{`${component} ${version}`}</div>
        { vulnList }
      </div>
    );
  }

  // Create text for tooltip body: `4 Critical Vulnerabilities`
  mkTooltipVulnLine = (severity, num) => {
    if (!num && severity !== SECURE) {
      return null;
    }
    const v = num === 1 ? 'Vulnerability' : 'Vulnerabilities';
    const text = severity === SECURE ?
      'No Known Vulnerabilities'
      : `${num} ${capitalize(severity)} ${v}`;
    return <div className={css.tooltipText}>{ text }</div>;
  }

  render() {
    const { isLayerSelected, isSelected, severity } = this.props;
    const squareClasses = classnames({
      [css.square]: true,
      // layer is not selected or this component is selected
      [css.solid]: !isLayerSelected || isSelected,
      // another component is selected
      [css.faded]: isLayerSelected && !isSelected,
      [css.outline]: isSelected,
      [css[severity]]: true,
    });
    const tooltipContent = this.mkTooltipContent();
    return (
      <Tooltip content={tooltipContent}>
        <div className={ squareClasses } onClick={ this.onClick }>
          &nbsp;
        </div>
      </Tooltip>
    );
  }
}
