import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import { CheckIcon } from 'common/Icon';
import moment from 'moment';
import { IN_PROGRESS, scanStatuses } from 'lib/constants/nautilus';
import { SMALL } from 'lib/constants/sizes';
const { array, oneOf, shape, string } = PropTypes;

export default class ScanHeader extends Component {

  static propTypes = {
    componentsBySeverity: shape({
      critical: array,
      major: array,
      minor: array,
      secure: array,
    }).isRequired,
    scan: shape({
      completed_at: string,
      latest_scan_status: oneOf(scanStatuses),
      reponame: string,
      sha256sum: string,
      tag: string,
    }).isRequired,

  }

  mkFeedback = (scan) => {
    const {
      completed_at,
      reponame,
      sha256sum,
      tag,
    } = scan;
    const mailTo = ['mailto:nautilus-feedback@docker.com?subject=Feedback',
      ` about the ${reponame}:${tag} image&body=${reponame}:${tag}`,
      ` image (${sha256sum}), last scanned at ${completed_at}`].join(' ');
    return (
      <div className={css.text}>
        <a href={mailTo} className={css.text}>Provide Feedback</a>
      </div>
    );
  }

  mkLastScanned = (scan) => {
    const { completed_at, latest_scan_status } = scan;
    const time = moment(completed_at).fromNow();
    let text;
    if (latest_scan_status === IN_PROGRESS) {
      text = `(New scan in progress, showing results from ${time})`;
    } else {
      text = `(Last scanned ${time})`;
    }
    return <div className={css.text}>{text}</div>;
  }

  mkVulnerabilityText = (componentsBySeverity) => {
    const { critical, major, minor, secure } = componentsBySeverity;
    const numTotalComponents = critical.length + major.length
      + minor.length + secure.length;
    const numVulnerableComponents = numTotalComponents - secure.length;
    const components = numVulnerableComponents === 1
      ? 'component' : 'components';
    const isOrAre = numVulnerableComponents === 1 ? 'is' : 'are';
    if (!numVulnerableComponents) {
      return (
        <div className={`${css.text} ${css.clean}`}>
          <CheckIcon size={SMALL} className={css.icon} />
          <span className={css.bold}>{'Your image is clean!'}</span>
          {' No known vulnerabilities were found.'}
        </div>
      );
    }
    return (
      <div className={css.text}>
        {`There ${isOrAre} `}
        <span className={css.bold}>{`${numVulnerableComponents} `}</span>
        {`vulnerable ${components}`}
      </div>
    );
  }

  render() {
    const { componentsBySeverity, scan } = this.props;
    return (
      <div className={css.row}>
        {this.mkVulnerabilityText(componentsBySeverity)}
        {this.mkLastScanned(scan)}
        {this.mkFeedback(scan)}
      </div>
    );
  }
}
