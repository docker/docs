import React, { Component, PropTypes } from 'react';
import { CheckIcon, DropdownIcon } from 'common/Icon';
import { mapCvss, mkComponentId } from 'lib/utils/nautilus-utils';
import { Tooltip } from 'common';
import capitalize from 'lodash/capitalize';
import classnames from 'classnames';
import css from './styles.css';
import filter from 'lodash/filter';
import forEach from 'lodash/forEach';
const { object, bool, array, string } = PropTypes;

export default class LayerVulnerabilitiesTable extends Component {
  static propTypes = {
    componentsSortedBySeverity: array.isRequired,
    selectedComponent: string,
    isExpanded: bool,
    vulnerabilities: object.isRequired,
  }

  static defaultProps = {
    selectedComponent: '',
  }

  state = {
    // default to hiding secure components
    hideSecureComponents: true,
  }

  toggleHideSecureComponents = () => {
    this.setState({
      hideSecureComponents: !this.state.hideSecureComponents,
    });
  }

  mkShowAllButton = () => {
    const { selectedComponent } = this.props;
    // only render a button when a component is *not* selected
    if (selectedComponent) {
      return null;
    }
    const { hideSecureComponents } = this.state;
    const allOrFewer = hideSecureComponents ? 'all' : 'fewer';
    const iconClass = hideSecureComponents ? css.icon : css.iconUp;
    return (
      <div
        className={css.showAllButton}
        onClick={this.toggleHideSecureComponents}
      >
        <span className={css.showAll}>{`Show ${allOrFewer} components`}</span>
        <DropdownIcon size="small" className={iconClass} />
      </div>
    );
  }

  mkHeader = () => {
    return (
      <div className={css.row}>
        <div className={css.headerText}>Component</div>
        <div className={css.headerText}>Vulnerability</div>
        <div className={`${css.headerText} ${css.severity}`}>Severity</div>
      </div>
    );
  }

  mkComponentArea = (fullName, license, license_type) => {
    return (
      <div className={css.wrapWords}>
        <div className={css.componentName}>
          { fullName }
        </div>
        <div className={css.license}>
          { `${license}: ${capitalize(license_type)} License` }
        </div>
      </div>
    );
  }

  mkVulnerabilityTooltipContent = ({ cve, summary }) => {
    return (
      <div className={css.tooltipContent}>
        <div className={css.tooltipTitle}>{ cve }</div>
        <div className={css.tooltipText}>{ summary }</div>
      </div>
    );
  }

  mkVulnerabilityArea = (fullName, vulns) => {
    if (!vulns.length) {
      return (
        <div className={css.cve}>
          No known vulnerabilities
        </div>
      );
    }
    const lines = vulns.map(v => {
      const tooltipContent = this.mkVulnerabilityTooltipContent(v);
      return (
        <Tooltip
          key={`${fullName}-${v.cve}-cvss`}
          content={tooltipContent}
          mouseEnterDelay={ 0.2 }
          mouseLeaveDelay={ 0.2 }
        >
          <div
            key={`${fullName}-${v.cve}`}
            className={css.cve}
          >
            <a href={ `https://cve.mitre.org/cgi-bin/cvename.cgi?name=${v.cve}` }
              target="_blank"
            >
              { v.cve }
            </a>
          </div>
        </Tooltip>
      );
    });
    return <div>{lines}</div>;
  }

  mkSeverityArea = (fullName, vulns) => {
    if (!vulns.length) {
      return <div className={`${css.secure} ${css.severity}`}>N/A</div>;
    }
    const lines = vulns.map(v => {
      const classes = classnames({
        [css[v.severity]]: true,
        [css.vulnerabilityLines]: true,
        [css.severity]: true,
      });
      return (
        <div key={`${fullName}-${v.cve}-cvss`} className={classes}>
          { capitalize(v.severity) }
        </div>
      );
    });
    return <div>{lines}</div>;
  }

  // sort vulnerabilities for this component
  mkSortedComponentVulns = (vulnerabilities, layerVulnerabilities) => {
    const componentVulns = [];
    forEach(vulnerabilities, (v) => {
      const vuln = layerVulnerabilities[v];
      // add severity key to vuln object: ex. severity: 'minor'
      componentVulns.push({ ...vuln, severity: mapCvss(vuln.cvss) });
    });
    // sort component vulnerabilities by cvss so most vulnerable is first
    componentVulns.sort((v1, v2) => v2.cvss - v1.cvss);
    return componentVulns;
  }

  mkComponentRow = (component) => {
    const {
      license,
      license_type,
      component: name,
      version,
      vulnerabilities,
    } = component;
    const { vulnerabilities: layerVulnerabilities } = this.props;

    const fullName = version ? `${name} ${version}` : name;
    const componentArea =
      this.mkComponentArea(fullName, license, license_type);
    const componentVulns =
      this.mkSortedComponentVulns(vulnerabilities, layerVulnerabilities);
    const vulnerabilityArea =
      this.mkVulnerabilityArea(fullName, componentVulns);
    const severityArea = this.mkSeverityArea(fullName, componentVulns);
    return (
      <div className={`${css.row} ${css.componentRow}`} key={fullName}>
        {componentArea}
        {vulnerabilityArea}
        {severityArea}
      </div>
    );
  }

  render() {
    const { componentsSortedBySeverity } = this.props;
    const componentsToShow = filter(componentsSortedBySeverity, (comp) => {
      const { selectedComponent } = this.props;
      const { vulnerabilities } = comp;
      const { hideSecureComponents } = this.state;
      const id = mkComponentId(comp);
      // is this is a secure component and should we hide it?
      const hideSecureComp = hideSecureComponents && !vulnerabilities.length;
      // is another component selected? --> hide this one
      const showOtherComp = selectedComponent && selectedComponent !== id;
      // always show a component that is selected
      const showThisComp = selectedComponent && selectedComponent === id;
      return (!hideSecureComp && !showOtherComp) || showThisComp;
    });
    let components;
    if (!componentsToShow.length) {
      components = (
        <div className={css.noVulnerableComponents}>
          <CheckIcon size="small" />
          No known vulnerable components
        </div>
      );
    } else {
      components = componentsToShow.map(this.mkComponentRow);
    }
    return (
      <div>
        { this.mkHeader() }
        { components }
        { this.mkShowAllButton() }
      </div>
    );
  }
}
