'use strict';

import React, { Component, PropTypes } from 'react';
import ComponentGrid from './ComponentGrid';
import LayerVulnerabilitiesTable from './LayerVulnerabilitiesTable';
import size from 'lodash/collection/size';
import map from 'lodash/collection/map';
const { object, shape, string, bool, func, number } = PropTypes;
import styles from './Layer.css';
import FontAwesome from 'common/FontAwesome';
import ui from 'redux-ui';
import numeral from 'numeral';
import forEach from 'lodash/collection/forEach';
import { getHighestComponentCvss, mapCvss } from '../nautilusUtils.js';
import Tooltip from 'rc-tooltip';
import classnames from 'classnames';
import { VelocityComponent } from 'velocity-react';
const debug = require('debug')('hub:Layer');

// Layer is a row in the table of layers for a ScannedTag
@ui({
  state: {
    selectedComponent: '',
    isExpanded: false,
    // if a component is selected and unselected, should the table still show?
    shouldReturnToExpanded: false
  }
})
export default class Layer extends Component {
  //vulnerabilities and components are just for this layer
  static propTypes = {
    ui: shape({
      selectedComponent: string,
      isExpanded: bool,
      shouldReturnToExpanded: bool
    }),
    updateUI: func,

    layerNum: number,
    components: object.isRequired,
    layer: object.isRequired,
    vulnerabilities: object.isRequired
  }

  truncateStringInMiddle = (str, max = 25, sep = `...`) => {
    //don't truncate if it's less than max chars
    const sepLen = sep.length;
    const len = str.length;
    if (len < max || sepLen > max) { return str; }
    const n = -0.5 * (max - len - sepLen);
    const center = len / 2;
    return str.substr(0, center - n) + sep + str.substr(len - center + n);
  }

  mkTitleRow = () => {
    const { components, layer, vulnerabilities, ui: { isExpanded } } = this.props;
    const truncateAtChars = 45;
    let titleText = layer.docker_command_line;
    if (titleText.length > truncateAtChars) {
      titleText = this.truncateStringInMiddle(layer.docker_command_line, truncateAtChars);
    }
    //add link to clickable rows
    let title = <div className={styles.dockerCommandLine}>{ titleText }</div>;
    if (size(layer.components)) {
      title = (
        <a onClick={ ::this.toggleExpanded }>
          <div className={styles.dockerCommandLine}>
            { titleText }
          </div>
        </a>
      );
    }
    //add tooltip
    const instruction = (
      <Tooltip overlay={ <div style={{ maxWidth: 400 }}>{ layer.docker_command_line }</div> }
         placement='top'
         mouseEnterDelay={ 0.1 }
         mouseLeaveDelay={ 0.3 }
         align={ { overflow: { adjustY: 0 } } }>
         { title }
      </Tooltip>
    );
    const layerSize = <div className={styles.layerSize}>Compressed size: { numeral(layer.size).format('0.0b') }</div>;
    let chevron;
    if (size(components)) {
      chevron = (
        <div className={styles.chevron} onClick={ ::this.toggleExpanded }>
          <VelocityComponent animation={{rotateX: !isExpanded ? 0 : -180}}>
            <FontAwesome icon='fa-chevron-down' />
          </VelocityComponent>
        </div>
      );
    }
    return <div className={styles.title}>{ instruction }{ layerSize } { chevron }</div>;
  }

  mkLine = (numVulnerableComps, numComps, layerType) => {
    let vulnComps, baseLayer;
    if (!numComps) {
      vulnComps = (
        <span className={styles.numVulnerableComps}>
          No components in this layer
        </span>
      );
    } else if (!numVulnerableComps && numComps) {
      vulnComps = (
        <span className={styles.numVulnerableComps}>
          <span className={styles.check}><FontAwesome icon='fa-check' /></span>
          &nbsp;No vulnerable components
        </span>
      );
    } else if (numVulnerableComps === 1) {
      vulnComps = <span className={styles.numVulnerableComps}>1 vulnerable component</span>;
    } else if (numComps) {
      //nothing for layers with no components
      vulnComps = <span className={styles.numVulnerableComps}>{numVulnerableComps} vulnerable components</span>;
    }

    // add base layer tag if relevant
    if (layerType === 'BASE') {
      baseLayer = <span className={styles.baseLayerLabel}>Base Layer</span>;
    }
    return <div>{vulnComps} {baseLayer}</div>;
  }

  toggleExpanded = () => {
    const { updateUI } = this.props;
    const { selectedComponent, isExpanded, shouldReturnToExpanded } = this.props.ui;
    if (!!selectedComponent && isExpanded) {
      //clear selected component when closing
      updateUI('selectedComponent', '');
    }
    updateUI('shouldReturnToExpanded', !shouldReturnToExpanded);
    updateUI('isExpanded', !isExpanded);
  }

  selectComponent = (componentId) => {
    const { updateUI } = this.props;
    const { isExpanded, shouldReturnToExpanded } = this.props.ui;
    if (!isExpanded && !!componentId) {
      updateUI('isExpanded', true);
    }
    //close the layer if it was closed before you selected then unselected a component
    if (isExpanded && !componentId && !shouldReturnToExpanded) {
      updateUI('isExpanded', false);
    }
    updateUI('selectedComponent', componentId);
  }

  viewAll = () => {
    const { updateUI } = this.props;
    const { selectedComponent } = this.props.ui;
    //can only be triggered when isExpanded && selectedComponent !== ''
    updateUI('selectedComponent', '');
  }

  sortComponentsBySeverity = (highestComponentCvss) => {
    //layer may not have any components
    const { components } = this.props;
    if (!size(components)) {
      return [];
    }
    return map(components, (c, key) => {
      return { ...c, key };
    }).sort( (c1, c2) => {
      return highestComponentCvss[c2.key] - highestComponentCvss[c1.key];
    });
  }

  getNumVulnerableComponents = (highestComponentCvss) => {
    const { components } = this.props;
    let numVulnComps = 0;
    forEach(components, (c, key) => {
      //get the component's highest cvss and the severity of it
      if (mapCvss(highestComponentCvss[key]) !== 'secure') {
        numVulnComps++;
      }
    });
    return numVulnComps;
  }

  render() {
    const { selectedComponent, isExpanded } = this.props.ui;
    const { components, layer, vulnerabilities } = this.props;
    // map of { componentName: highestCvss }
    const highestComponentCvss = getHighestComponentCvss(components, vulnerabilities);
    const componentsSortedBySeverity = this.sortComponentsBySeverity(highestComponentCvss);
    //nothing displayed for non-base layers with no components
    let maybeLine, maybeTable;
    const numComps = size(components);
    if (numComps && isExpanded) {
      maybeTable = (
        <LayerVulnerabilitiesTable
          componentsSortedBySeverity={ componentsSortedBySeverity }
          viewAll={ ::this.viewAll }
          isExpanded={ isExpanded }
          selectedComponent={ selectedComponent }
          vulnerabilities={ vulnerabilities } />
      );
    }
    if (!isExpanded) {
      const numVulnerableComponents = this.getNumVulnerableComponents(highestComponentCvss);
      maybeLine = this.mkLine(numVulnerableComponents, numComps, layer.type);
    }
    let maybeGrid;
    //text is displayed instead of table when there are no components (#NOP)
    if (numComps) {
      maybeGrid = (
        <ComponentGrid componentsSortedBySeverity={ componentsSortedBySeverity }
          selectedComponent={ selectedComponent }
          onClick={ ::this.selectComponent }
          vulnerabilities={ vulnerabilities } />
      );
    }
    const rowClasses = classnames({
      'row': true,
      [styles.unselectedRow]: !isExpanded,
      [styles.selectedRow]: isExpanded
    });
    return (
      <div className={rowClasses}>
        <div className='columns large-5'>
          <span className={ styles.layerNum }>{ (this.props.layerNum + 1) }</span>
          { this.mkTitleRow() }
          { maybeLine }
          { maybeTable }
        </div>
        <div className='columns large-7'>
          { maybeGrid }
        </div>
      </div>
    );
  }
}
