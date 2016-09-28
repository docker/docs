'use strict';

import React, { PropTypes, Component } from 'react';
import CopyCodeBox from 'common/CopyCodeBox';
import { DEB, RPM } from 'common/data/csEngineInstructions';
import styles from './CSEngineBox.css';

export default class CSEngineBox extends Component {
	render() {
		return (
      <div className={styles.downloadFlexItem}>
        <div className={styles.curlWrap}>
          <h3>Install CS Engine</h3>
          <div className='row'>
            <div className={'columns large-12 ' + styles.curlHelp }>
              * Copy and run either the RPM or the DEB specific instructions in your terminal.
            </div>
          </div>
          <div className='row'>
            <div className='columns large-12'>
              RPM
              <CopyCodeBox content={RPM}
                           dollar={true}
                           lines={4} />
            </div>
          </div>
          <div className='row'>
            <div className='columns large-12'>
              DEB
              <CopyCodeBox content={DEB}
                           dollar={true}
                           lines={5} />

                <span className={styles.curlHelp}>
                  For more details about this installation,&nbsp;
                  <a href='http://docs.docker.com/docker-trusted-registry/install/install-csengine/'
                     target='_blank'>
                    view our documentation
                  </a>
                </span>
            </div>
          </div>
        </div>
      </div>
    );
	}
}
