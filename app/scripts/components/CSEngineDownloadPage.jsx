'use strict';

import React, { PropTypes, Component } from 'react';
import CSEngineBox from 'common/CSEngineBox';
import styles from './CSEngineDownloadPage.css';

//publicly accessible page for only CS engine download
export default class CSEngineDownloadPage extends Component {
	render() {
		return (
      <div className={styles.pageWrapper}>
        <div className='row'>
          <div className='columns large-12'>
            <CSEngineBox />
          </div>
        </div>
      </div>
    );
	}
}
