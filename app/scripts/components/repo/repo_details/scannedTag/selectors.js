'use strict';

import { Map } from 'immutable';
import { createSelector } from 'reselect';
import { mapCvss, getHighestSeverity } from '../nautilusUtils.js';
import values from 'lodash/object/values';
import forEach from 'lodash/collection/forEach';


// Returns the specific tag's scan for the current route
export const getScan = state => {
  // When requesting the scan from the API we only ever get one scan back.
  // This means that to get our scan we use the first value in our scan object.
  // TODO: When we move to redux-simple-router we can use the router's state
  //       to create the object key as you'd expect.
  const vals = values(state.scans.get('scan', new Map()).toJS());
  if (vals.length === 0) {
    return null;
  }
  return vals[0];
};

export const getVulnerabilities = (state) =>
  state.scans.get('vulnerability', new Map()).toJS();

export const getComponents = (state) =>
  state.scans.get('component', new Map()).toJS();

export const getLayers = (state) =>
  state.scans.get('blob', new Map()).toJS();

/**
 * Returns a JS object containing component keys sorted into severity adjectives
 * @return object  {critical: [CompKey, CompKey], major: [], minor: [], secure: []}
 */
export const getComponentsBySeverity = createSelector(
  [getComponents, getVulnerabilities],
  (comps, vulns) => {
    let map = {
      critical: [],
      major: [],
      minor: [],
      secure: []
    };
    //for each component in this scan get all the vulnerability cvss scores in an array
    //and return the highest
    forEach(comps, (componentDetail, componentKey) => {
      //componentDetail.vulnerabilities is either null or has vuln keys
      const { vulnerabilities } = componentDetail;
      if (vulnerabilities === null || !vulnerabilities.length) {
        map.secure.push(componentKey);
      } else {
        let cvssScores = [];
        forEach(vulnerabilities, (v) => {
          cvssScores.push(vulns[v].cvss);
        });
        map[getHighestSeverity(...cvssScores)].push(componentKey);
      }
    });
    return map;
  }
);

/**
 * Returns a JS object containing layer.index as keys (sha is not unique),
 * with an object containing only that layer's vulnerabilities as the value.
 * The object with vulnerabilies will be a subset of all vulnerabilities
 * @return object { 0: [Vuln, Vuln], 1: [Vuln], ... }
 */
export const getVulnerabilitiesByLayer = createSelector(
  [getLayers, getComponents, getVulnerabilities],
  (layers, components, vulnerabilities) => {
    let map = {};
    forEach(layers, l => {
      let layerVulnerabilities = {};
      forEach(l.components, c => {
        forEach(components[c].vulnerabilities, v => {
          layerVulnerabilities[v] = vulnerabilities[v];
        });
      });
      map[l.index] = layerVulnerabilities;
    });
    return map;
  }
);
