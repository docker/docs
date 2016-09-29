'use strict';

import React, { Component, PropTypes } from 'react';
const { object, bool, array, func, string } = PropTypes;
import styles from './LayerVulnerabilitiesTable.css';
import { mapCvss, mkComponentId } from '../nautilusUtils.js';
import forEach from 'lodash/collection/forEach';
import capitalize from 'lodash/string/capitalize';
import FontAwesome from 'common/FontAwesome';
import classnames from 'classnames';
import Tooltip from 'rc-tooltip';

export default class LayerVulnerabilitiesTable extends Component {
  static propTypes = {
    componentsSortedBySeverity: array.isRequired,
    selectedComponent: string,
    isExpanded: bool,
    viewAll: func,
    vulnerabilities: object.isRequired
  }

  static defaultProps = {
    selectedComponent: ''
  }

  mkHeader = () => {
    return (
      <div className='row'>
        <span className={styles.headerText}>
          <div className='columns large-5'>Component</div>
          <div className='columns large-4'>Vulnerability</div>
          <div className={`columns large-3 ${styles.severity}`}>Severity</div>
        </span>
      </div>
    );
  }

  mkComponentArea = (fullName, license, license_type) => {
    return (
      <div className={styles.wrapWords}>
        <div><b>{ fullName }</b></div>
        <div className={styles.licenseType}>
          { `${license}:${capitalize(license_type)} License` }
        </div>
      </div>
    );
  }

  mkVulnerabilityArea = (fullName, vulns) => {
    if (!vulns.length) {
      return <div className={styles.vulnerabilityLines}>No known vulnerabilities</div>;
    }
    return vulns.map( v => {
      const tooltipContent = (
        <div style={ { maxWidth: 300 } }>
          <div><b>{ v.cve }</b></div>
          <div>{ v.summary }</div>
        </div>
      );
      return (
      <Tooltip key={ fullName + '-' + v.cve + '-cvss' }
        overlay={ tooltipContent }
        placement='top'
        mouseEnterDelay={ 0.2 }
        mouseLeaveDelay={ 0.2 }
        align={ { overflow: { adjustY: 0 } } }
        trigger={ ['hover'] }>
        <div className={styles.vulnerabilityLines} key={ fullName + '-' + v.cve }>
          <a href={ `https://cve.mitre.org/cgi-bin/cvename.cgi?name=${v.cve}` }
            target='_blank'>
            { v.cve }
          </a>
        </div>
      </Tooltip>
      );
    });
  }

  mkSeverityArea = (fullName, vulns) => {
    if (!vulns.length) {
      let classes = classnames({
        [styles.secure]: true,
        [styles.vulnerabilityLines]: true,
        [styles.severity]: true
      });
      return <div className={classes}>N/A</div>;
    }
    return vulns.map( v => {
      let classes = classnames({
        [styles[v.severity]]: true,
        [styles.vulnerabilityLines]: true,
        [styles.severity]: true
      });
      return (
        <div key={ fullName + '-' + v.cve + '-cvss' } className={classes}>
          { capitalize(v.severity) }
        </div>
      );
    });
  }

  mkComponentRow = (component) => {
    const {
      license,
      license_type,
      component: name,
      version,
      vulnerabilities
    } = component;
    const { selectedComponent, vulnerabilities: layerVulnerabilities } = this.props;
    const id = mkComponentId(component);
    // only render this row if there is NOT a selected component OR this is the selected component
    let content;
    if (!selectedComponent || selectedComponent === id) {
      const fullName = version ? `${name} ${version}` : name;
      const componentArea = this.mkComponentArea(fullName, license, license_type);
      let componentVulns = [];
      forEach(vulnerabilities, (v) => {
        const vuln = layerVulnerabilities[v];
        //add severity key to vuln object: ex. severity: 'minor'
        componentVulns.push({ ...vuln, severity: mapCvss(vuln.cvss)});
      });
      //sort component vulnerabilities by cvss so most vulnerable is first
      componentVulns.sort((v1, v2) => v2.cvss - v1.cvss);
      const vulnerabilityArea = this.mkVulnerabilityArea(fullName, componentVulns);
      const severityArea = this.mkSeverityArea(fullName, componentVulns);
      content = (
        <div className='row' key={fullName}>
          <div className='columns large-12'><hr className={styles.border} /></div>
          <div className='columns large-5'>{componentArea}</div>
          <div className='columns large-4'>{vulnerabilityArea}</div>
          <div className='columns large-3'>{severityArea}</div>
        </div>
      );
    }
    return content;
  }

  render() {
    const {
      componentsSortedBySeverity,
      selectedComponent,
      isExpanded,
      viewAll
    } = this.props;
    let viewAllSection;
    if (selectedComponent) {
      viewAllSection = (
        <div className={styles.viewAll} onClick={viewAll}>
          { `View All `}<FontAwesome icon='fa-chevron-down' />
        </div>
      );
    }
    return (
      <div>
        { this.mkHeader() }
        { componentsSortedBySeverity.map(this.mkComponentRow) }
        { viewAllSection }
      </div>
    );
  }
}
