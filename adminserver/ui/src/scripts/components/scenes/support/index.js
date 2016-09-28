'use strict';

import React, { Component } from 'react';
import { connect } from 'react-redux';
import * as SettingsActions from 'actions/settings';

import PageHeader from 'components/common/pageHeader';
import CommonStyles from 'components/common/styles.css';
import DescriptParagraph from 'components/common/descriptParagraph';
import styles from './styles.css';
import { mapActions } from 'utils';
import autoaction from 'autoaction';

const mapState = (state) => ({ settings: state.settings });

@connect(mapState, mapActions(SettingsActions))
@autoaction({
  // We need to ensure that we load the user's license to show the correct
  // hourly support component, if necessary.
  getSettings: []
}, SettingsActions)
export default class Support extends Component {

  static propTypes = {
    actions: React.PropTypes.object,
    settings: React.PropTypes.object
  }

  static defaultProps = {
    settings: {
      license: {}
    }
  }

  render() {
    let license = this.props.settings.license || {};
    if (license.type === 'Hourly') {
      return <HourlySupport />;
    }
    return <StandardSupport />;
  }
}

class HourlySupport extends Component {
  render() {
    return (
      <div>
        <PageHeader title='Support' />
        <p className={ styles.hourlySupport }>For support please email <a href='mailto:support@docker.com'>support@docker.com</a>.</p>
      </div>
    );
  }
}

class StandardSupport extends Component {
  render() {
    return (
      <div>
        <PageHeader title='Support' />
        <DescriptParagraph>Technical experts are standing by to answer your questions, solve questions and issue patches and updates to keep your Docker environment secure and up to date. Docker offers commercial support subscription for selected product releases with a variety of support plans to align to your application SLA's.</DescriptParagraph>

        <div className={ styles.container }>
          <div className={ styles.text }>
            <h4>What is commercial support?</h4>
            <p>Commercial support subscription refers to "Supported Software" meaning the Docker or third party software identified on the Order Form as software for which Docker or its authorized resellers agree to provide Subscription Services to Customer. Read the Subscription Services terms.</p>

            <h4>How do I purchase support?</h4>
            <p>Support for Docker products are included as part of the Docker Subscription bundles. To upgrade your support contact sales via <a href='http://goto.docker.com/sales-inquiry.html' target='_blank'>this off-site link</a>.</p>

            <h4>Can I mix and match support?</h4>
            <p>Yes. You will need to keep the Docker components with different subscription tiers on separate infrastructure.</p>
          </div>

          <div className={ styles.tableContainer }>
            <h4 className={ CommonStyles.centered }>Commercial Support Service Levels</h4>
            <table className={ styles.table }>
              <thead>
                <tr>
                  <td></td>
                  <td>Email Support</td>
                  <td>Business Day Support</td>
                  <td>Business Critical Support</td>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Hours</td>
                  <td>N/A</td>
                  <td>6am - 6pm Mon-Fri</td>
                  <td>24x7x365</td>
                </tr>
                <tr>
                  <td>Email</td>
                  <td>Yes</td>
                  <td>Yes</td>
                  <td>Yes</td>
                </tr>
                <tr>
                  <td>Phone</td>
                  <td></td>
                  <td>Yes</td>
                  <td>Yes</td>
                </tr>
                <tr>
                  <td>Support Portal</td>
                  <td>Yes</td>
                  <td>Yes</td>
                  <td>Yes</td>
                </tr>
                <tr>
                  <td>Contacts</td>
                  <td>1</td>
                  <td>4</td>
                  <td>8</td>
                </tr>
                <tr>
                  <td>Response Times</td>
                  <td>Best Effort</td>
                  <td>1 Day</td>
                  <td>2 Hour</td>
                </tr>
                <tr>
                  <td>Availability</td>
                  <td><a href='http://goto.docker.com/sales-inquiry.html' target='_blank'>Get Pricing</a></td>
                  <td><a href='http://goto.docker.com/sales-inquiry.html' target='_blank'>Get Pricing</a></td>
                  <td><a href='http://goto.docker.com/sales-inquiry.html' target='_blank'>Get Pricing</a></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    );
  }
}
