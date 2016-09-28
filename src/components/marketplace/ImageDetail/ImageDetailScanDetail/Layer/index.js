import React, { Component, PropTypes } from 'react';
import ComponentGrid from '../ComponentGrid';
import LayerVulnerabilitiesTable from '../LayerVulnerabilitiesTable';
import size from 'lodash/size';
import map from 'lodash/map';
import filter from 'lodash/filter';
const { object, number } = PropTypes;
import css from './styles.css';
import numeral from 'numeral';
import { CheckIcon, ChevronIcon } from 'common/Icon';
import { Tooltip } from 'common';
import { getHighestComponentCvss, mapCvss } from 'lib/utils/nautilus-utils';
import classnames from 'classnames';
import { BASE, SECURE } from 'lib/constants/nautilus';
import { SMALL } from 'lib/constants/sizes';

// Layer is a row in the table of layers for a ScannedTag
export default class Layer extends Component {
  // vulnerabilities and components are just for this layer
  static propTypes = {
    components: object.isRequired,
    layer: object.isRequired,
    layerNum: number,
    vulnerabilities: object.isRequired,
  }

  state = {
    selectedComponent: '',
    isExpanded: false,
    // if a component is selected and unselected, should the table still show?
    shouldReturnToExpanded: false,
  }

  getNumVulnerableComponents = (highestComponentCvss) => {
    const { components } = this.props;
    return filter(components, (c, key) => {
      return mapCvss(highestComponentCvss[key]) !== SECURE;
    }).length;
  }

  mkLine = (numVulnerableComps, numComps) => {
    let text;
    let icon;
    if (!numComps) {
      text = 'No components in this layer';
    } else if (!numVulnerableComps && numComps) {
      text = 'No known vulnerable components';
      icon = <CheckIcon size={SMALL} className={css.check} />;
    } else if (numVulnerableComps === 1) {
      text = '1 vulnerable component';
    } else if (numComps) {
      // nothing for layers with no components
      text = `${numVulnerableComps} vulnerable components`;
    }
    return (
      <div className={css.line}>
        {icon}
        <span className={css.numVulnerableComps}>
          {text}
        </span>
      </div>
    );
  }

  mkTitleRow = () => {
    const { components, layer } = this.props;
    const { isExpanded } = this.state;
    const { type, docker_command_line: command } = layer;
    const truncateAtChars = 35;
    const titleText = this.truncateStringInMiddle(command, truncateAtChars);
    // add link to clickable rows
    let title = (
      <div className={`${css.dockerCommandLine} ${css.notClickable}`}>
        { titleText }
      </div>
    );
    if (size(layer.components)) {
      title = (
        <a onClick={ this.toggleExpanded }>
          <div className={`${css.dockerCommandLine} ${css.clickable}`}>
            { titleText }
          </div>
        </a>
      );
    }
    // add tooltip
    const tooltipContent = (
      <div className={css.tooltipContent}>{ command }</div>
    );
    const titleWithTooltip = (
      <Tooltip
        content={tooltipContent}
        mouseLeaveDelay={ 0.3 }
        align={ { overflow: { adjustY: 0 } } }
        theme="dark"
      >
        { title }
      </Tooltip>
    );
    const layerSize = (
      <div className={css.layerSize}>
        { numeral(layer.size).format('0.0b') }
      </div>
    );
    let chevron;
    if (size(components)) {
      const facesUp = isExpanded ? css.chevronUp : '';
      const chevronClasses = `${css.chevron} ${css.clickable} ${facesUp}`;
      chevron = (
        <div onClick={ this.toggleExpanded }>
          <ChevronIcon className={chevronClasses} />
        </div>
      );
    }
    let baseLayer;
    if (type === BASE) {
      baseLayer = <div className={css.baseLayerLabel}>base layer</div>;
    }
    return (
      <div className={css.title}>
        <div className={css.information}>
          { titleWithTooltip }
          { layerSize }
          { baseLayer }
        </div>
        { chevron }
      </div>
    );
  }

  selectComponent = (componentId) => {
    const { isExpanded, shouldReturnToExpanded } = this.state;
    const nextState = { selectedComponent: componentId };
    if (!isExpanded && !!componentId) {
      nextState.isExpanded = true;
    }
    if (isExpanded && !componentId && !shouldReturnToExpanded) {
      // close the layer if it was closed before you selected then
      // unselected a component
      nextState.isExpanded = false;
    }
    this.setState(nextState);
  }

  sortComponentsBySeverity = (highestComponentCvss) => {
    // layer may not have any components
    const { components } = this.props;
    if (!size(components)) {
      return [];
    }
    return map(components, (c, key) => {
      return { ...c, key };
    }).sort((c1, c2) => {
      return highestComponentCvss[c2.key] - highestComponentCvss[c1.key];
    });
  }

  toggleExpanded = () => {
    const {
      isExpanded,
      selectedComponent,
      shouldReturnToExpanded,
    } = this.state;
    const nextState = {
      isExpanded: !isExpanded,
      shouldReturnToExpanded: !shouldReturnToExpanded,
    };
    if (!!selectedComponent && isExpanded) {
      // clear selected component when closing
      nextState.selectedComponent = undefined;
    }
    this.setState(nextState);
  }

  truncateStringInMiddle = (str, max = 25, sep = '...') => {
    // don't truncate if it's less than max chars
    const sepLen = sep.length;
    const len = str.length;
    if (len < max || sepLen > max) { return str; }
    const n = -0.5 * (max - len - sepLen);
    const center = len / 2;
    return str.substr(0, center - n) + sep + str.substr(len - center + n);
  }

  render() {
    const { isExpanded, selectedComponent } = this.state;
    const { components, layer, vulnerabilities } = this.props;
    // map of { componentName: highestCvss }
    const highestComponentCvss =
      getHighestComponentCvss(components, vulnerabilities);
    const componentsSortedBySeverity =
      this.sortComponentsBySeverity(highestComponentCvss);
    // nothing is displayed for non-base layers with no components
    let maybeLine;
    let maybeTable;
    const numComps = size(components);
    if (numComps && isExpanded) {
      maybeTable = (
        <LayerVulnerabilitiesTable
          componentsSortedBySeverity={ componentsSortedBySeverity }
          isExpanded={ isExpanded }
          selectedComponent={ selectedComponent }
          vulnerabilities={ vulnerabilities }
        />
      );
    }
    if (!isExpanded) {
      const numVulnerableComponents =
        this.getNumVulnerableComponents(highestComponentCvss);
      maybeLine = this.mkLine(numVulnerableComponents, numComps, layer.type);
    }
    let maybeGrid;
    // text is displayed instead of table when there are no components (#NOP)
    if (numComps) {
      maybeGrid = (
        <ComponentGrid componentsSortedBySeverity={ componentsSortedBySeverity }
          selectedComponent={ selectedComponent }
          onClick={ this.selectComponent }
          vulnerabilities={ vulnerabilities }
        />
      );
    }
    const classes = classnames({
      [css.layer]: true,
      [css.layerWrapper]: true,
      [css.selectedLayer]: isExpanded,
    });
    return (
      <div className={classes}>
        <section className={css.section}>
          { this.mkTitleRow() }
          { maybeLine }
          { maybeTable }
        </section>
        <section>
          { maybeGrid }
        </section>
      </div>
    );
  }
}
