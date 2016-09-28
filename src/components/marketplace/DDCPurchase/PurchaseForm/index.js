import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import {
  chain,
  find,
  has,
  isEmpty,
  map,
  merge,
} from 'lodash';
import {
  Button,
  Card,
  Checkbox,
  Input,
  RadioButton,
  RadioButtonGroup,
  Select,
  Slider,
} from 'common';
import { isPaidPlan } from 'lib/utils/billing-plan-utils';
import { countryOptions } from 'lib/constants/countries';
import {
  DDC_BUSINESS_DAY,
  DDC_ID,
  EUSA_LINK,
} from 'lib/constants/eusa';
import {
  accountFetchUser,
  accountSelectNamespace,
} from 'actions/account';
import {
  billingFetchProfile,
  billingFetchPaymentMethods,
} from 'actions/billing';
import CardExpiration from 'components/marketplace/Purchase/CardExpiration';
import AccountSelect from 'components/marketplace/Purchase/AccountSelect';
import EmailSelect from 'components/marketplace/Purchase/EmailSelect';
import validate from '../purchaseValidations';
import css from './styles.css';

const { array, string, object, func, shape } = PropTypes;

class PurchaseForm extends Component {
  static propTypes = {
    accountEmails: array.isRequired,
    accountOptions: array.isRequired,
    billingPaymentError: object,
    billingPaymentMethods: array,
    billingProductDetails: object,
    billingProfile: object,
    billingSubscriptions: object,
    components: object,
    currentUser: object.isRequired,
    fields: object.isRequired,
    initialized: object.isRequired,
    location: object.isRequired,
    fetchedNamespaces: object.isRequired,
    productDetails: object,
    profileError: string,
    selectedPlanObject: object.isRequired,
    // actions
    accountFetchUser: func.isRequired,
    accountSelectNamespace: func.isRequired,
    billingFetchPaymentMethods: func.isRequired,
    billingFetchProfile: func.isRequired,
    changePlan: func.isRequired,
    updateComponents: func.isRequired,
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  onPlanChange = (event, value) => {
    const { pathname } = this.props.location;
    const { router } = this.context;
    router.push({
      pathname,
      query: {
        plan: value,
      },
    });
    this.props.changePlan(value);
  }

  onCheck = () => {
    const { accepted } = this.props.fields;
    const checked = accepted.value === 'checked' ? '' : 'checked';
    accepted.onChange(checked);
  }

  onSelectChange = (field) => (data) => {
    if (field === 'account') {
      const {
        accountFetchUser: fetchUser,
        accountSelectNamespace: selectNamespace,
        billingFetchPaymentMethods: fetchPaymentMethods,
        billingFetchProfile: fetchBillingProfiles,
        currentUser,
        fetchedNamespaces,
      } = this.props;
      if (fetchedNamespaces[data.value]) {
        const userOrg = fetchedNamespaces[data.value];
        const { type, id: docker_id } = userOrg;
        const isOrg = type === 'Organization';
        selectNamespace({ namespace: data.value });
        fetchBillingProfiles({ docker_id, isOrg }).then((res) => {
          if (res.value.profile) {
            fetchPaymentMethods({ docker_id });
          }
        });
      } else {
        const isOrg = data.value !== currentUser.username;
        // fetching user info vs. org info locally hits different api's which
        // requires knowing which endpoint to hit.
        // unnecessary in production since v2/users will redirect to v2/orgs
        fetchUser({ namespace: data.value, isOrg }).then((userRes) => {
          const { id: docker_id } = userRes.value;
          return Promise.all([
            selectNamespace({ namespace: data.value }),
            fetchBillingProfiles({ docker_id, isOrg })
              .then((billingRes) => {
                if (billingRes.value.profile) {
                  fetchPaymentMethods({ namespace: data.value, docker_id });
                }
              }),
          ]);
        });
      }
    }
    const fieldObject = this.props.fields[field];
    fieldObject.onChange(data.value);
  }

  renderPricingComponent = (component, idx) => {
    const {
      components,
      updateComponents,
    } = this.props;
    const { tiers, name } = component;
    const initialTier = tiers[0];
    const lower = initialTier.lower_threshold;
    const upper =
      initialTier.upper_threshold > 100 ? 100 : initialTier.upper_threshold;
    const selected = components[name] || lower;
    return (
      <div className={css.component} key={idx}>
        <div>{component.label} ({selected})</div>
        <Slider
          min={lower}
          max={upper}
          step={1}
          value={components[name]}
          onChange={updateComponents(name)}
        />
      </div>
    );
  }

  renderDDCEngines() {
    const { selectedPlanObject } = this.props;
    const { pricing_components } = selectedPlanObject;
    return (
      <div className={css.engines}>
        <div>
          <div className={css.sectionTitle}>
            Specify number of Nodes
          </div>
          <div className={css.sectionSub}>
            <span>Have more than 100 nodes environment?</span>
            <a
              href={
                'mailto:sales@docker.com?Subject=Docker%20Datacenter%20Purchase'
              }
            >
              Contact us
            </a>
          </div>
        </div>
        <div className={css.pricingComponents}>
          {map(pricing_components, this.renderPricingComponent)}
        </div>
      </div>
    );
  }

  renderDDCPlans() {
    const { billingProductDetails } = this.props;
    if (
      billingProductDetails.isFetching ||
      !has(billingProductDetails, 'rate_plans')
    ) {
      return null;
    }
    let plans = billingProductDetails.rate_plans;
    if (billingProductDetails.id === DDC_ID) {
      plans = chain(plans)
        .filter(isPaidPlan)
        .sortBy((plan) => {
          // sort plans with business_day_monthly first
          return plan.name !== DDC_BUSINESS_DAY;
        })
        .value();

      if (isEmpty(plans)) {
        console.error('expected at least one ddc paid plan');
      }
    }

    const renderRadio = (plan) => {
      return (
        <RadioButton
          key={plan.name}
          value={plan.name}
          label={plan.label}
          labelStyle={{ width: null }}
          iconStyle={{ marginRight: '4px' }}
          style={{ width: null, display: 'inline-flex', marginRight: '16px' }}
        />
      );
    };

    return (
      <div className={css.plans}>
        <div className={css.sectionTitle}>
          Select support level
        </div>
        <RadioButtonGroup
          name="rate_plans"
          defaultSelected={plans[0].name}
          onChange={this.onPlanChange}
          className={css.plansRadio}
        >
          {map(plans, renderRadio)}
        </RadioButtonGroup>
      </div>
    );
  }

  renderPaymentForm() {
    const {
      billingPaymentError: error,
      billingPaymentMethods,
      fields: propFields,
      initialized,
    } = this.props;
    const {
      address,
      city,
      postCode,
      cardNumber,
      country,
      cvv,
      firstName,
      lastName,
      state,
    } = propFields;
    if (billingPaymentMethods.length === 0) {
      const errClass = error.message ? css.softErr : '';
      return (
        <div>
          <div className={`${css.soft} ${errClass}`}>
            Credit card information
          </div>
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
              readOnly={!!initialized.cardNumber}
              id={'card'}
              placeholder="Card number"
              errorText={
                (cardNumber.touched && cardNumber.error) ||
                (error.type === 'card_error' && error.message)
              }
            />
            <Input
              {...cvv}
              readOnly={!!initialized.cvv}
              id={'cvv'}
              placeholder="CVV"
              errorText={
                (cvv.touched && cvv.error) || error.type === 'card_error'
              }
            />
            <CardExpiration
              fields={propFields}
              initialized={initialized}
              onSelectChange={this.onSelectChange}
            />
          </div>
          <div className={css.soft}>Billing information</div>
          <div className={css.row}>
            <Input
              {...address}
              id={'address'}
              placeholder="Address"
              errorText={address.touched && address.error}
            />
            <Input
              {...city}
              id={'city'}
              placeholder="City"
              errorText={city.touched && city.error}
            />
          </div>
          <div className={css.postalRow}>
            <Input
              {...state}
              id={'state'}
              placeholder="State/Province"
              errorText={state.touched && state.error}
            />
            <Input
              {...postCode}
              id={'postCode'}
              placeholder="Postal Code"
              errorText={postCode.touched && postCode.error}
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
        </div>
      );
    }
    // For now only render the default
    // Users can update their default in the billing profile page - however they
    // Should be warned that changing your default payment will change the
    // payment method for ALL their subscriptions
    const defaultPayment = find(billingPaymentMethods, (method) => {
      return method.default;
    });
    return this.renderPaymentMethod(defaultPayment);
  }

  renderPaymentMethod = (method, idx) => {
    const {
      selectedCard,
    } = this.props.fields;
    const {
      bf_payment_method_id: bfId,
      description,
      expiry_month,
      expiry_year,
      card_type,
      default: defaultCard,
    } = method;
    const endsWith = description.replace(/#*/, 'x');
    const selected =
      bfId === selectedCard.value ? css.selected : '';
    const selectCard = () => { selectedCard.onChange(bfId); };
    return (
      <Card
        className={`${css.paymentMethod} ${selected}`}
        onClick={selectCard}
        key={idx}
      >
        <div>
          {defaultCard}
          <div>
            {card_type} ending in {endsWith}
          </div>
          <div>
            Expires: {expiry_month}/{expiry_year}
          </div>
        </div>
      </Card>
    );
  }

  renderProfileForm() {
    // If NO BILLING PROFILE - Show billing info forms to create billing profile
    const {
      fields: propFields,
      accountEmails,
      billingProfile,
      initialized,
    } = this.props;
    const {
      company,
      job,
    } = propFields;
    let form;
    if (!billingProfile || billingProfile.error) {
      form = (
        <div>
          <div className={css.soft}>Profile information</div>
          <div className={css.row}>
            <Input
              {...company}
              id={'company'}
              placeholder="Company"
              errorText={company.touched && company.error}
            />
            <Input
              {...job}
              id={'job'}
              placeholder="Title"
              errorText={job.touched && job.error}
            />
          </div>
          <div className={css.row}>
            <EmailSelect
              fields={propFields}
              initialized={initialized}
              accountEmails={accountEmails}
              onSelectChange={this.onSelectChange('email')}
            />
          </div>
        </div>
      );
    }
    return form;
  }

  renderAdditionalOptions() {
    const {
      coupon,
    } = this.props.fields;
    return (
      <div>
        <div className={css.sectionTitle}>
          Additional Options
        </div>
        <div className={css.couponRow}>
          <Input
            {...coupon}
            id={'coupon'}
            placeholder="Coupon code"
            errorText={coupon.touched && coupon.error}
          />
          <Button inverted disabled>Apply</Button>
        </div>
      </div>
    );
  }

  render() {
    const {
      accountOptions,
      fields: propFields,
      productDetails,
      billingPaymentError,
      billingSubscriptions,
      profileError,
    } = this.props;
    const {
      accepted,
    } = propFields;

    const title = `Purchase ${productDetails.name}`;
    const subtext = `
      This subscription gives you access to the tools you need to manage your
      Docker images, containers, and applications within your firewall.
    `;
    let globalErr;
    if (
      billingPaymentError.message && billingPaymentError.type !== 'card_error'
    ) {
      globalErr = (
        <div>{billingPaymentError.message}</div>
      );
    } else if (billingSubscriptions.error) {
      globalErr = (
        <div>{billingSubscriptions.error}</div>
      );
    } else if (profileError) {
      globalErr = (
        <div>{profileError}</div>
      );
    }
    return (
      <Card
        className={css.card}
        title={title}
      >
        <div>
          <div className={css.subtitle}>
            {subtext}
          </div>
          <AccountSelect
            options={accountOptions}
            fields={propFields}
            onSelectChange={this.onSelectChange('account')}
          />
          {this.renderDDCEngines()}
          {this.renderDDCPlans()}
          <div className={css.paymentDetails}>
            <div className={css.sectionTitle}>
              Payment Details
            </div>
            {this.renderPaymentForm()}
            {this.renderProfileForm()}
          </div>
          {/* this.renderAdditionalOptions() */}
          <div className={css.eusa}>
            <Checkbox
              {...accepted}
              checked={accepted.value === 'checked'}
              onBlur={() => {}}
              onCheck={this.onCheck}
              style={{ width: '' }}
            />
            <div>
              I have read and agree to Docker's&nbsp;
              <a href={EUSA_LINK} >
                End User Subscription Agreement
              </a>
            </div>
          </div>
          <div className={css.error}>
            {(accepted.touched && accepted.error) || globalErr}
          </div>
        </div>
      </Card>
    );
  }
}

const fields = [
  'accepted',
  'account',
  'coupon',
  // Billing Contact Info
  'address',
  'city',
  'company',
  'country',
  'email',
  'state',
  'job',
  'postCode',
  // Billing Card Info
  'cardNumber',
  'cvv',
  'firstName',
  'lastName',
  'expMonth',
  'expYear',
  'selectedCard',
];

const mapStateToProps = ({ account, billing }, props) => {
  const {
    profiles,
    paymentMethods,
    subscriptions,
  } = billing;
  const {
    currentUser,
    namespaceObjects,
    selectedNamespace,
  } = account;
  const selectedUser = namespaceObjects.results[selectedNamespace] || {};
  const billingProfile = profiles.results[selectedUser.id];
  const initialData = {
    account: selectedNamespace,
    email: props.accountEmails[0],
  };
  let initialized = {};
  const defaultPayment = find(paymentMethods.results, (card) => {
    return card.default;
  });
  let initialCard = {};
  let selectedCard;
  if (defaultPayment) {
    selectedCard = defaultPayment.bf_payment_method_id;
    initialCard = {
      cardNumber: defaultPayment.description,
      cvv: '000',
      expMonth: '20',
      expYear: '4000',
    };
  }
  if (!!billingProfile) {
    const address =
      billingProfile.addresses && billingProfile.addresses[0] || {};
    initialized = {
      address: address.address_line_1,
      cardNumber: initialCard.cardNumber,
      city: address.city,
      company: billingProfile.company_name,
      country: address.country,
      cvv: initialCard.cvv,
      email: billingProfile.email,
      expMonth: initialCard.expMonth,
      expYear: initialCard.expYear,
      firstName: billingProfile.first_name,
      job: billingProfile.job_function,
      lastName: billingProfile.last_name,
      phone: billingProfile.phone_primary,
      postCode: address.post_code,
      selectedCard,
      state: address.province,
      accepted: 'unchecked',
    };
    merge(initialData, initialized);
  }
  return {
    billingPaymentError: paymentMethods.error,
    billingPaymentMethods: paymentMethods.results,
    billingProfile,
    billingSubscriptions: subscriptions,
    currentUser,
    initialized,
    initialValues: initialData,
    fetchedNamespaces: namespaceObjects.results,
    profileError: profiles.error,
  };
};

const dispatcher = {
  accountFetchUser,
  accountSelectNamespace,
  billingFetchProfile,
  billingFetchPaymentMethods,
};

export default reduxForm({
  form: 'purchaseForm',
  fields,
  validate,
},
mapStateToProps,
dispatcher,
)(PurchaseForm);
