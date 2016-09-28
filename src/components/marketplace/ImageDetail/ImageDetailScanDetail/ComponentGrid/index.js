import React, { Component, PropTypes } from 'react';
const { array, func, object, string } = PropTypes;
import { mkComponentId, mapCvss } from 'lib/utils/nautilus-utils';
import css from './styles.css';
import Square from '../Square';
import forEach from 'lodash/forEach';
import { CRITICAL, MAJOR, MINOR, SECURE } from 'lib/constants/nautilus';

export default class ComponentGrid extends Component {
  static propTypes = {
    componentsSortedBySeverity: array.isRequired,
    onClick: func,
    selectedComponent: string,
    vulnerabilities: object,
  }

  static defaultProps = {
    selectedComponent: '',
  }

  getNumVulnsBySeverity = (component) => {
    const { vulnerabilities: vulns } = this.props;
    const numVulns = {
      [CRITICAL]: 0,
      [MAJOR]: 0,
      [MINOR]: 0,
      [SECURE]: 0,
    };
    forEach(component.vulnerabilities, (v) => {
      numVulns[mapCvss(vulns[v].cvss)]++;
    });
    return numVulns;
  }

  mkComponentSquare = (component) => {
    const { selectedComponent, onClick } = this.props;
    const componentFullName = mkComponentId(component);
    const isSelected = selectedComponent === componentFullName;
    const numVulnsBySeverity = this.getNumVulnsBySeverity(component);
    let severity = SECURE;
    if (numVulnsBySeverity.critical) {
      severity = CRITICAL;
    } else if (numVulnsBySeverity.major) {
      severity = MAJOR;
    } else if (numVulnsBySeverity.minor) {
      severity = MINOR;
    }

    return (
      <Square
        key={ componentFullName }
        component={ component }
        isSelected={ isSelected }
        isLayerSelected={ !!selectedComponent }
        numVulnsBySeverity={ numVulnsBySeverity }
        onClick={ onClick }
        severity={ severity }
      />
    );
  }

  render() {
    const { componentsSortedBySeverity } = this.props;
    return (
      <div className={ css.componentGrid }>
        { componentsSortedBySeverity.map(this.mkComponentSquare) }
      </div>
    );
  }
}
