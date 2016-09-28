import React, { Component, PropTypes } from 'react';
import Layer from '../Layer';
import ScanHeader from '../ScanHeader';
import forEach from 'lodash/forEach';
const { array, bool, object, shape, string } = PropTypes;

/* The parent component of `NautilusScan` that is redux-aware
 * and uses `connect` can use the `mapNatilusState` selector
 * from `Nautilus/lib` to get the data in the right format
 * and pass it down to `NautilusScan` as props.
 */
class NautilusScan extends Component {
  static propTypes = {
    componentsBySeverity: shape({
      critical: array,
      major: array,
      minor: array,
      secure: array,
    }).isRequired,
    fullScan: object.isRequired,
    isFetching: bool.isRequired,
    // pathname to view all tags of this repo
    pathToTags: string.isRequired,
    vulnerabilitiesByLayer: object.isRequired,
  }

  mkLayer = (layerIndex) => {
    const { fullScan, vulnerabilitiesByLayer } = this.props;
    const layer = fullScan.layers[layerIndex];
    const layerVulnerabilities = vulnerabilitiesByLayer[layerIndex];
    const layerComponents = {};
    // layer.components is an array with ids
    forEach(layer.components, c => {
      layerComponents[c] = fullScan.components[c];
    });
    return (
      <Layer
        key={ layerIndex }
        components={ layerComponents }
        layer={ layer }
        layerNum={ layerIndex }
        vulnerabilities={ layerVulnerabilities }
      />
    );
  };

  render() {
    const {
      componentsBySeverity,
      fullScan,
      isFetching,
      pathToTags,
    } = this.props;
    const { blobs, scan_id } = fullScan;
    // No scan_id ==> first scan has failed or is in progress
    const scanError = !blobs || !scan_id;

    // Loading states and error handling should be handled by the parent
    if (isFetching || scanError) {
      return null;
    }
    // Use blobs (ordered array of layer Ids) to preserve API ordering
    return (
      <div>
        <div>
          <ScanHeader
            componentsBySeverity={componentsBySeverity}
            pathToTags={pathToTags}
            scan={fullScan}
          />
          <div>
            {blobs.map(this.mkLayer)}
          </div>
        </div>
      </div>
    );
  }
}

export default NautilusScan;
