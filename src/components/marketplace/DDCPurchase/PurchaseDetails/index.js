import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import map from 'lodash/map';
import merge from 'lodash/merge';
import find from 'lodash/find';
import reduce from 'lodash/reduce';
import isEmpty from 'lodash/isEmpty';
import {
  Button,
  Card,
  DDCGraphic,
  ImageWithFallback,
} from 'common';
import { DDC } from 'lib/constants/landingPage';
import { FALLBACK_IMAGE_SRC, FALLBACK_ELEMENT } from 'lib/constants/fallbacks';
import { SUBSCRIPTIONS } from 'lib/constants/overlays';
import getLargestLogo from 'lib/utils/get-largest-logo';
import css from './styles.css';
import validate from '../purchaseValidations';
import {
  accountToggleMagicCarpet,
} from 'actions/account';
import {
  billingCreatePaymentMethod,
  billingCreateProfile,
  billingCreateSubscription,
  billingFetchPaymentMethods,
  billingSetDefaultPaymentMethod,
} from 'actions/billing';

const { string, array, bool, object, func, shape } = PropTypes;

class PurchaseDetails extends Component {
  static propTypes = {
    fields: object.isRequired,
    formSubmitting: bool.isRequired,
    billingPaymentMethods: array,
    billingProductDetails: object,
    billingProfile: object,
    components: object, // selected pricing component values
    currentUser: object,
    location: object.isRequired,
    selectedUser: object.isRequired,
    productDetails: shape({
      name: string.isRequired,
      short_description: string,
      logo_url: object,
    }),
    selectedPlanObject: object,
    // actions
    handleSubmit: func.isRequired,
    accountToggleMagicCarpet: func.isRequired,
    billingCreatePaymentMethod: func.isRequired,
    billingCreateProfile: func.isRequired,
    billingCreateSubscription: func.isRequired,
    billingFetchPaymentMethods: func.isRequired,
    billingSetDefaultPaymentMethod: func.isRequired,
  }

  static contextTypes = {
    router: object.isRequired,
  }

  submitForm = (values) => {
    const {
      billingProfile,
      billingPaymentMethods,
      billingProductDetails,
      components,
      location,
      selectedUser,
      selectedPlanObject,
      // actions
      accountToggleMagicCarpet: toggleMagicCarpet,
      billingCreateProfile: createProfile,
      billingCreatePaymentMethod: createPaymentMethod,
      billingCreateSubscription: createSubscription,
      billingFetchPaymentMethods: fetchPaymentMethods,
      billingSetDefaultPaymentMethod: setDefaultPaymentMethod,
    } = this.props;
    const isOrg = selectedUser.type === 'Organization';
    let billforward_id = billingProfile ? billingProfile.billforward_id : '';
    const auth = {
      billforward_id,
      isOrg,
      docker_id: selectedUser.id,
    };
    const profile = {
      first_name: values.firstName,
      last_name: values.lastName,
      email: values.email,
      phone_primary: values.phone,
      company_name: values.company,
      job_function: values.job,
      addresses: [{
        province: values.state,
        country: values.country,
        address_line_1: values.address,
        city: values.city,
        post_code: values.postCode,
      }],
    };
    const paymentInfo = {
      name_first: values.firstName,
      name_last: values.lastName,
      cvc: values.cvv,
      number: values.cardNumber,
      exp_month: values.expMonth,
      exp_year: values.expYear,
    };
    const date = new Date();
    const { name: planName, pricing_components } = selectedPlanObject;
    const submitPricingComponents = [];
    pricing_components.forEach((component) => {
      const name = component.name;
      const minimum = component.tiers[0].lower_threshold;
      const value = components[name] || minimum;
      submitPricingComponents.push({ name, value });
    });
    const prodName = billingProductDetails.label || '';
    const subData = {
      name: `${prodName} Subscription ${date.toDateString()}`,
      product_rate_plan: planName,
      pricing_components: submitPricingComponents,
      product_id: billingProductDetails.id,
      eusa: {
        accepted: !!values.accepted,
      },
    };
    const router = this.context.router;
    // grabbing path from location since we cannot tell whether this is a bundle
    // or an image otherwise
    const detailRoute = location.pathname.split('/purchase')[0];
    if (!billingProfile || billingProfile.error) {
      const data = { ...auth, profile };
      createProfile(data).then((profRes) => {
        billforward_id = profRes.action.payload.profile.billforward_id;
        const paymentData = { ...paymentInfo, billforward_id };
        createPaymentMethod(paymentData).then(() => {
          // create subscription uses default payment method
          Promise.when([
            createSubscription(merge({}, subData, auth)).then(() => {
              // redirect to details page on success with subscriptions
              // carpet open
              toggleMagicCarpet({ magicCarpet: SUBSCRIPTIONS });
              router.push(`${detailRoute}?overlay=subscriptions`);
              analytics.track('ddc_purchase', {
                plan: planName,
              });
            }),
            fetchPaymentMethods(auth),
          ]);
        });

        const keys = selectedUser.type === 'User' ? {
          selectedUser: selectedUser.id,
          Docker_Hub_User_Name__c: selectedUser.username,
        } : {
          dockerUUID: selectedUser.id,
          Docker_Hub_Organization_Name__c: selectedUser.orgname,
        };

        analytics.track('account_info_update', {
          ...keys,
          email: values.email,
          firstName: values.firstName,
          lastName: values.lastName,
          jobFunction: values.job,
          address: values.address1,
          postalCode: values.postCode,
          state: values.state,
          city: values.city,
          country: values.country,
          phone: values.phone,
        });
      });
    } else if (billingPaymentMethods.length === 0) {
      const data = merge({}, auth, paymentInfo);
      // Create Payment Method - Then create Subscription
      createPaymentMethod(data).then(() => {
        Promise.when([
          createSubscription(merge({}, subData, auth)).then(() => {
            // redirect to details page on success with subscription carpet open
            toggleMagicCarpet({ magicCarpet: SUBSCRIPTIONS });
            router.push(`${detailRoute}?overlay=subscriptions`);
            analytics.track('ddc_purchase', {
              plan: planName,
            });
          }),
          fetchPaymentMethods(auth),
        ]);
      });
    } else {
      const defaultPayment = find(billingPaymentMethods, (card) => {
        return card.default;
      });
      // selectedCard is the bf card ID
      if (values.selectedCard !== defaultPayment.bf_payment_method_id) {
        // Update default card > THEN create Subscription
        setDefaultPaymentMethod({
          card_id: values.selectedCard,
          ...auth,
        }).then(() => {
          createSubscription(merge({}, subData, auth)).then(() => {
            // redirect to details page on success with subscription carpet open
            toggleMagicCarpet({ magicCarpet: SUBSCRIPTIONS });
            router.push(`${detailRoute}?overlay=subscriptions`);
            analytics.track('ddc_purchase', {
              plan: planName,
            });
          });
        });
      } else {
        // Just create Subscription
        createSubscription(merge({}, subData, auth)).then(() => {
          // redirect to details page on success with subscription carpet open
          toggleMagicCarpet({ magicCarpet: SUBSCRIPTIONS });
          router.push(`${detailRoute}?overlay=subscriptions`);
          analytics.track('ddc_purchase', {
            plan: planName,
          });
        });
      }
    }
  }

  calculateTierBuckets = (component) => {
    const {
      components: selectedComponents,
    } = this.props;
    const {
      name,
      tiers,
    } = component;
    const value = selectedComponents[name];
    if (!value) {
      const {
        price,
        lower_threshold,
      } = tiers[0];
      return price * lower_threshold;
    }
    return reduce(tiers, (sum, bucket) => {
      const {
        price,
        upper_threshold,
        lower_threshold,
      } = bucket;
      const bucketCount =
        Math.min(upper_threshold, (value - lower_threshold + 1));
      if (bucketCount < 1) {
        return sum;
      }
      return sum + (bucketCount * price);
    }, 0);
  }

  calculateTotalCost = () => {
    const {
      selectedPlanObject,
    } = this.props;
    return reduce(
      selectedPlanObject.pricing_components,
      (sum, component) => {
        return sum + this.calculateTierBuckets(component);
      },
      0,
    );
  }

  renderTierBreakdown = (component, idx) => {
    const {
      components: selectedComponents,
    } = this.props;
    const {
      name,
      tiers,
    } = component;
    const value = selectedComponents[name];
    if (!value) {
      const minimum = tiers[0].lower_threshold;
      const price = tiers[0].price;
      return (
        <div className={css.pricingTier} key={idx}>
          <div>{name} @ {price} x {minimum}</div>
          <div>{`$${price * minimum}`}</div>
        </div>
      );
    }
    const tierBreakdown = [];
    reduce(tiers, (sum, bucket) => {
      const {
        price,
        upper_threshold,
        lower_threshold,
      } = bucket;

      // See how many units are within the bucket for this tier;
      const bucketCount =
        Math.min(upper_threshold, (value - lower_threshold + 1));
      // If there are we have units in this bucket add to breakdown
      if (bucketCount > 0) {
        tierBreakdown.push(
          <div className={css.pricingTier} key={lower_threshold}>
            <div>{name} @ {price} x {bucketCount}</div>
            <div>{`$${price * bucketCount}`}</div>
          </div>
        );
      }
      return sum + upper_threshold;
    }, 0);
    return tierBreakdown;
  }

  renderComponentBreakdown() {
    const {
      selectedPlanObject,
    } = this.props;
    const { pricing_components } = selectedPlanObject;
    return map(pricing_components, this.renderTierBreakdown);
  }

  render() {
    const {
      formSubmitting,
      productDetails,
      selectedPlanObject,
      handleSubmit,
      currentUser,
    } = this.props;
    let details;
    if (selectedPlanObject.error) {
      details = (
        <div>404 Plan does not exist</div>
      );
    } else if (!selectedPlanObject) {
      details = (
        <div>fetching ...</div>
      );
    } else {
      const {
        duration_period,
        currency,
       } = selectedPlanObject;
      const { name } = productDetails;
      const logoUrl = getLargestLogo(productDetails.logo_url);
      let logo;
      if (!logoUrl) {
        logo = (
          <DDCGraphic
            size="xlarge"
            className={css.ddcgraphic}
            variant="dull"
          />
        );
      } else {
        logo = (
          <ImageWithFallback
            src={logoUrl}
            className={css.icon}
            fallbackImage={FALLBACK_IMAGE_SRC}
            fallbackElement={FALLBACK_ELEMENT}
          />
        );
      }

      const cost = this.calculateTotalCost();
      const period = duration_period && duration_period.replace(/s$/, '');
      const included = DDC.whatsIncluded;
      let breakdown;
      const price = <span className={css.price}>${cost}</span>;
      if (!isEmpty(currentUser)) {
        breakdown = (
          <div>
            <hr />
            {this.renderComponentBreakdown()}
            <hr />
            <div className={css.cost}>
              {price} / {period} {currency}
            </div>
            <div>
              <Button
                id="submit"
                className={css.submit}
                disabled={formSubmitting}
                type="submit"
                onClick={handleSubmit(this.submitForm)}
              >
                Get a {name} Subscription Today
              </Button>
            </div>
          </div>
        );
      }

      details = (
        <div className={css.details}>
          <div>
            {logo}
          </div>
          <div className={css.included}>
            <div className={css.title}>What&#39;s Included</div>
            {selectedPlanObject.label}
            <div className={css.includeItems}>
            {included.map((item, idx) => (<div key={idx}>{item}</div>))}
            </div>
          </div>
          {breakdown}
        </div>
      );
    }
    return (
      <Card
        className={css.card}
      >
      {details}
      </Card>
    );
  }
}

const fields = [
  'accepted',
  'account',
  'address',
  'cardNumber',
  'city',
  'company',
  'country',
  'coupon',
  'cvv',
  'email',
  'expMonth',
  'expYear',
  'firstName',
  'job',
  'lastName',
  'postCode',
  'selectedCard',
  'state',
];

const mapStateToProps = ({ account, billing }) => {
  const {
    profiles,
    paymentMethods,
  } = billing;
  const {
    selectedNamespace,
    namespaceObjects,
    currentUser,
  } = account;
  const user = namespaceObjects.results[selectedNamespace] || {};
  const billingProfile = profiles.results[user.id];
  return {
    billingProfile,
    billingPaymentMethods: paymentMethods.results,
    currentUser,
    selectedUser: user,
  };
};

const dispatcher = {
  accountToggleMagicCarpet,
  billingCreatePaymentMethod,
  billingCreateProfile,
  billingCreateSubscription,
  billingFetchPaymentMethods,
  billingSetDefaultPaymentMethod,
};

export default reduxForm({
  form: 'purchaseForm',
  fields,
  validate,
},
mapStateToProps,
dispatcher,
)(PurchaseDetails);
