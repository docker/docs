'use strict';

import React, { Component, PropTypes } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
const debug = require('debug')('EnterpriseTrialForm');
import Card, { Block } from '@dux/element-card';
import FA from '../common/FontAwesome.jsx';
import EnterpriseTrialFormStore from 'stores/EnterpriseTrialFormStore';
import updateEnterpriseTrialFormField from 'actions/updateEnterpriseTrialFormField';
import createNewLicenseTrial from 'actions/createNewLicenseTrial';
import {
  ATTEMPTING,
  DEFAULT,
  FACEPALM,
  SUCCESSFUL_SIGNUP
} from 'stores/enterprisetrialstore/Constants';
import SimpleInput from 'components/common/SimpleInput';
import onChange from 'actions/common/onChangeUtil';
import Button from '@dux/element-button';
import getValues from 'lodash/object/values';
import includes from 'lodash/collection/includes';
import omit from 'lodash/object/omit';
import map from 'lodash/collection/map';
import countries from 'common/data/countries.js';
import provinces from 'common/data/provinces.js';
import states from 'common/data/states.js';
import jobFunctions from 'common/data/jobFunctions.js';
import styles from './EnterpriseTrialForm.css';
import classnames from 'classnames';

const {
  array,
  func,
  object,
  shape,
  string
} = PropTypes;

// the set of countries that have states
const UNITED_STATES = 'United States';
const CANADA = 'Canada';

const COUNTRIES_WITH_STATES = [UNITED_STATES, CANADA];

class EnterpriseTrialForm extends Component {

  static propTypes = {
    JWT: string.isRequired,
    enterpriseTrial: shape({
      fields: object,
      globalFormError: string.isRequired,
      orgs: array.isRequired,
      STATUS: string.isRequired,
      values: object
    }).isRequired
  }

  static contextTypes = {
    executeAction: func.isRequired
  }

  _mkCountryOptions = (country) => {
    return (
      <option key={country.code} value={country.name}>{country.name}</option>
    );
  }

  _mkOrgOptions = (orgname) => {
    return (
      <option key={orgname} value={orgname}>{orgname}</option>
    );
  }

  _mkJobOptions = (jobTitle) => {
    return (
      <option key={jobTitle} value={jobTitle}>{jobTitle}</option>
    );
  }

  _mkStateOptions = (state) => {
    return (
      <option key={state.abbr} value={state.abbr}>{state.name}</option>
    );
  }

  _onChange = onChange({
    storePrefix: 'ENTERPRISE_TRIAL'
  })

  _onSelectChange = (fieldKey) => {
    return (e) => {
      const fieldValue = e.target.options[e.target.selectedIndex].value;
      this.context.executeAction(updateEnterpriseTrialFormField, {
        fieldKey,
        fieldValue: fieldValue
      });
      //if country is changed, clear the `state` selection
      if (fieldKey === 'country') {
        this.context.executeAction(updateEnterpriseTrialFormField, {
          fieldKey: 'state',
          fieldValue: ''
        });
      }
    };
  }

  onSubmit = (e) => {
    e.preventDefault();
    const { JWT } = this.props;
    const {
      companyName,
      country,
      state,
      email,
      firstName,
      lastName,
      jobFunction,
      namespace,
      phoneNumber
    } = this.props.enterpriseTrial.values;
    const packageName = 'Trial';
    // TODO: is this true:
    //namespace must be an organization
    this.context.executeAction(createNewLicenseTrial, {
      JWT,
      companyName,
      country,
      state,
      email,
      firstName,
      lastName,
      jobFunction,
      namespace,
      packageName,
      phoneNumber
    });
  }

  isStateRequired = (country) => includes(COUNTRIES_WITH_STATES, country);

  isFormComplete = () => {
    const { values } = this.props.enterpriseTrial;
    // if country !== US or Canada, state is not required
    const requiredFields = this.isStateRequired(values.country) ? values : omit(values, 'state');
    return !includes(requiredFields, '');
  }

  render() {
    const {
      fields,
      globalFormError,
      orgs,
      STATUS,
      values
    } = this.props.enterpriseTrial;
    const disabled = !this.isFormComplete();
    let disabledText = null;
    if (disabled) {
      disabledText = <div className={styles.label}>All fields are required.</div>;
    }
    let globalError = null;
    if(globalFormError) {
      const extraText = (
        <div>
          <a href='https://goto.docker.com/contact-us.html'>Contact us</a>
          {' to extend your trial and discuss purchasing options'}
        </div>
      );
      globalError = (
        <div className={styles.error}>
          {globalFormError}
          { STATUS !== FACEPALM ? extraText : '' }
        </div>
      );
    }
    let buttonText;
    let spinner = <FA icon='fa-spinner fa-spin' />;
    if (STATUS === ATTEMPTING) {
      buttonText = 'Registering ';
    } else if (STATUS === SUCCESSFUL_SIGNUP) {
      buttonText = 'Redirecting ';
    } else {
      buttonText = 'Start Your Free Trial';
      spinner = '';
    }
    const buttonClasses = classnames({
      [styles.fullWidth]: true,
      [styles.disabled]: disabled
    });

    const button = (
      <div className={buttonClasses}>
        <Button type='submit'
                variant={ STATUS === SUCCESSFUL_SIGNUP ? 'success' : 'primary'}
                disabled={disabled}>
          { buttonText } { spinner }
        </Button>
      </div>
    );

    // Do not show state option unless country is US or Canada
    let stateArea;
    if (includes(COUNTRIES_WITH_STATES, values.country)) {
      const statesOrProvinces = values.country === UNITED_STATES ? states : provinces;
      stateArea = (
        <div>
          <label className={styles.label}>State/Province*</label>
          <div className={styles.error}>
            { fields.state.hasError ? fields.state.error : '' }
          </div>
          <select className={styles.select}
            ref='state'
            value={values.state}
            onChange={this._onSelectChange('state')}>
            {map(statesOrProvinces, this._mkStateOptions)}
          </select>
        </div>
      );
    }
    return (
      <Card>
        <Block>
          <div className={styles.formWrapper}>
            <p>
              This trial gives you access to the tools you need to manage
             your Docker images, containers, and applications within your firewall.
            </p>
            <form onSubmit={this.onSubmit}>
              <div className={styles.title}>Attach License to Docker Hub Account</div>
              <div className={styles.error}>
                { fields.namespace.hasError ? fields.namespace.error : '' }
              </div>
              <select className={styles.select}
                      ref='orgSelect'
                      autoFocus
                      value={values.namespace}
                      onChange={this._onSelectChange('namespace')}>
                {map(orgs, this._mkOrgOptions)}
              </select>

              <div className={styles.title}>Contact Information</div>
              <label className={styles.label}>First Name*</label>
              <div className={styles.error}>
                { fields.firstName.hasError ? fields.firstName.error : '' }
              </div>
              <SimpleInput onChange={this._onChange('firstName').bind(this)}
                           hasError={fields.firstName.hasError}
                           value={values.firstName} />

              <label className={styles.label}>Last Name*</label>
              <div className={styles.error}>
                { fields.lastName.hasError ? fields.lastName.error : '' }
              </div>
              <SimpleInput onChange={this._onChange('lastName').bind(this)}
                           hasError={fields.lastName.hasError}
                           value={values.lastName} />

              <label className={styles.label}>Company*</label>
              <div className={styles.error}>
                { fields.companyName.hasError ? fields.companyName.error : '' }
              </div>
              <SimpleInput onChange={this._onChange('companyName').bind(this)}
                           hasError={fields.companyName.hasError}
                           value={values.companyName} />

              <label className={styles.label}>Job Function*</label>
              <div className={styles.error}>
                { fields.jobFunction.hasError ? fields.jobFunction.error : '' }
              </div>
              <select className={styles.select}
                      ref='jobFunction'
                      value={values.jobFunction}
                      onChange={this._onSelectChange('jobFunction')}>
                {map(jobFunctions, this._mkJobOptions)}
              </select>

              <label className={styles.label}>Company Email Address*</label>
              <div className={styles.error}>
                { fields.email.hasError ? fields.email.error : '' }
              </div>
              <SimpleInput onChange={this._onChange('email').bind(this)}
                           type='email'
                           hasError={fields.email.hasError}
                           value={values.email} />

              <label className={styles.label}>Phone Number*</label>
              <div className={styles.error}>
                { fields.phoneNumber.hasError ? fields.phoneNumber.error : '' }
              </div>
              <SimpleInput onChange={this._onChange('phoneNumber').bind(this)}
                           hasError={fields.phoneNumber.hasError}
                           value={values.phoneNumber} />

              <label className={styles.label}>Country*</label>
              <div className={styles.error}>
                { fields.country.hasError ? fields.country.error : '' }
              </div>
              <select className={styles.select}
                      ref='countrySelect'
                      value={values.country}
                      onChange={this._onSelectChange('country')}>
                {map(countries, this._mkCountryOptions)}
              </select>
              {stateArea}

              {globalError}
              {disabledText}
              {button}
              <div className={styles.termsText}>
                By registering for the free trial, you agree to Docker's
                <Link to='/enterprise/trial/terms/'> software evaluation terms.</Link>
              </div>
            </form>
          </div>
        </Block>
      </Card>
    );
  }
}

export default connectToStores(EnterpriseTrialForm,
  [
    EnterpriseTrialFormStore
  ],
  function({ getStore }, props) {
    return {
      enterpriseTrial: getStore(EnterpriseTrialFormStore).getState()
    };
  });
