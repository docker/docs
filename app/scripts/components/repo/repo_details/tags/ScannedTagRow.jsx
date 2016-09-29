'use strict';

import React, { PropTypes, Component } from 'react';
import { Link } from 'react-router';
const { func, object, shape, string, number, bool, instanceOf } = PropTypes;
import DeleteTagArea from './DeleteTagArea';
import { FlexRow, FlexItem } from 'common/FlexTable.jsx';
import VulnerabilityBar from './VulnerabilityBar';
import InProgressBar from './InProgressBar';
import ErrorBar from './ErrorBar';
import bytesToSize from '../../../utils/bytesToSize';
const debug = require('debug')('ScannedTagRow');
import FontAwesome from 'common/FontAwesome';
import classnames from 'classnames';
import styles from './ScannedTagRow.css';
import { StatusRecord } from 'records';
import moment from 'moment';
import { consts } from '../nautilusUtils';

//Renders a tag row for the Tags table with Nautilus Scan information
export default class ScannedTagRow extends Component {

  static propTypes = {
    actions: shape({
      deleteRepoTag: func
    }),

    tag: shape({
      name: string.isRequired,
      full_size: number,
      critical: number,
      healthy: number,
      major: number,
      minor: number
    }),
    JWT: string.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    canEdit: bool.isRequired,
    status: instanceOf(StatusRecord)
  }

  static defaultProps = {
    canEdit: false
  }

  deleteTag = (e) => {
    const { JWT, name, namespace, tag, actions } = this.props;
    const tagName = tag.name;
    actions.deleteRepoTag({ JWT, namespace, name, tagName });
  }

  render = () => {
    const {
      tag,
      canEdit,
      namespace,
      name: repoName,
      status
    } = this.props;
    const {
      name,
      full_size,
      last_scanned,
      critical,
      healthy: secure,
      major,
      minor,
      latest_scan_status
    } = tag;
    // this is the default value, indicating that no scan results exist
    const isFirstScan = moment(last_scanned).isSame('0001-01-01T00:00:00Z');
    // A scan is in progress or failed but we have stale results to show
    const scanInProgress = latest_scan_status === consts.IN_PROGRESS && !isFirstScan;
    // A scan is in progress or failed and we don't have prior results to show
    const newScanInProgress = latest_scan_status === consts.IN_PROGRESS && isFirstScan;
    const newScanFailed = latest_scan_status === consts.FAILED && isFirstScan;

    const numVulnerableComponents = critical + major + minor;
    const vulnerabilityTextClass = classnames({
      [styles.lineSpacing]: true,
      [styles.grey]: true
    });
    const vulnerabilityIconClass = classnames({
      [styles.lineSpacing]: true,
      [styles.green]: !numVulnerableComponents,
      [styles.grey]: numVulnerableComponents || newScanInProgress || newScanFailed,
      [styles.icon]: true
    });
    let icon, text, bar, lastScanned, titleLink;
    if (newScanInProgress) {
      icon = 'fa-refresh';
      text = `Scanning image for vulnerabilities...`;
      bar = <InProgressBar />;
      lastScanned = `Scan in progress...`;
      titleLink = name;
    } else if (newScanFailed) {
      icon = 'fa-minus-circle';
      text = `Could not complete image scan`;
      bar = <ErrorBar />;
      lastScanned = `Scan failed`;
      titleLink = name;
    } else {
      const numVulns = numVulnerableComponents ? 'vulnerabilities' : 'no known vulnerabilities';
      text = `This image has ${numVulns}`;
      icon = numVulnerableComponents ? 'fa-exclamation-circle' : 'fa-check';
      bar = (
        <VulnerabilityBar critical={critical}
          secure={secure}
          major={major}
          minor={minor} />
      );
      //Do not expose refresh scanner failure if scan results exist
      lastScanned = scanInProgress ? `New scan in progress...` : `Scanned ${moment(last_scanned).fromNow()}`;
      titleLink = <Link to={`/r/${namespace}/${repoName}/tags/${name}/`}>{name}</Link>;
    }
    const vulnerabilityArea = (
      <div>
        <div className={styles.linePadding}>
          <span className={vulnerabilityIconClass}>
            <FontAwesome icon={icon} />
          </span>
          <span className={vulnerabilityTextClass}>
            {text}
          </span>
        </div>
        <div className={styles.vulnerabilityBarWrapper}>
          {bar}
        </div>
      </div>
    );
    const lastScannedClasses = classnames({
      [styles.lineSpacing]: true,
      [styles.smallText]: true,
      [styles.grey]: true
    });
    const tagNameClasses = classnames({
      [styles.tagName]: true,
      [styles.grey]: newScanInProgress || newScanFailed
    });
    const titleArea = (
      <div>
        <div className={styles.linePadding}>
          <span className={tagNameClasses}>
            {titleLink}
          </span>
          <span className={styles.tagSize}>
            Compressed size: {bytesToSize(full_size)}
          </span>
        </div>
        <div className={lastScannedClasses}>
          {lastScanned}
        </div>
      </div>
    );
    let deleteArea;
    if (canEdit) {
      deleteArea = (
        <DeleteTagArea
          {...this.props}
          status={status}
          deleteTag={this.deleteTag} />
      );
    }

    return (
      <FlexRow>
        <FlexItem>{titleArea}</FlexItem>
        <FlexItem grow={2}> {vulnerabilityArea} </FlexItem>
        <FlexItem end> {deleteArea} </FlexItem>
      </FlexRow>
    );
  }
}
