import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import BackButtonArea from 'common/BackButtonArea';
import { Link } from 'react-router';
import NautilusScan from './NautilusScan';
import makeRepositoryId from 'lib/utils/repo-image-name';
import ImageHeaderWithLogo from '../ImageHeaderWithLogo';
import { createStructuredSelector } from 'reselect';
import routes from 'lib/constants/routes';
const { bool, object, shape, string } = PropTypes;
import { getLogo, mapNautilusState } from 'lib/utils/nautilus-utils';

const mapStateToProps = createStructuredSelector({
  scanDetail: mapNautilusState,
  logo: getLogo,
});

@connect(mapStateToProps)
export default class ImageDetailScanDetail extends Component {
  static propTypes = {
    logo: string,
    scanDetail: shape({
      componentsBySeverity: object,
      fullScan: object,
      isFetching: bool,
      vulnerabilitiesByLayer: object,
    }),
    params: shape({
      namespace: string,
      reponame: string,
      tag: string,
    }),
  }

  renderIsFetching = () => {
    return <div>Fetching... </div>;
  }

  render() {
    const { logo, params, scanDetail } = this.props;
    const { fullScan, isFetching } = scanDetail;
    const { namespace, reponame, tag } = params;
    const repositoryId = makeRepositoryId({ namespace, reponame });
    const overview = routes.imageDetail({ namespace, reponame });
    const tagsPath = routes.imageDetailTags({ namespace, reponame });
    const viewAllTags = <Link to={tagsPath}>View All Tags</Link>;
    // No scan_id ==> first scan has failed or is in progress
    const scanError = !isFetching && fullScan &&
      (!fullScan.blobs || !fullScan.scan_id);

    const layerInfo = `${repositoryId}:${tag}`;
    let pageContent;
    if (isFetching) {
      pageContent = this.renderIsFetching();
    } else if (scanError) {
      // TODO Kristie 3/22/16 Get designs for better error handling and 404
      pageContent = <div>{`Scan results for ${layerInfo} unavailable`}</div>;
    } else {
      pageContent = (
        <div>
          <NautilusScan pathToTags={tagsPath} {...scanDetail} />
        </div>
      );
    }
    return (
      <div>
        <BackButtonArea pathname={overview} text="Back" />
        <ImageHeaderWithLogo
          logo={logo}
          helpText={viewAllTags}
          title={layerInfo}
        />
        {pageContent}
      </div>
    );
  }
}
