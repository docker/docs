'use strict';

import React, { PropTypes, Component } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import { Link } from 'react-router';
import Card, { Block } from '@dux/element-card';
import Button from '@dux/element-button';
import { PageHeader } from 'dux';
import FA from 'common/FontAwesome';
import styles from './EnterpriseTrialSuccess.css';
import StepsTab from './EnterpriseTrialSuccess/StepsTab';
import InstallCSEngine from './EnterpriseTrialSuccess/InstallCSEngine';
import InstallDTR from './EnterpriseTrialSuccess/InstallDTR';
import AddLicense from './EnterpriseTrialSuccess/AddLicense';
import RouteNotFound404Page from 'common/RouteNotFound404Page';
import SecureDTR from './EnterpriseTrialSuccess/SecureDTR';
import EnterpriseTrialSuccessStore from 'stores/EnterpriseTrialSuccessStore';
import _parseInt from 'lodash/string/parseInt';
import { DEFAULT,
         ERROR } from 'stores/enterprisetrialsuccessstore/Constants';
const { func, string, object } = PropTypes;

class EnterpriseTrialSuccess extends Component {

  static propTypes = {
    JWT: string.isRequired,
    license: object.isRequired,
    STATUS: string.isRequired
  }

  renderTab = (currentStep) => {
    const { JWT, license } = this.props;
    switch (currentStep) {
      case 1:
        return <InstallCSEngine />;
      case 2:
        return <InstallDTR />;
      case 3:
        return <AddLicense JWT={JWT} license={license} />;
      case 4:
        return <SecureDTR />;
      default:
        return <InstallCSEngine />;
    }
  }

  render() {
    const { STATUS } = this.props;
    if (STATUS === ERROR) {
      //trial license for this namespace cannot be found
      return <RouteNotFound404Page />;
    }
    const { namespace,
            step } = this.props.location.query;
    const currentStep = step ? _parseInt(step) : 1;
    const instructions = this.renderTab(currentStep);
    let nextButton;
    if (currentStep === 4) {
      nextButton = null;
    } else {
      const query = { namespace,
                      step: currentStep + 1 };
      nextButton = (
        <Link to='/enterprise/trial/success/' query={query}>
          <Button>
            Next
          </Button>
        </Link>
      );
    }
    return (
      <div>
        <PageHeader title={`Success! Your trial subscription for ${ namespace } is ready.`} />
        <div className={styles.pageWrapper}>
          <div className='row'>
            <div className='columns large-12'>
              <Card>
                <Block>
                  <h3>Installation Steps</h3>
                  <div className={styles.title}>
                    Follow the steps below to install Docker Datacenter components.
                  </div>
                  <div className={styles.tabsWrapper}>
                    <StepsTab step={1}
                              title='Install CS Engine'
                              namespace={namespace}
                              currentStep={currentStep} />
                    <StepsTab step={2}
                              title='Install DTR + UCP'
                              namespace={namespace}
                              currentStep={currentStep} />
                    <StepsTab step={3}
                              title='Add License'
                              namespace={namespace}
                              currentStep={currentStep} />
                    <StepsTab step={4}
                              title='Secure and Configure'
                              namespace={namespace}
                              currentStep={currentStep} />
                  </div>
                  <div className={styles.instructions}>
                    { instructions }
                  </div>
                  <hr />
                  { nextButton }
                </Block>
              </Card>
            </div>
          </div>
        </div>
      </div>
    );

	}
}

export default connectToStores(EnterpriseTrialSuccess,
  [ EnterpriseTrialSuccessStore ],
  function({ getStore }, props) {
    return getStore(EnterpriseTrialSuccessStore).getState();
  });
