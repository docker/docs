'use strict';

import React from 'react';
import styles from './ScanHeader.css';
import capitalize from 'lodash/string/capitalize';
import size from 'lodash/collection/size';
import FontAwesome from 'common/FontAwesome';
import moment from 'moment';
import { consts } from '../nautilusUtils';

const ScanHeader = ({ scan, componentsBySeverity }) => {
  const { critical, major, minor, secure } = componentsBySeverity;
  const { reponame, tag, sha256sum, completed_at, latest_scan_status } = scan;
  const numTotalComponents = critical.length + major.length + minor.length + secure.length;
  const numVulnerableComponents = numTotalComponents - secure.length;
  const components = numTotalComponents === 1 ? `component` : `components`;
  const isOrAre = numVulnerableComponents === 1 ? `is` : `are`;
  const lastScanned = latest_scan_status === consts.IN_PROGRESS ? `New scan in progress, showing results from ${moment(completed_at).fromNow()}` : `Scanned ${moment(completed_at).fromNow()}`;
  let vulnText, vulnArea;
  if (!numVulnerableComponents) {
    vulnText = <div><b>Your image is clean!</b>&nbsp;No known vulnerabilities were found.</div>;
    vulnArea = (
      <div>
        <div className={styles.bigCheck}>
          <FontAwesome icon='fa-check' size='3x' />
        </div>
        <div className={styles.inlineBlock}>
          <div className={styles.vulnerabilityText}>{vulnText}</div>
          <div className={styles.lastScanned}>{lastScanned}</div>
        </div>
      </div>
    );
  } else {
    vulnText = `${numVulnerableComponents} of ${numTotalComponents} ${components} ${isOrAre} vulnerable`;
    vulnArea = (
      <div>
        <div className={styles.vulnerabilityText}>{vulnText}</div>
        <div className={styles.lastScanned}>{lastScanned}</div>
      </div>
    );
  }
  const feedback = (
    <a href={`mailto:nautilus-feedback@docker.com?subject=Feedback about the ${reponame}:${tag} image&body=${reponame}:${tag} image (${sha256sum}), last scanned at ${completed_at}`}
       className={styles.feedbackLink}>
      Provide Feedback
    </a>
  );
  return (
    <div>
      <div className={'row ' + styles.headerRow}>
        <div className='columns large-9'>
          {vulnArea}
        </div>
        <div className={`columns large-3 ${styles.feedback}`}>{feedback}</div>
      </div>
      <div className={`row ${styles.tableHeader}`}>
        <div className='columns large-5'><b>Layers</b></div>
        <div className='columns large-7'><b className={styles.componentsTitle}>Components</b></div>
      </div>
    </div>
  );
};

export default ScanHeader;
