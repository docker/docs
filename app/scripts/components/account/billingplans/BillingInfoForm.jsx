'use strict';

import React, { PropTypes } from 'react';
const { string, number, object, bool, func, shape, oneOfType } = PropTypes;
import { Link } from 'react-router';
import includes from 'lodash/collection/includes';
import map from 'lodash/collection/map';
import classnames from 'classnames';

import DUXInput from '../../common/DUXInput.jsx';
import FA from '../../common/FontAwesome.jsx';
import acceptedCards from '../../common/data/acceptedCards.js';
import updateBillingInfoFormField from '../../../actions/updateBillingInfoFormField.js';
import validateBillingInfo from '../../../actions/common/validateBillingInfo.js';
import { STATUS } from 'stores/billingformstore/Constants.js';
import { Button } from 'dux';
import Card, { Block } from '@dux/element-card';

var countries = require('common/data/countries.js');
var states = require('common/data/states.js');
var months = require('common/data/months.js');
var years = require('common/data/years.js');

import styles from './BillingInfoForm.css';
var debug = require('debug')('BillingInfoForm:');

var _mkOptions = function(list_item){
  return (
      <option key={list_item.abbr} value={list_item.abbr}>{list_item.name}</option>
  );
};

var _mkCountryOptions = function(country) {
  return (
    <option key={country.code} value={country.name}>{country.name}</option>
  );
};

var BillingInfoForm = React.createClass({
  contextTypes: {
    getStore: func.isRequired,
    executeAction: func.isRequired
  },
  propTypes: {
    isOrg: bool.isRequired,
    username: string.isRequired,
    accountInfo: shape({
      account_code: string,
      username: string,
      email: string,
      first_name: string,
      last_name: string,
      company_name: string
    }),
    billingInfo: shape({
      city: string,
      state: string,
      zip: string,
      first_name: string,
      last_name: string,
      address1: string,
      address2: string,
      country: string
    }),
    card: shape({
      number: string,
      cvv: string,
      month: oneOfType([number, string]),
      year: oneOfType([number, string]),
      type: string
    }),
    fieldErrors: object.isRequired,
    errorMessage: string,
    STATUS: string.isRequired,
    history: object.isRequired
  },
  _onChange(field, fieldKey) {
    return (e) => {
      this.context.executeAction(updateBillingInfoFormField, {
        field,
        fieldKey,
        fieldValue: e.target.value
      });
    };
  },
  _onSelectChange(field, fieldKey) {
    return (e) => {
      var fieldValue = e.target.options[e.target.selectedIndex].value;
      if (fieldKey === 'month' || fieldKey === 'year') {
        fieldValue = parseInt(fieldValue, 10);
      }
      this.context.executeAction(updateBillingInfoFormField, {
        field,
        fieldKey,
        fieldValue: fieldValue
      });
    };
  },
  _onBackClick(e) {
    e.preventDefault();
    if (this.props.isOrg) {
      this.props.history.pushState(null, `/u/${this.props.username}/dashboard/billing/`);
    } else {
      this.props.history.pushState(null, '/account/billing-plans/');
    }
  },
  _mkIconFromCardType(cardType) {
    const icon = acceptedCards[cardType];
    if (icon) {
      return (<FA icon={`fa-cc-${icon}`} size='2x'/>);
    }
  },
  formSubmit(e){
    e.preventDefault();
    const billing = this.props.billingInfo;
    let validate = window.recurly.validate;
    const fieldErrors = {
      number: !validate.cardNumber(this.props.card.number),
      expiry: !validate.expiry(this.props.card.month, this.props.card.year),
      cvv: !validate.cvv(this.props.card.cvv),
      first_name: !billing.first_name,
      last_name: !billing.last_name,
      city: !billing.city,
      address1: !billing.address1,
      country: !billing.country
    };
    const accountErr = { hasError: !this.props.accountInfo.email };
    if (includes(fieldErrors, true) || accountErr.hasError ) {
      const hasError = { fieldErrors, accountErr };
      this.context.executeAction(validateBillingInfo({storePrefix: 'BILLING'}), hasError);
    } else {
      this.props.submitAction();
    }
  },
  render: function() {
    const { card } = this.props;
    const cardIcon = card.number ? this._mkIconFromCardType(card.type) : null;
    var expiryClass = classnames({
      [styles.error]: this.props.fieldErrors.expiry,
      [styles.billingDropdown]: true
    });
    var countryClass = classnames({
      [styles.billingDropdown]: true,
      [styles.error]: this.props.fieldErrors.country
    });
    let submit = 'Submit';
    if (this.props.STATUS === STATUS.ATTEMPTING) {
      submit = (
        <div>
          Submitting <FA icon='fa-spinner fa-spin'/>
        </div>
      );
    }
    var intent;
    if (this.props.STATUS === STATUS.SUCCESS) {
      intent = 'success';
      submit = (
        <div>
          Redirecting <FA icon='fa-spinner fa-spin'/>
        </div>
      );
    } else if (this.props.STATUS === STATUS.FORM_ERROR) {
      intent = 'alert';
    }
    const submitButton = (
      <Button type="submit"
              intent={ intent }
              disabled={this.props.STATUS === STATUS.ATTEMPTING}
              size='small'>{submit}</Button>
    );
    return (
      <Card>
        <Block>
          <form className={styles.form} onSubmit={this.formSubmit} >
            <div className="row">
              <div className={styles.billingFormSection}>Contact Info:</div>
              <div className="columns large-12">
                <div className="row">
                  <div className="columns large-6">
                    <DUXInput label="First Name"
                              onChange={this._onChange('account', 'first_name')}
                              value={this.props.accountInfo.first_name}/>
                  </div>
                  <div className="columns large-6">
                    <DUXInput label="Last Name"
                              onChange={this._onChange('account', 'last_name')}
                              value={this.props.accountInfo.last_name}/>
                  </div>
                </div>

                <div className="row">
                  <div className='columns large-8 end'>
                    <DUXInput label="Company"
                              onChange={this._onChange('account', 'company_name')}
                              value={this.props.accountInfo.company_name}/>
                  </div>
                </div>

                <div className="row">
                  <div className='columns large-8 end'>
                    <DUXInput label="Email"
                              hasError={this.props.accountInfo.hasError}
                              error='Required'
                              onChange={this._onChange('account', 'email')}
                              value={this.props.accountInfo.email}/>
                  </div>
                </div>

              </div>
            </div>
            <div className="row">
              <div className={styles.marginBottom}>Billing Info:</div>
              <div className="columns large-12">
                <div className="row">
                  <div className="columns large-6">
                    <DUXInput label="First Name"
                            hasError={this.props.fieldErrors.first_name}
                            error='Required'
                            onChange={this._onChange('billing', 'first_name')}
                            value={this.props.billingInfo.first_name}/>
                  </div>
                  <div className="columns large-6">
                    <DUXInput label="Last Name"
                            hasError={this.props.fieldErrors.last_name}
                            error='Required'
                            onChange={this._onChange('billing', 'last_name')}
                            value={this.props.billingInfo.last_name}/>
                  </div>
                </div>
                <div className="row">
                  <div className="columns large-7">
                    <DUXInput label="Credit Card Number"
                            hasError={this.props.fieldErrors.number}
                            error='Invalid Card Value'
                            required='required'
                            onChange={this._onChange('card', 'number')}
                            value={this.props.card.number}/>
                  </div>
                  <div className="columns large-3">
                    <DUXInput label="cvv"
                            hasError={this.props.fieldErrors.cvv}
                            error='Invalid cvv'
                            onChange={this._onChange('card', 'cvv')}
                            value={this.props.card.cvv}/>
                  </div>
                  <div className={'columns large-2 ' + styles.cardIcon}>
                    {cardIcon}
                  </div>
                </div>
                <div className="row">
                  <div className={'columns large-2 ' + styles.dateText}>
                    Expires
                  </div>
                  <div className="columns large-2">
                    <select className={expiryClass}
                            value={this.props.card.month}
                            onChange={this._onSelectChange('card', 'month')}>
                      {months.map(_mkOptions, this)}
                    </select>
                  </div>
                  <div className={'columns large-1 ' + styles.dateText}>
                    /
                  </div>
                  <div className="columns large-2 end">
                    <select className={expiryClass}
                            value={this.props.card.year}
                            onChange={this._onSelectChange('card', 'year')}>
                      {years.map(_mkOptions, this)}
                    </select>
                  </div>
                </div>
                <div className="row">
                  <div className='columns large-8 end'>
                    <DUXInput label="Billing Address"
                              hasError={this.props.fieldErrors.address1}
                              error="Address Required"
                              onChange={this._onChange('billing', 'address1')}
                              value={this.props.billingInfo.address1}/>
                  </div>
                  <div className='columns large-8 end'>
                    <DUXInput label="Apt/Suite"
                            onChange={this._onChange('billing', 'address2')}
                            value={this.props.billingInfo.address2}/>
                  </div>
                  <div className='columns large-8 end'>
                    <DUXInput label="City"
                              hasError={this.props.fieldErrors.city}
                              error='Required'
                              onChange={this._onChange('billing', 'city')}
                              value={this.props.billingInfo.city}/>
                  </div>
                </div>
                <div className="row">
                  <div className="columns large-12">
                    <select className={countryClass}
                            value={this.props.billingInfo.country}
                            onChange={this._onSelectChange('billing', 'country')}>
                      {map(countries, _mkCountryOptions)}
                    </select>
                  </div>
                </div>
                <PostalComponent billingInfo={this.props.billingInfo}
                                 fieldErrors={this.props.fieldErrors}
                                 onChange={this._onChange}
                                 onSelectChange={this._onSelectChange}/>
              </div>
            </div>
            <div className={styles.globalError }>
              { this.props.errorMessage }
            </div>
            <div className='row'>
              <div className='columns large-9'>
                { submitButton }
              </div>
              <div className='columns large-3'>
                <Button intent='secondary'
                        size='small'
                        onClick={this._onBackClick}>Back</Button>
              </div>
            </div>
          </form>
        </Block>
      </Card>
    );
  }
});

var PostalComponent = React.createClass({
  propTypes: {
    billingInfo: shape({
      city: string,
      state: string,
      zip: string,
      first_name: string,
      last_name: string,
      address1: string,
      address2: string,
      country: string
    }),
    onChange: func.isRequired,
    onSelectChange: func.isRequired
  },
  render: function() {
    var countryState;
    var country = this.props.billingInfo.country;
    var stateClass = this.props.fieldErrors.state ? 'billing-dropdown error' : 'billing-dropdown';
    //IF US TERRITORY - ADD STATES SELECT
    if (includes(['US', 'UM'], country)) {
      countryState = (
        <div className="columns large-8">
          <select className={stateClass}
                  value={this.props.billingInfo.state}
                  onChange={this.props.onSelectChange('billing', 'state')}>
            {states.map(_mkOptions, this)}
          </select>
        </div>
      );
    } else {
      countryState = (
        <div className="columns large-8">
          <DUXInput label="State/Province"
                    onChange={this.props.onChange('billing', 'state')}
                    value={this.props.billingInfo.state}/>
        </div>
      );
    }
    return (
      <div className="row">
        {countryState}
        <div className="columns large-4">
          <DUXInput label="Postal Code"
                  onChange={this.props.onChange('billing', 'zip')}
                  value={this.props.billingInfo.zip}/>
        </div>
      </div>
    );
  }
});

export default BillingInfoForm;
