'use strict';

import React, { PropTypes, Component } from 'react';
import Card, { Block } from '@dux/element-card';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';
import {
  getComponents,
  getComponentsBySeverity,
  getLayers,
  getScan,
  getVulnerabilities,
  getVulnerabilitiesByLayer
} from './scannedTag/selectors';
import { ERROR } from 'reduxConsts';
import forEach from 'lodash/collection/forEach';
import ScanHeader from './scannedTag/ScanHeader.jsx';
import Layer from './scannedTag/Layer.jsx';
import styles from './ScannedTag.css';
import { Map } from 'immutable';
import { getStatus } from 'selectors/status';

const { object, array, instanceOf, shape } = PropTypes;

// TODO: conversion to records
// Right now we're storing only *one* scan in the reducer, therefore all entities
// can be merged in to the scan as we know they all belong to this scan.
//
// This produces a denormalized, nested scan with all entities as the child of the scan.
const getFullScan = createSelector(
  [getScan, getLayers, getComponents, getVulnerabilities],
  (scan, layers, components, vulns) => ({
    ...scan,
    layers,
    components,
    vulnerabilities: vulns
  })
);

let mapState = createStructuredSelector({
  componentsBySeverity: getComponentsBySeverity,
  scan: getFullScan,
  status: getStatus,
  vulnerabilitiesByLayer: getVulnerabilitiesByLayer
});

/**
 * This component is the detail view of particular scan/tag combination, showing
 * vulnerability and component information for a tag.
 */
@connect(mapState)
export default class ScannedTag extends Component {

  static propTypes = {
    componentsBySeverity: shape({
      critical: array,
      major: array,
      minor: array,
      secure: array
    }),
    scan: object,
    status: instanceOf(Map),
    vulnerabilitiesByLayer: object
  }

  mkLayer = (layerIndex) => {
    const { scan, vulnerabilitiesByLayer } = this.props;
    const layer = scan.layers[layerIndex];
    const layerVulnerabilities = vulnerabilitiesByLayer[layerIndex];
    let layerComponents = {};
    //layer.components is an array with ids
    forEach(layer.components, c => {
      layerComponents[c] = scan.components[c];
    });
    return (
      <Layer key={ layerIndex }
        components={ layerComponents }
        layer={ layer }
        layerNum={ layerIndex }
        vulnerabilities={ layerVulnerabilities } />
    );
  };

  render() {
    const {
      componentsBySeverity,
      params,
      namespace,
      scan,
      status,
      vulnerabilitiesByLayer
    } = this.props;
    const { blobs, reponame, tag, scan_id } = scan;
    //TODO change to use redux-simple-router params when we include it
    const ns = namespace ? namespace : params.user;
    const rn = reponame ? reponame : params.splat;
    const tn = tag ? tag : params.tagname;
    //No scan_id ==> first scan has failed or is in progress
    const scanError = !blobs || !scan_id;
    if (status.getIn(['getScanForTag', ns, rn, tn, 'status']) === ERROR || scanError) {
      return (
        <Card>
          <Block>
            <h5>Scan results unavailable.</h5>
          </Block>
        </Card>
      );
    }
    const layerInfo = <b>{`${reponame}:${tag}`}</b>;
    // blobs is an ordered array of layer ids, so we must use that to preserve API ordering
    return (
      <Card heading={<span className={styles.scanTitle}>Scan results for {layerInfo}</span>}>
        <Block>
          <div className={styles.wrapper}>
            <ScanHeader
              scan={scan}
              componentsBySeverity={componentsBySeverity}/>
            {blobs.map(this.mkLayer)}
          </div>
        </Block>
      </Card>
    );
  }
}
