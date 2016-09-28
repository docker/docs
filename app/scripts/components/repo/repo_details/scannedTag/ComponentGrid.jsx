'use strict';

import React, { Component, PropTypes } from 'react';
const { object, array, func, string } = PropTypes;
import { mapCvss, getHighestSeverity } from '../nautilusUtils.js';
import styles from './ComponentGrid.css';
import Square from './Square';
import forEach from 'lodash/collection/forEach';
import { mkComponentId } from '../nautilusUtils';

export default class ComponentGrid extends Component {
  static propTypes = {
    componentsSortedBySeverity: array.isRequired,
    onClick: func,
    selectedComponent: string,
    vulnerabilities: object
  }

  static defaultProps = {
    selectedComponent: ''
  }

  getNumVulnsBySeverity = (component) => {
    const { vulnerabilities: vulns } = this.props;
    let numVulns = {
      critical: 0,
      major: 0,
      minor: 0,
      secure: 0
    };
    forEach(component.vulnerabilities, (v) => {
      numVulns[mapCvss(vulns[v].cvss)]++;
    });
    return numVulns;
  }

  mkComponentSquare = (component) => {
    const { selectedComponent, onClick, vulnerabilities: vulns } = this.props;
    const componentFullName = mkComponentId(component);
    const isSelected = selectedComponent === componentFullName;
    const numVulnsBySeverity = this.getNumVulnsBySeverity(component);
    let severity = 'secure';
    if (numVulnsBySeverity.critical) {
      severity = 'critical';
    } else if (numVulnsBySeverity.major) {
      severity = 'major';
    } else if (numVulnsBySeverity.minor) {
      severity = 'minor';
    }

    return (
      <Square key={ componentFullName }
        component={ component }
        isSelected={ isSelected }
        isLayerSelected={ !!selectedComponent }
        numVulnsBySeverity={ numVulnsBySeverity }
        onClick={ onClick }
        severity={ severity } />
    );
  }

  render() {
    const { componentsSortedBySeverity } = this.props;
    return (
      <div className={ styles.componentGrid }>
          { componentsSortedBySeverity.map(this.mkComponentSquare) }
      </div>
    );
  }
}
