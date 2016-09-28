import React, { Component } from 'react';
import css from './styles.css';

/* eslint-disable max-len */
export default class Footer extends Component {
  render() {
    return (
      <div className={css.footerWrapper}>
        <div className={css.contentWrapper}>
          <div>
            <div className={css.description}>
              Build, Ship, Run. An open platform for distributed applications for developers and sysadmins
            </div>
            <div>Copyright &copy; 2016 Docker Inc. All rights reserved.</div>
          </div>
          <div className={css.links}>
            <a href="https://hub.docker.com" target="_blank">Hub</a>
            <a href="https://cloud.docker.com" target="_blank">Cloud</a>
            <a href="https://www.docker.com/legal" target="_blank">Legal</a>
            <a href="http://www.docker.com" target="_blank">Home</a>
          </div>
        </div>
      </div>
    );
  }
}
/* eslint-enable */
