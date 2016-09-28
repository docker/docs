import { createSelector, createStructuredSelector } from 'reselect';
import forEach from 'lodash/forEach';
import trim from 'lodash/trim';
import max from 'lodash/max';
import makeRepositoryId from 'lib/utils/repo-image-name';
// last argument to get is a default value if the target / path is undefined
import get from 'lodash/get';

/* Nautilus helpers */

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
  const highestComponentCvss = {};
  forEach(components, (c, compName) => {
    let cvssMax = 0;
    forEach(c.vulnerabilities, v => {
      const { cvss } = vulnerabilities[v];
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
  return mapCvss(max(cvss));
};


/* Nautilus selectors */

// TODO KRISTIE 3/24/16 Check out sharing selectors across multiple components
// https://github.com/reactjs/reselect

// Parse the url information to get the id info, and construct the path to the
// scan in the state
const getScanPath = ({ params }) => {
  // Expects the url to be of the form .../:namespace/:reponame/tags/:tag
  const { namespace, reponame, tag } = params;
  const repositoryId = makeRepositoryId({ namespace, reponame });
  return ['marketplace', 'images', repositoryId, 'scanDetail', tag];
};

// Returns the specific tag's scan for the current route
export const getScan = (state, props) => {
  const path = getScanPath(props);
  const { scan } = get(state, path, {});
  const { tag } = props.params;
  // scan is saved under
  // marketplace.images[repositoryId].scanDetail[tag].scan[tag]
  // because 'scan' is a normalizr entity

  return scan && scan[tag] || {};
};

export const getVulnerabilities = (state, props) => {
  return get(state, [...getScanPath(props), 'vulnerability'], {});
};

export const getComponents = (state, props) => {
  return get(state, [...getScanPath(props), 'component'], {});
};

export const getLayers = (state, props) => {
  return get(state, [...getScanPath(props), 'blob'], {});
};

export const getIsFetching = (state, props) => {
  return get(state, [...getScanPath(props), 'isFetching'], {});
};

export const getFullScan = createSelector(
  [getScan, getLayers, getComponents, getVulnerabilities],
  (scan, layers, components, vulns) => ({
    ...scan,
    layers,
    components,
    vulnerabilities: vulns,
  })
);

export const getLogo = (state, { params }) => {
  const { namespace, reponame } = params;
  const repositoryId = makeRepositoryId({ namespace, reponame });
  return get(state, ['marketplace', 'images', repositoryId, 'logo_url'], '');
};

export const getOrderedTagsAndScans = (state, { params, location }) => {
  const { namespace, reponame } = params;
  const { page = 1 } = location && location.query;
  const repositoryId = makeRepositoryId({ namespace, reponame });
  const path = ['marketplace', 'images', repositoryId, 'tagsAndScans'];
  const tagsAndScans = get(state, path, {});
  const { isFetching, count } = tagsAndScans;
  // Return an array in the order that was returned by the API
  const orderedIds = get(state, [...path, 'pages', page, 'orderedIds'], []);
  const results = get(state, [...path, 'pages', page], {});
  const orderedTags = orderedIds.map((tagId) => results[tagId]);
  return {
    isFetching,
    count,
    tags: orderedTags,
  };
};

/*
 * Returns a JS object containing component keys sorted into severity adjectives
 * {
 *   critical: ['openssl:2.1', 'bgaes:3.2'],
 *   major: ['pcre:3.1.2'],
 *   minor: [],
 *   secure: ['abc:new']
 * }
 */
export const getComponentsBySeverity = createSelector(
  [getComponents, getVulnerabilities],
  (comps, vulns) => {
    const map = {
      critical: [],
      major: [],
      minor: [],
      secure: [],
    };
    // For each component, get array of CVSS scores and return the highest
    forEach(comps, (componentDetail, componentKey) => {
      // componentDetail.vulnerabilities is either null or has vulns as keys
      const { vulnerabilities } = componentDetail;
      if (vulnerabilities === null || !vulnerabilities.length) {
        map.secure.push(componentKey);
      } else {
        const cvssScores = [];
        forEach(vulnerabilities, (v) => {
          cvssScores.push(vulns[v].cvss);
        });
        map[getHighestSeverity(...cvssScores)].push(componentKey);
      }
    });
    return map;
  }
);

/*
 * Returns a JS object containing layer.index as keys (sha is not unique),
 * with an object containing only that layer's vulnerabilities as the value.
 * The object with vulnerabilies will be a subset of all vulnerabilities
 * {
 *   0: [Vuln, Vuln],
 *   1: [Vuln],
 *   ...
 * }
 */
export const getVulnerabilitiesByLayer = createSelector(
  [getLayers, getComponents, getVulnerabilities],
  (layers, components, vulnerabilities) => {
    const map = {};
    forEach(layers, l => {
      const layerVulnerabilities = {};
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

// Selectors needed to `connect` a parent component of NautilusScan
export const mapNautilusState = createStructuredSelector({
  componentsBySeverity: getComponentsBySeverity,
  fullScan: getFullScan,
  isFetching: getIsFetching,
  vulnerabilitiesByLayer: getVulnerabilitiesByLayer,
});
