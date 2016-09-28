import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import {
  countryOptions,
} from '../../lib/countries';
import {
  Button,
  Input,
  Select,
} from '../../lib/common';
import EmailSelect from '../EmailSelect';
import sortBy from 'lodash/sortBy';
import merge from 'lodash/merge';
import css from './styles.css';
const { array, func, object } = PropTypes;

class ContactForm extends Component {
  static propTypes = {
    emailsArray: array.isRequired,
    fields: object.isRequired,
    handleSubmit: func.isRequired,
  }

  onSelectChange = (field) => (data) => {
    const fieldObject = this.props.fields[field];
    fieldObject.onChange(data.value);
  }

  render() {
    const {
      emailsArray,
      fields: propFields,
    } = this.props;
    const {
      address1,
      city,
      company,
      country,
      firstName,
      job,
      lastName,
      phone,
      postalCode,
      province,
    } = propFields;
    return (
      <div>
        <div className={css.title} >
          Contact Information
        </div>
        <div className={css.form}>
          <div className={css.row}>
            <Input
              {...firstName}
              id={'first'}
              placeholder="First Name"
              errorText={firstName.touched && firstName.error}
            />
            <Input
              {...lastName}
              id={'last'}
              placeholder="Last Name"
              errorText={lastName.touched && lastName.error}
            />
          </div>
          <div className={css.row}>
            <Input
              {...company}
              id={'company'}
              placeholder="Company"
            />
            <EmailSelect
              accountEmails={emailsArray}
              fields={propFields}
              onSelectChange={this.onSelectChange('email')}
            />
          </div>
          <div className={css.row}>
            <Input
              {...job}
              id={'job'}
              placeholder="Job"
            />
            <Input
              {...phone}
              id={'phone'}
              placeholder="Phone"
            />
          </div>
          <div className={css.title} >
            Address
          </div>
          <div className={css.fullRow}>
            <Input
              {...address1}
              style={{ width: '100%' }}
              id={'address'}
              placeholder="Address"
            />
          </div>
          <div className={css.row}>
            <Input
              {...city}
              id={'city'}
              placeholder="City"
            />
            <Input
              {...postalCode}
              id={'postalCode'}
              placeholder="Postal Code"
            />
          </div>
          <div className={css.row}>
            <Input
              {...province}
              id={'province'}
              placeholder="Province"
            />
            <Select
              {...country}
              onBlur={() => {}}
              className={css.select}
              placeholder="Country"
              style={{ marginBottom: '10px', width: '' }}
              options={countryOptions}
              onChange={this.onSelectChange('country')}
              ignoreCase
              clearable={false}
              errorText={
                country.touched && country.error ? country.error : ''
              }
            />
          </div>
          <Button
            className={css.submit}
            type="submit"
            onClick={this.props.handleSubmit}
          >
            Update Contact Information
          </Button>
        </div>
      </div>
    );
  }
}

const fields = [
  'address1',
  'city',
  'company',
  'country',
  'email',
  'firstName',
  'job',
  'lastName',
  'phone',
  'postalCode',
  'province',
];

const mapStateToProps = ({ account, billing }) => {
  const {
    profiles,
  } = billing;
  const {
    namespaceObjects,
    selectedNamespace,
    userEmails,
  } = account;
  const emailList = sortBy(userEmails.results, o => !o.primary);
  const emailsArray = emailList.map((emailObject) => {
    return emailObject.email;
  });
  const selectedUser = namespaceObjects.results[selectedNamespace] || {};
  const billingProfile = profiles.results[selectedUser.id];
  const initialData = {
    account: selectedNamespace,
    email: emailList[0],
  };
  let initialized = {};
  if (!!billingProfile) {
    const selectedEmail = billingProfile.email;
    if (emailsArray.indexOf(selectedEmail) < 0) {
      emailsArray.unshift(selectedEmail);
    }
    const address =
      billingProfile.addresses && billingProfile.addresses[0] || {};
    initialized = {
      address1: address.address_line_1,
      city: address.city,
      company: billingProfile.company_name,
      country: address.country,
      email: billingProfile.email,
      firstName: billingProfile.first_name,
      job: billingProfile.job_function,
      lastName: billingProfile.last_name,
      phone: billingProfile.phone_primary,
      postalCode: address.post_code,
      province: address.province,
    };
    merge(initialData, initialized);
  }
  return {
    emailsArray,
    initialValues: initialData,
  };
};

export default reduxForm({
  form: 'ContactForm',
  fields,
  // validate,
},
mapStateToProps,
)(ContactForm);
