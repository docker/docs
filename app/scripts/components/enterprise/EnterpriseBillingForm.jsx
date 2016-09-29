'use strict';

import React, {
  PropTypes,
  Component
  } from 'react';
import { findDOMNode } from 'react-dom';
const { string, func, shape, object, array } = PropTypes;
import includes from 'lodash/collection/includes';
import merge from 'lodash/object/merge';
import has from 'lodash/object/has';
import connectToStores from 'fluxible-addons-react/connectToStores';
import classnames from 'classnames';
var debug = require('debug')('EnterprisePaid');

import DUXInput from '../common/DUXInput';
import FA from '../common/FontAwesome.jsx';
import acceptedCards from '../common/data/acceptedCards';
import EnterprisePaidFormStore from 'stores/EnterprisePaidFormStore';
import { Button } from 'dux';
import Card, { Block } from '@dux/element-card';
import validateBillingInfo from '../../actions/common/validateBillingInfo.js';
import onChange from '../../actions/common/onChangeUtil';
import styles from './EnterpriseBillingForm.css';
import EnterprisePanel from './EnterpriseBillingForm/EnterprisePanel.jsx';
import { STATUS as BASE_STATUS } from 'stores/common/Constants';

import months from 'common/data/months.js';
import years from 'common/data/years.js';

class EnterprisePaid extends Component {

  static contextTypes = {
    executeAction: func.isRequired
  }

  static propTypes = {
    user: shape({username: string.isRequired}).isRequired,
    JWT: string.isRequired,
    submitAction: func.isRequired,
    orgAction: func,
    enterpriseType: string.isRequired,
    orgs: array,
    currentPlan: object,
    billingInfo: object,
    accountInfo: object
  }

  state = {
    checkboxError: false
  }

  _onCheckboxChange = (e) => {
    // When the checkbox is clicked, clear any potential errors
    this.setState({
      checkboxError: false
    });
  }

  onSubmit = (e) => {
    e.preventDefault();
    if(!findDOMNode(this.refs.eulaCheckbox).checked) {
      this.setState({
        checkboxError: true
      });
    } else if (has(this.props.currentPlan, 'plan') && this.props.currentPlan.plan) {
      let username = findDOMNode(this.refs.orgSelect).value; //Or orgname
      this.props.submitAction(username, this.props.values);
    } else {
      const validate = window.recurly.validate;
      const { values } = this.props;
      let fieldErrors = {
        number: !validate.cardNumber(values.number),
        expiry: !validate.expiry(values.month, values.year),
        cvv: !validate.cvv(values.cvv),
        first_name: !values.first_name,
        last_name: !values.last_name,
        city: !values.city,
        address1: !values.address1,
        country: !values.country,
        state: !values.state,
        postal_code: !values.postal_code
      };
      fieldErrors.email = this.props.enterpriseType === 'cloud' ? !values.email : false;
      if (includes(fieldErrors, true)) {
        this.context.executeAction(validateBillingInfo({storePrefix: 'ENTERPRISE_PAID'}), fieldErrors);
      } else {
        const namespace = findDOMNode(this.refs.orgSelect).value; //Or orgname
        this.props.submitAction(namespace, values);
      }
    }
  }

  _onChange = onChange({
    storePrefix: 'ENTERPRISE_PAID'
  })

  _onOrgChange = (e) => {
    e.preventDefault();
    if (this.props.orgAction) {
      // Should update the billing info/plans given org/user namespace
      this.props.orgAction(e.target.value);
    }
  }

  _mkIconFromCardType(cardType) {
    const icon = acceptedCards[cardType];
    if (icon) {
      return (<FA icon={`fa-cc-${icon}`} size='2x'/>);
    }
  }

  _renderOptions(list_item){
    return (
      <option key={list_item.abbr} value={list_item.name}>{list_item.name}</option>
    );
  }

  _renderOrgOptions(org) {
    return (
      <option key={org} value={org}>{org}</option>
    );
  }

  render() {
    const {
      accountInfo,
      billingInfo,
      currentPlan,
      fields,
      orgs,
      STATUS,
      values,
      enterpriseType
      } = this.props;

    let formBody = null;
    let contactForm = null;
    if (has(currentPlan, 'plan') && currentPlan.plan) {
      // currentPlan will only be passed in for Cloud
      formBody = (
        <EnterprisePanel heading='Payment Information'>
          <div className={styles.formBody}>
            <div className="row">
              <div className="columns large-4">
                <div>Name:</div>
              </div>
              <div className={'columns large-8 ' + styles.infoContent}>
                <div>{accountInfo.first_name} {accountInfo.last_name}</div>
              </div>
            </div>
            <div className="row">
              <div className="columns large-4">
                <div>Email:</div>
              </div>
              <div className={'columns large-8 ' + styles.infoContent}>
                <div>{accountInfo.email}</div>
              </div>
            </div>
            <div className="row">
              <div className="columns large-4">
                <div>Company:</div>
              </div>
              <div className={'columns large-8 ' + styles.infoContent}>
                <div>{accountInfo.company_name}</div>
              </div>
            </div>
            <div className="row">
              <div className="columns large-4">Card Info:</div>
              <div className={'columns large-8 ' + styles.infoContent}>
                {billingInfo.first_name} {billingInfo.last_name}<br/>
                {billingInfo.card_type} card ending with x{billingInfo.last_four} <br/>
                Expiration: {billingInfo.month}/{billingInfo.year}
              </div>
            </div>
            <div className="row">
              <div className="columns large-4">Billing Address:</div>
              <div className={'columns large-8 ' + styles.infoContent}>
                {billingInfo.address1} <br/>
                {billingInfo.address2} <br/>
                {billingInfo.city} {billingInfo.state} {billingInfo.zip} <br/>
                {billingInfo.country}
              </div>
            </div>
          </div>
        </EnterprisePanel>
      );
    } else {
      // No Plan = Signing up for Server -OR- NEW Cloud plan
      const monthClass = classnames({
        [styles.dropdownError]: fields.expiry.hasError,
        [styles.dropdown]: true
      });
      const yearClass = classnames({
        [styles.dropdownError]: fields.expiry.hasError,
        [styles.dropdown]: true
      });
      const cardIcon = values.number ? this._mkIconFromCardType(values.card_type) : null;
      formBody = (
        <EnterprisePanel heading='Payment Information'>
          <div className={styles.formBody}>
            <div className='row'>
              <div className='large-6 columns'>
                <DUXInput label='First Name'
                          onChange={this._onChange('first_name').bind(this)}
                          value={values.first_name}
                  {...fields.first_name} />
              </div>
              <div className='large-6 columns end'>
                <DUXInput label='Last Name'
                          onChange={this._onChange('last_name').bind(this)}
                          value={values.last_name}
                  {...fields.last_name} />
              </div>
            </div>
            <div className='row'>
              <div className='large-6 columns'>
                <DUXInput label={'Credit Card Number'}
                          onChange={this._onChange('number').bind(this)}
                          value={values.number}
                  {...fields.number} />
              </div>
              <div className='large-2 columns'>
                <DUXInput label='CVV'
                          onChange={this._onChange('cvv').bind(this)}
                          value={values.cvv}
                  {...fields.cvv} />
              </div>
              <div className='large-2 columns end'>
                {cardIcon}
              </div>
            </div>
            <div className='row'>
              <div className='large-2 columns'>
                Expiration
              </div>
              <div className='large-2 columns'>
                <select className={monthClass}
                        onChange={this._onChange('month').bind(this)}
                        value={values.month}>
                  {months.map(this._renderOptions, this)}
                </select>
              </div>
              <div className='large-2 columns end'>
                <select className={yearClass}
                        onChange={this._onChange('year').bind(this)}
                        value={values.year}>
                  {years.map(this._renderOptions, this)}
                </select>
              </div>
            </div>
            <div className='row'>
              <div className='large-6 columns'>
                <DUXInput label='Billing Address'
                          onChange={this._onChange('address1').bind(this)}
                          value={values.address1}
                  {...fields.address1} />
              </div>
              <div className='large-6 columns end'>
                <DUXInput label='City'
                          onChange={this._onChange('city').bind(this)}
                          value={values.city}
                  {...fields.city}/>
              </div>
            </div>
            <div className='row'>
              <div className='large-6 columns'>
                <DUXInput label='State'
                          onChange={this._onChange('state').bind(this)}
                          value={values.state}
                  {...fields.state} />
              </div>
              <div className='large-6 columns end'>
                <DUXInput label='Postal Code'
                          onChange={this._onChange('postal_code').bind(this)}
                          value={values.postal_code}
                  {...fields.postal_code}/>
              </div>
            </div>
            <div className='row'>
              <div className='large-6 columns end'>
                <p>Country</p>
                <select onChange={this._onChange('country').bind(this)}
                        value={values.country}>
                  <option value='US'>United States of America</option>
                  <option value='CA'>Canada</option>
                </select>
              </div>
            </div>
          </div>
        </EnterprisePanel>
      );
      if (accountInfo) {
        contactForm = (
          <EnterprisePanel heading='Contact Information'>
            <div className={styles.formBody}>
              <div className="row">
                <div className="columns large-6">
                  <DUXInput label="First Name"
                            onChange={this._onChange('account_first').bind(this)}
                            value={values.account_first}/>
                </div>
                <div className="columns large-6">
                  <DUXInput label="Last Name"
                            onChange={this._onChange('account_last').bind(this)}
                            value={values.account_last}/>
                </div>
              </div>

              <div className="row">
                <div className='columns large-8 end'>
                  <DUXInput label="Company"
                            onChange={this._onChange('company_name').bind(this)}
                            value={values.company_name}/>
                </div>
              </div>

              <div className="row">
                <div className='columns large-8 end'>
                  <DUXInput label="Email"
                            type="email"
                            onChange={this._onChange('email').bind(this)}
                            value={values.email}
                    {...fields.email} />
                </div>
              </div>
            </div>
          </EnterprisePanel>
        );
      }
    }

    let checkboxError = null;
    if (this.state.checkboxError) {
      checkboxError = (
        <p className='alert-box alert'>The Terms must be accepted to continue</p>
      );
    }

    let submitError = null;
    if (STATUS === BASE_STATUS.FACEPALM) {
      debug('FACEPALLLLM');
      submitError = <span className={styles.error}>There was an error submitting the form. Please <a
        href='https://support.docker.com/'>contact support</a>.</span>;
    } else if (this.props.globalFormError) {
      submitError = <span className={styles.error}>{this.props.globalFormError}</span>;
    }

    let submitButton = (
      <Button type='submit'>Submit</Button>
    );
    if (STATUS === BASE_STATUS.ATTEMPTING) {
      debug('ATTEMPTING');
      submitButton = (
        <Button disabled type='submit'>
          Submitting <FA icon='fa-spinner fa-spin'/>
        </Button>
      );
    }
    let intent;
    if (STATUS === BASE_STATUS.SUCCESSFUL) {
      intent = 'success';
      submitButton = (
        <Button intent={intent} disabled type='submit'>
          Redirecting <FA icon='fa-spinner fa-spin'/>
        </Button>
      );
    } else if (STATUS === BASE_STATUS.FACEPALM || STATUS === BASE_STATUS.ERROR) {
      submitButton = (
        <Button intent={'alert'} type='submit'>Submit</Button>
      );
    }
    const attach = enterpriseType === 'cloud' ? 'package' : 'license';
    return (
      <div>
        <Card>
          <Block>
            <EnterprisePanel heading={'Attach this ' + attach + ' to a Docker Hub Organization'}>
              <div className={styles.formBody}>
                <select ref='orgSelect' onChange={this._onOrgChange}>
                  {orgs.map(this._renderOrgOptions, this)}
                </select>
              </div>
            </EnterprisePanel>
            <hr/>
            {contactForm}
            {formBody}
            <hr/>
            <EnterprisePanel heading='End User Software Agreement'>
              {checkboxError}
              <input type='checkbox'
                     ref='eulaCheckbox'
                     onChange={this._onCheckboxChange}/>
              <span> I have read and agree to the <a href='/enterprise/eusa/'>Docker End User Software
                Agreement</a></span>
            </EnterprisePanel>
          </Block>
        </Card>

        <form className={styles.formSubmit} onSubmit={this.onSubmit}>
          {submitButton} {submitError}
        </form>
      </div>
    );
  }
}

export default connectToStores(EnterprisePaid,
  [
    EnterprisePaidFormStore
  ],
  function({ getStore }, props) {
    return getStore(EnterprisePaidFormStore).getState();
  });
