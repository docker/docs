import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import map from 'lodash/map';
import isEmpty from 'lodash/isEmpty';
import sortBy from 'lodash/sortBy';
import find from 'lodash/find';
import {
  billingCreatePaymentMethod,
  billingDeletePaymentMethod,
  billingFetchPaymentMethods,
  billingSetDefaultPaymentMethod,
  billingUpdateProfile,
} from '../actions';
import {
  Card,
  Button,
  FetchingError,
  FullscreenLoading,
  WarningIcon,
} from '../lib/common';
import {
  TINY,
} from '../lib/constants';
import PaymentForm from './PaymentForm';
import ContactForm from './ContactForm';
import css from './styles.css';

const dispatcher = {
  billingCreatePaymentMethod,
  billingDeletePaymentMethod,
  billingFetchPaymentMethods,
  billingSetDefaultPaymentMethod,
  billingUpdateProfile,
};

const mapStateToProps = ({ billing }) => {
  const {
    paymentMethods,
  } = billing;
  return {
    paymentMethods,
  };
};

const { object, bool, shape, func } = PropTypes;

@connect(mapStateToProps, dispatcher)
export default class PaymentMethodsView extends Component {
  static propTypes = {
    billingProfile: object.isRequired,
    paymentMethods: shape({
      isFetching: bool.isRequired,
      results: object.isrequired,
    }),
    selectedUser: object.isRequired,
    // actions
    billingCreatePaymentMethod: func.isRequired,
    billingDeletePaymentMethod: func.isRequired,
    billingFetchPaymentMethods: func.isRequired,
    billingSetDefaultPaymentMethod: func.isRequired,
    billingUpdateProfile: func.isRequired,
  }

  state = {
    selectedCardId: '',
    defaultPaymentWarning: '',
  }

  setDefaultCard = (bfId) => () => {
    const {
      selectedUser,
      billingSetDefaultPaymentMethod: setDefaultPaymentMethod,
      billingFetchPaymentMethods: fetchPaymentMethods,
    } = this.props;
    const { id: docker_id } = selectedUser;
    setDefaultPaymentMethod({
      card_id: bfId,
      docker_id,
    }).then(() => {
      fetchPaymentMethods({ docker_id });
    });
  }

  addCard = (data, shouldFetch) => {
    const {
      billingProfile,
      selectedUser,
      billingCreatePaymentMethod: createPaymentMethod,
      billingFetchPaymentMethods: fetchPaymentMethods,
    } = this.props;
    const { id: docker_id } = selectedUser;
    const billforward_id = billingProfile && billingProfile.billforward_id;
    const paymentInfo = {
      billforward_id,
      name_first: data.firstName,
      name_last: data.lastName,
      cvc: data.cvv,
      number: data.cardNumber,
      exp_month: data.expMonth,
      exp_year: data.expYear,
    };
    if (!!shouldFetch) {
      return createPaymentMethod(paymentInfo).then(() => {
        fetchPaymentMethods({ docker_id });
      });
    }
    return createPaymentMethod(paymentInfo);
  }

  replaceCard = (bfId) => (data) => {
    const addCard = this.addCard;
    const deleteCard = this.deleteCard;
    addCard(data).then(() => {
      deleteCard(bfId)();
    });
  }

  deleteCard = (bfId) => () => {
    const {
      selectedUser,
      billingDeletePaymentMethod: deletePaymentMethod,
      billingFetchPaymentMethods: fetchPaymentMethods,
    } = this.props;
    const { id: docker_id } = selectedUser;
    deletePaymentMethod({ docker_id, card_id: bfId }).then(() => {
      fetchPaymentMethods({ docker_id });
    });
  }

  updateBillingProfile = (values) => {
    const {
      selectedUser,
      billingUpdateProfile: updateProfile,
    } = this.props;
    const {
      address1,
      city,
      company,
      country,
      email,
      firstName,
      job,
      lastName,
      phone,
      postalCode,
      province,
    } = values;
    const docker_id = selectedUser.id;
    const address = {
      address_line_1: address1,
      city,
      province,
      country,
      post_code: postalCode,
      primary_address: true,
    };
    const submitData = {
      addresses: [
        address,
      ],
      company_name: company,
      docker_id,
      email,
      first_name: firstName,
      job_function: job,
      last_name: lastName,
      phone_primary: phone,
    };
    updateProfile(submitData).then(() => {
      const keys = selectedUser.type === 'User' ? {
        dockerUUID: selectedUser.id,
        Docker_Hub_User_Name__c: selectedUser.username,
      } : {
        dockerUUID: selectedUser.id,
        Docker_Hub_Organization_Name__c: selectedUser.orgname,
      };

      analytics.track('account_info_update', {
        ...keys,
        email,
        firstName,
        lastName,
        jobFunction: job,
        address: address1,
        postalCode,
        state: province,
        city,
        country,
        phone,
      });
    });
  }

  renderCardInfo = (card, idx) => {
    const {
      bf_payment_method_id: cardId,
      description,
      expiry_month,
      expiry_year,
      card_type,
      default: defaultCard,
    } = card;
    const isDefault = defaultCard ? 'Default' : '';
    const endsWith = description.replace(/#*/, 'x');
    let selected =
      cardId === this.state.selectedCardId ? css.selected : '';
    if (!this.state.selectedCardId && defaultCard) {
      selected = css.selected;
    }
    const selectCard = () => { this.setState({ selectedCardId: cardId }); };
    return (
      <Card
        className={`${css.paymentMethod} ${selected}`}
        onClick={selectCard}
        key={idx}
      >
        {isDefault}
        <div>
          {card_type} ending in {endsWith}
        </div>
        <div>
          Expires: {expiry_month}/{expiry_year}
        </div>
      </Card>
    );
  }

  render() {
    // TODO design - nathan make this GORGEOUS. because it's ugly af right now
    const {
      paymentMethods,
    } = this.props;
    const {
      selectedCardId,
    } = this.state;
    if (paymentMethods.isFetching) {
      return <FullscreenLoading />;
    } else if (!isEmpty(paymentMethods.fetchingError)) {
      return (
        <div className={css.fetchingError}>
          <FetchingError resource="your payment methods" />
        </div>
      );
    } else if (paymentMethods.results.length < 1) {
      return (
        <Card shadow>
          No Payment Methods to show
        </Card>
      );
    }
    let cardActions;
    const defaultCard = find(paymentMethods.results, card => card.default);
    let defaultSelected = false;
    let paymentFormSubmit = this.addCard;
    if (selectedCardId && defaultCard.bf_payment_method_id !== selectedCardId) {
      cardActions = (
        <div className={css.updateCardRow}>
          <div>
            <Button
              className={css.submit}
              onClick={this.setDefaultCard(selectedCardId)}
            >
              Make Default Card {this.state.hover}
            </Button>
            <Button
              className={css.submit}
              onClick={this.deleteCard(selectedCardId)}
            >
              Delete Payment Method
            </Button>
          </div>
          <div className={css.warning}>
            <WarningIcon size={TINY} />&nbsp;
            Warning: Changing your default payment method will
              change your default payment for ALL subscriptions
          </div>
        </div>
      );
    } else if (
      defaultCard.bf_payment_method_id === selectedCardId || !selectedCardId
    ) {
      defaultSelected = true;
      paymentFormSubmit = this.replaceCard(defaultCard.bf_payment_method_id);
    }

    const sortedPaymentMethods =
      sortBy(paymentMethods.results, method => !method.default);

    const splitForms = (
      <div className={css.splitForms}>
        <PaymentForm
          onSubmit={paymentFormSubmit}
          defaultSelected={defaultSelected}
        />
        <ContactForm onSubmit={this.updateBillingProfile} />
      </div>
    );

    return (
      <Card title="Payment Information" shadow>
        <div className={css.infoRow}>
          <div className={css.title}>Credit Cards</div>
          {map(sortedPaymentMethods, this.renderCardInfo)}
        </div>
        {cardActions}
        {splitForms}
      </Card>
    );
  }
}
