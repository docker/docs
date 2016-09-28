import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import isEmpty from 'lodash/isEmpty';
import Expiration from './CardExpiration';
import validate from './validations';
import {
  Button,
  Input,
} from '../../lib/common';

import css from './styles.css';
const { bool, func, object } = PropTypes;

class PaymentMethodsForm extends Component {
  static propTypes = {
    billingPaymentError: object,
    defaultSelected: bool.isRequired,
    fields: object.isRequired,
    handleSubmit: func.isRequired,
  }

  onSelectChange = (field) => (data) => {
    const fieldObject = this.props.fields[field];
    fieldObject.onChange(data.value);
  }

  render() {
    const {
      billingPaymentError: error,
      defaultSelected,
      fields: propFields,
    } = this.props;
    const {
      cardNumber,
      cvv,
      firstName,
      lastName,
    } = propFields;
    const errClass = !isEmpty(error) ? css.error : '';
    const warning =
      'This will update your payment method for all subscriptions';
    let submitText = 'Change Card';
    if (!defaultSelected) {
      submitText = 'Add Card';
    }

    return (
      <div>
        <div className={`${css.title} ${errClass}`}>
          Credit card information
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
          <div className={css.cardRow}>
            <Input
              {...cardNumber}
              id={'card'}
              placeholder="Card number"
              errorText={
                (cardNumber.touched && cardNumber.error) ||
                (error.type === 'card_error' && error.message)
              }
            />
            <Input
              {...cvv}
              id={'cvv'}
              placeholder="CVV"
              errorText={
                (cvv.touched && cvv.error) || error.type === 'card_error'
              }
            />
            <Expiration
              fields={propFields}
              onSelectChange={this.onSelectChange}
            />
          </div>
          <Button
            className={css.submit}
            type="submit"
            onClick={this.props.handleSubmit}
          >
            {submitText}
          </Button>
          {warning}
        </div>
      </div>
    );
  }
}

const fields = [
  'cardNumber',
  'cvv',
  'firstName',
  'lastName',
  'expMonth',
  'expYear',
];

const mapStateToProps = ({ account, billing }) => {
  const {
    profiles,
    paymentMethods,
  } = billing;
  const {
    namespaceObjects,
    selectedNamespace,
  } = account;
  const selectedUser = namespaceObjects.results[selectedNamespace] || {};
  const billingProfile = profiles.results[selectedUser.id];
  return {
    billingPaymentError: paymentMethods.error,
    billingProfile,
  };
};

export default reduxForm({
  form: 'paymentMethodsForm',
  fields,
  validate,
},
mapStateToProps,
)(PaymentMethodsForm);
