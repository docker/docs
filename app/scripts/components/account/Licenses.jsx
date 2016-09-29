'use strict';

import React, { PropTypes, createClass, Component } from 'react';
import _ from 'lodash';
import moment from 'moment';
import classnames from 'classnames';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import request from 'superagent';

// Blob is a polyfill
require('vendor/Blob');
import { saveAs } from 'vendor/FileSaver';
import FA from 'common/FontAwesome';

import { PageHeader } from 'dux';
import AccountSettingsLicensesStore from '../../stores/AccountSettingsLicensesStore';
import Row from '../common/Row';
import CSEngineBox from 'common/CSEngineBox';

import styles from './Licenses.css';
import DocumentTitle from 'react-document-title';

var debug = require('debug')('AccountSettingsLicenses');

class LicenseBox extends Component {

  state = {
    hasAcceptedTerms: false,
    hasError: false,
    isDownloading: false,
    // Paid licenses require terms
    // https://github.com/docker/dhe-license-server/blob/master/tiers/tiers.go
    requiresAcceptedTerms: this.props.tier !== 'Trial' && this.props.tier !== 'Evaluation'
  }

  onCheckboxChange = () => {
    this.setState({
      hasAcceptedTerms: !this.state.hasAcceptedTerms
    });
  }

  onClick = (e) => {
    this.setState({
      isDownloading: true,
      hasError: false
    });
    request.get(process.env.REGISTRY_API_BASE_URL
              + '/api/licensing/v3/license/'
              + this.props.orgname
              + '/'
              + this.props.keyId
              + '/')
      .set('Authorization', 'JWT ' + this.props.JWT)
           .end((err, res) => {
             if(err) {
               this.setState({
                 isDownloading: false,
                 hasError: true
               });
             } else {

               this.setState({
                 isDownloading: false,
                 hasError: false
               });

               const blob = new Blob([res.text], {
                 type: 'text/plain;charset=utf-8'
               });
               saveAs(blob, `docker_subscription.lic`);
             }

           });
  }

  render() {
    const { expiration, alias, orgname, end, tier, maxEngines } = this.props;
    const { requiresAcceptedTerms, hasAcceptedTerms } = this.state;
    const exp = moment(expiration).fromNow();
    // This is a temporary work around until we can enable click to accept
    // in Store
    let maybeTermsAndConditions = <div className={styles.legal}></div>;
    if (requiresAcceptedTerms) {
      const EUSALink = 'https://www.docker.com/docker-software-end-user-subscription-agreement';
      const checkbox = (
        <input
          checked={hasAcceptedTerms}
          className={styles.checkbox}
          onChange={this.onCheckboxChange}
          ref='eusaCheckbox'
          type='checkbox'
        />
      );
      maybeTermsAndConditions = (
        <div className={styles.legal}>
          { checkbox }
          I agree to Docker's subscription <a href={EUSALink} target="_blank">terms</a>
        <div className={styles.finePrint}>This will not override any pre-negotiated terms.</div>
        </div>
      );
    }

    const preventDownload = requiresAcceptedTerms && !hasAcceptedTerms;
    let icon;
    let onClick;
    if (preventDownload) {
      icon = 'Please accept terms to download';
      onClick = () => {};
    } else {
      icon = <FA icon='fa-cloud-download'/>;
      onClick = this.onClick;
      if(this.state.isDownloading) {
        icon = <FA icon='fa-circle-o-notch' animate='spin'/>;
      }
    }

    const classes = classnames({
      'large-3 columns': true,
      'end': end
    });

    const downloadClasses = classnames({
      [styles.download]: !this.state.hasError,
      [styles.downloadAlert]: this.state.hasError
    });
    return (
      <div className={classes}>
        <div className={styles.license}>
            <h3>{alias}</h3>
            <p>Organization: {orgname}</p>
            <p>Engines: {maxEngines}</p>
            <p><FA icon='fa-clock'/>Expires: <span>{exp}</span></p>
            { maybeTermsAndConditions }
        </div>
        <a onClick={onClick}>
          <div className={downloadClasses}>{icon}</div></a>
      </div>
    );
  }
}

class Licenses extends Component {
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  }

  propTypes: {
    JWT: PropTypes.string.isRequired,
    licenses: PropTypes.array.isRequired
  }

  state = {
    rpm: {
      isDownloading: false,
      hasError: false
    },
    deb: {
      isDownloading: false,
      hasError: false
    }
  }


  render() {
    if(this.props.licenses.length <= 0) {
      return (
        <DocumentTitle title='Licenses - Docker Hub'>
          <div>
            <PageHeader title='Licenses'></PageHeader>
            <div className={styles.pageWrapper}>
              <Row>
                <div className='large-12 columns module'>You don't seem to have any
                  licenses; Get a
                  <a href='https://store.docker.com/bundles/docker-datacenter/purchase?plan=free-trial' target='_blank'> Trial.</a>
                </div>
              </Row>
            </div>
          </div>
        </DocumentTitle>
      );
    } else {
      const { rpm, deb } = this.state;

      const icon = <FA icon='fa-cloud-download' />;
      const iconDownloading = <FA icon='fa-circle-o-notch' animate='spin'/>;

      const rpmIcon = rpm.isDownloading ? iconDownloading : icon;
      const debIcon = deb.isDownloading ? iconDownloading : icon;

      const rpmIntent = rpm.hasAlert ? 'alert' : null;
      const debIntent = deb.hasAlert ? 'alert' : null;

      return (
        <DocumentTitle title='Licenses - Docker Hub'>
          <div>
            <PageHeader title='Licenses' />
            <div className={styles.pageWrapper}>
              <div className='row'>
                <div className='large-12 columns'>
                  <CSEngineBox />
                </div>
              </div>
              <div className='row'>
                {this.props.licenses.map((license, i, arr) => { return <LicenseBox {...license} key={license.keyId} JWT={this.props.JWT} end={i === arr.length - 1}/>; } )}
              </div>
            </div>
          </div>
        </DocumentTitle>
      );
    }
  }
}

export default connectToStores(Licenses,
                               [
                                 AccountSettingsLicensesStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(AccountSettingsLicensesStore).getState();
                               });
