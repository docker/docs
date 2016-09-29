'use strict';

const keyMirror = require('keymirror');
import forEach from 'lodash/collection/forEach';
import trim from 'lodash/string/trim';

export const mkComponentId = (component) => {
  const { component: name, version } = component;
  const id = version ? `${name}:${version}` : `${name}:`;
  return trim(id);
};

export const mapCvss = (cvss) => {
  if (cvss <= 0) {
    return 'secure';
  }
  if (cvss < 4) {
    return 'minor';
  }
  if (cvss < 7) {
    return 'major';
  }
  return 'critical';
};

// Get a map of component key to its highest CVSS score
// { bgaes: 4.5, }
export const getHighestComponentCvss = (components, vulnerabilities) => {
  let highestComponentCvss = {};
  forEach(components, (c, compName) => {
    let cvssMax = 0;
    forEach(c.vulnerabilities, v => {
      let { cvss } = vulnerabilities[v];
      if (cvss > cvssMax) {
        cvssMax = cvss;
      }
    });
    highestComponentCvss[compName] = cvssMax;
  });
  return highestComponentCvss;
};

// Given a list of severities returns the highest severity string using mapCvss
// Usage:
//   getHighestSeverity([0.5, 2, 4])
//   getHighestSeverity(...cvssList)
export const getHighestSeverity = (...cvss) => {
  if (cvss.length === 0) {
    return -1;
  }
  let maxScore = 0;
  for (let i in cvss) {
    if (cvss[i] > maxScore) {
      maxScore = cvss[i];
    }
  }
  return mapCvss(maxScore);
};

export const consts = keyMirror({
  FAILED: null,
  IN_PROGRESS: null
});
