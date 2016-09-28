import billingReducer, { DEFAULT_STATE } from 'reducers/billing';
import {
  BILLING_CREATE_PAYMENT_METHOD,
  BILLING_CREATE_PAYMENT_TOKEN,
  BILLING_CREATE_PROFILE,
  BILLING_CREATE_SUBSCRIPTION,
  BILLING_DELETE_SUBSCRIPTION,
  BILLING_FETCH_INVOICES,
  BILLING_FETCH_LICENSE_DETAIL,
  BILLING_FETCH_PAYMENT_METHODS,
  BILLING_FETCH_PRODUCT,
  BILLING_FETCH_PROFILE_SUBSCRIPTIONS,
  BILLING_FETCH_PROFILE,
  BILLING_UPDATE_PROFILE,
  BILLING_UPDATE_SUBSCRIPTION,
} from 'actions/billing';
import { expect, assert } from 'chai';
import merge from 'lodash/merge';

describe('billing reducer', () => {
  const subscription_id = 'sub-1234';

  it('should return the initial state', () => {
    expect(billingReducer(undefined, {})).to.deep.equal(DEFAULT_STATE);
  });

  //----------------------------------------------------------------------------
  // BILLING_CREATE_PROFILE
  //----------------------------------------------------------------------------
  it('should handle BILLING_CREATE_PROFILE_REQ', () => {
    const payload = {};
    const action = {
      type: `${BILLING_CREATE_PROFILE}_REQ`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.profiles.isFetching).to.be.true;
    expect(reducer.profiles.isSubmitting).to.be.true;
  });
  it('should handle BILLING_CREATE_PROFILE_ACK', () => {
    const docker_id = 'test123';
    const payload = {
      profile: {
        email: 'matt@docker.com',
        first_name: 'Matt',
        last_name: 'Tescher',
        addresses: [
          {
            address_line_1: 'address line 1',
            address_line_2: 'address line 2',
            address_line_3: 'address line 3',
            city: 'SF',
            province: 'CA',
            country: 'United States',
            post_code: '94107',
            primary_address: true,
          },
        ],
        phone_primary: '1234567890',
        job_function: 'Developer',
        company_name: 'Docker, Inc',
      },
    };
    const action = {
      type: `${BILLING_CREATE_PROFILE}_ACK`,
      meta: { docker_id },
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.profiles.results[docker_id]).to.equal(payload.profile);
    expect(reducer.profiles.isFetching).to.be.false;
    expect(reducer.profiles.isSubmitting).to.be.false;
  });
  it('should handle BILLING_CREATE_PROFILE_ERR', () => {
    const payload = {};
    const action = {
      type: `${BILLING_CREATE_PROFILE}_ERR`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.profiles.isFetching).to.be.false;
    expect(reducer.profiles.isSubmitting).to.be.false;
    expect(reducer.profiles.error).to.exist;
  });

  //----------------------------------------------------------------------------
  // BILLING_UPDATE_PROFILE
  //----------------------------------------------------------------------------

  it('should handle BILLING_UPDATE_PROFILE_REQ', () => {
    const payload = {};
    const action = {
      type: `${BILLING_UPDATE_PROFILE}_REQ`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.profiles.isFetching).to.be.true;
    expect(reducer.profiles.isSubmitting).to.be.true;
  });

  it('should handle BILLING_UPDATE_PROFILE_ACK', () => {
    const docker_id = 'test123';
    const payload = {
      email: 'matt@docker.com',
      first_name: 'Matt',
      last_name: 'Tescher',
      addresses: [
        {
          address_line_1: 'address line 1',
          address_line_2: 'address line 2',
          address_line_3: 'address line 3',
          city: 'SF',
          province: 'CA',
          country: 'United States',
          post_code: '94107',
          primary_address: true,
        },
      ],
      phone_primary: '1234567890',
      job_function: 'Developer',
      company_name: 'Docker, Inc',
    };
    const action = {
      type: `${BILLING_UPDATE_PROFILE}_ACK`,
      meta: { docker_id },
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.profiles.results[docker_id]).to.equal(payload.profile);
    expect(reducer.profiles.isFetching).to.be.false;
    expect(reducer.profiles.isSubmitting).to.be.false;
  });

  it('should handle BILLING_UPDATE_PROFILE_ERR', () => {
    const payload = {};
    const action = {
      type: `${BILLING_UPDATE_PROFILE}_ERR`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.profiles.isFetching).to.be.false;
    expect(reducer.profiles.isSubmitting).to.be.false;
    expect(reducer.profiles.error).to.exist;
  });

  //----------------------------------------------------------------------------
  // BILLING_CREATE_SUBSCRIPTION
  //----------------------------------------------------------------------------
  it('should handle BILLING_CREATE_SUBSCRIPTION_REQ', () => {
    const payload = {};
    const action = {
      type: `${BILLING_CREATE_SUBSCRIPTION}_REQ`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.isFetching).to.be.true;
    expect(reducer.subscriptions.isSubmitting).to.be.true;
  });

  it('should handle BILLING_CREATE_SUBSCRIPTION_ACK', () => {
    const payload = { subscription_id };
    const action = {
      type: `${BILLING_CREATE_SUBSCRIPTION}_ACK`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.isFetching).to.be.false;
    expect(reducer.subscriptions.isSubmitting).to.be.false;
    expect(reducer.subscriptions.results[subscription_id]).to.equal(payload);
  });

  it('should handle BILLING_CREATE_SUBSCRIPTION_ERR', () => {
    const payload = { subscription_id };
    const action = {
      type: `${BILLING_CREATE_SUBSCRIPTION}_ERR`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.isFetching).to.be.false;
    expect(reducer.subscriptions.isSubmitting).to.be.false;
    expect(reducer.subscriptions.error).to.exist;
  });

  //----------------------------------------------------------------------------
  // BILLING_FETCH_INVOICES
  //----------------------------------------------------------------------------
  it('should handle BILLING_FETCH_INVOICES_REQ', () => {
    const payload = [];
    const action = {
      type: `${BILLING_FETCH_INVOICES}_REQ`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.invoices.isFetching).to.be.true;
    expect(reducer.invoices.error).to.be.empty;
  });

  it('should handle BILLING_FETCH_INVOICES_ACK', () => {
    const payload = [
      {
        cost: 70,
        currency: 'USD',
        invoice_id: 'INV-351969A6-2A8A-4C31-BB0C-1FEF5BF1',
        issued: '2016-05-15T16:07:52Z',
        paid: 70,
        payment_received: '2016-05-15T16:23:08Z',
        period_end: '2016-06-15T16:07:52Z',
        period_start: '2016-05-15T16:07:52Z',
        state: 'Paid',
        subscription_id: 'docker-sub-id',
        invoice_lines: [
          {
            cost: 70,
            product_id: 'docker-prod-id',
            product: 'Docker Hub',
            product_rate_plan: 'Docker Hub',
            pricing_component: {
              name: 'Private Repos',
              value: 10,
            },
          },
        ],
      },
    ];
    const action = {
      type: `${BILLING_FETCH_INVOICES}_ACK`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.invoices.isFetching).to.be.false;
    expect(reducer.invoices.error).to.be.empty;
    expect(reducer.invoices.results).to.have.lengthOf(1);
    expect(reducer.invoices.results).to.include(payload[0]);
  });

  it('should handle BILLING_FETCH_INVOICES_ERR', () => {
    const payload = [];
    const action = {
      type: `${BILLING_FETCH_INVOICES}_ERR`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.invoices.isFetching).to.be.false;
    expect(reducer.invoices.error).to.not.be.empty;
  });

  //----------------------------------------------------------------------------
  // BILLING_FETCH_PROFILE
  //----------------------------------------------------------------------------
  it('should handle BILLING_FETCH_PROFILES_REQ', () => {
    const payload = {};
    const reqAction = {
      type: `${BILLING_FETCH_PROFILE}_REQ`,
      payload,
    };
    const reducer = billingReducer(undefined, reqAction);
    expect(reducer.profiles.isFetching).to.be.true;
    expect(reducer.paymentMethods.isFetching).to.be.false;
    expect(reducer.paymentMethods.results).to.be.empty;
    expect(reducer.paymentMethods.error).to.be.empty;
  });

  it('should handle BILLING_FETCH_PROFILES_ACK', () => {
    const payload = {
      docker_id: 'acct-l337',
      profile: {
        email: 'matt@docker.com',
        first_name: 'acct-l337',
        last_name: 'acct-l337',
        addresses: [
          {
            address_line_1: 'address line 1',
            address_line_2: 'address line 2',
            address_line_3: 'address line 3',
            city: 'SF',
            province: 'CA',
            country: 'United States',
            post_code: '94107',
            primary_address: true,
          },
        ],
        company_name: 'acct-l337',
      },
    };
    const ackAction = {
      type: `${BILLING_FETCH_PROFILE}_ACK`,
      meta: { docker_id: payload.docker_id },
      payload,
    };
    let reducer = billingReducer(undefined, ackAction);
    assert.isObject(reducer.profiles.results);
    assert.property(reducer.profiles.results, payload.docker_id);
    // eslint-disable-next-line
    expect(reducer.profiles.results[payload.docker_id]).to.deep.equal(payload.profile);

    const errAction = {
      type: `${BILLING_FETCH_PROFILE}_ERR`,
      meta: { docker_id: payload.docker_id },
      payload: '',
    };
    reducer = billingReducer(undefined, errAction);
    assert.property(reducer.profiles, 'error');
    expect(reducer.profiles.error).to.be.string;
  });

  //----------------------------------------------------------------------------
  // BILLING_FETCH_PRODUCT
  //----------------------------------------------------------------------------
  it('should handle BILLING_FETCH_PRODUCT', () => {
    const id = 'billing_fetch_product_test_id';
    const payload = {
      id,
      publisher_id: '<publisher_id>',
      name: 'ddc',
      label: 'Docker Data Center',
      description: '',
      rate_plans: [
        {
          name: 'trial',
          label: 'Trial',
          duration: 1,
          duration_period: 'months',
          trial: 0,
          trial_period: 'months',
          currency: 'USD',
          pricing_components: [
            {
              name: 'docker-engines',
              label: 'Docker Engines',
              charge_type: 'subscription',
              charge_model: 'tiered',
              tiers: [
                {
                  lower_threshold: 1,
                  upper_threshold: 10000,
                  pricing_type: 'unit',
                  price: 150.00000,
                },
              ],
            },
          ],
        },
      ],
    };
    const action = {
      type: `${BILLING_FETCH_PRODUCT}_ACK`,
      meta: { id },
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.products[payload.id])
      .to.deep.equal({ ...payload, isFetching: false });
  });

  //----------------------------------------------------------------------------
  // BILLING_DELETE_SUBSCRIPTION
  //----------------------------------------------------------------------------
  it('should handle BILLING_DELETE_SUBSCRIPTION_REQ', () => {
    const action = {
      type: `${BILLING_DELETE_SUBSCRIPTION}_REQ`,
      meta: { subscription_id },
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.delete.isDeleting).to.equal(true);
    expect(reducer.subscriptions.delete.subscription_id)
      .to.equal(subscription_id);
    expect(reducer.subscriptions.delete.error).to.equal('');
  });

  it('should handle BILLING_DELETE_SUBSCRIPTION_ACK', () => {
    const action = {
      type: `${BILLING_DELETE_SUBSCRIPTION}_ACK`,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.delete.isDeleting).to.equal(false);
    expect(reducer.subscriptions.delete.subscription_id).to.equal('');
    expect(reducer.subscriptions.delete.error).to.equal('');
  });

  it('should handle BILLING_DELETE_SUBSCRIPTION_ERR', () => {
    const action = {
      type: `${BILLING_DELETE_SUBSCRIPTION}_ERR`,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.delete.isDeleting).to.equal(false);
    expect(reducer.subscriptions.delete.error).to.exist;
  });

  //----------------------------------------------------------------------------
  // BILLING_UPDATE_SUBSCRIPTION
  //----------------------------------------------------------------------------
  it('should handle BILLING_UPDATE_SUBSCRIPTION_REQ', () => {
    const action = {
      type: `${BILLING_UPDATE_SUBSCRIPTION}_REQ`,
      meta: { subscription_id },
    };
    const reducer = billingReducer(undefined, action);
    const { meta } = action;
    expect(reducer.subscriptions.update[meta.subscription_id].isUpdating)
      .to.equal(true);
    expect(reducer.subscriptions.update[meta.subscription_id].error)
      .to.equal('');
  });

  it('should handle BILLING_UPDATE_SUBSCRIPTION_ACK', () => {
    const action = {
      type: `${BILLING_UPDATE_SUBSCRIPTION}_ACK`,
      meta: { subscription_id },
    };
    const REQ_STATE = merge({}, DEFAULT_STATE, {
      subscriptions: {
        update: {
          [subscription_id]: {
            isUpdating: true,
            error: '',
          },
        },
      },
    });
    const reducer = billingReducer(REQ_STATE, action);
    const { meta } = action;
    expect(reducer.subscriptions.update[meta.subscription_id].isUpdating)
      .to.equal(false);
    expect(reducer.subscriptions.update[meta.subscription_id].error)
      .to.equal('');
  });

  it('should handle BILLING_UPDATE_SUBSCRIPTION_ERR', () => {
    const action = {
      type: `${BILLING_UPDATE_SUBSCRIPTION}_ERR`,
      meta: { subscription_id },
    };
    const REQ_STATE = merge({}, DEFAULT_STATE, {
      subscriptions: {
        update: {
          [subscription_id]: {
            isUpdating: true,
            error: '',
          },
        },
      },
    });
    const reducer = billingReducer(REQ_STATE, action);
    const { meta } = action;
    expect(reducer.subscriptions.update[meta.subscription_id].isUpdating)
      .to.equal(false);
    expect(reducer.subscriptions.update[meta.subscription_id].error).to.exist;
  });


  //----------------------------------------------------------------------------
  // BILLING_FETCH_PROFILE_SUBSCRIPTIONS
  //----------------------------------------------------------------------------
  it('should handle BILLING_FETCH_PROFILE_SUBSCRIPTIONS_REQ', () => {
    const action = {
      type: `${BILLING_FETCH_PROFILE_SUBSCRIPTIONS}_REQ`,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.isFetching).to.be.true;
  });

  it('should handle BILLING_FETCH_PROFILE_SUBSCRIPTIONS_ACK', () => {
    const subscription_1 = { subscription_id };
    const payload = [subscription_1];
    const action = {
      type: `${BILLING_FETCH_PROFILE_SUBSCRIPTIONS}_ACK`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.isFetching).to.be.false;
    expect(reducer.subscriptions.results[subscription_id])
      .to.equal(subscription_1);
  });

  it('should handle BILLING_FETCH_PROFILE_SUBSCRIPTIONS_ERR', () => {
    const err = { error: 'this is an error' };
    const action = {
      type: `${BILLING_FETCH_PROFILE_SUBSCRIPTIONS}_ERR`,
      payload: err,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.isFetching).to.false;
    expect(reducer.subscriptions.error).to.exist;
  });

  //----------------------------------------------------------------------------
  // BILLING_FETCH_LICENSE_DETAIL
  //----------------------------------------------------------------------------
  it('should handle BILLING_FETCH_LICENSE_DETAIL_REQ', () => {
    const action = {
      type: `${BILLING_FETCH_LICENSE_DETAIL}_REQ`,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.licenses.isFetching).to.be.true;
  });

  it('should handle BILLING_FETCH_LICENSE_DETAIL_ACK', () => {
    const payload = { key: '1', expires: 'sometime' };
    const action = {
      type: `${BILLING_FETCH_LICENSE_DETAIL}_ACK`,
      meta: { subscription_id },
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.licenses.isFetching).to.be.false;
    expect(reducer.subscriptions.licenses.results[subscription_id])
      .to.deep.equal(payload);
  });

  it('should handle BILLING_FETCH_LICENSE_DETAIL_ERR', () => {
    const status = '404';
    const action = {
      type: `${BILLING_FETCH_LICENSE_DETAIL}_ERR`,
      payload: {
        response: {
          error: {
            status,
          },
        },
      },
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.subscriptions.licenses.isFetching).to.false;
    expect(reducer.subscriptions.licenses.error).to.equal(status);
  });

  //----------------------------------------------------------------------------
  // BILLING_FETCH_PAYMENT_METHODS
  //----------------------------------------------------------------------------
  it('should handle BILLING_FETCH_PAYMENT_METHODS_ACK', () => {
    const payload = [
      {
        payment_method_id: 'pmd-1337',
        description: '############4242',
        expiry_month: 6,
        expiry_year: 21,
        card_type: 'Visa',
        country: 'US',
        default: true,
      },
    ];
    const action = {
      type: `${BILLING_FETCH_PAYMENT_METHODS}_ACK`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.paymentMethods.isFetching).to.be.false;
    expect(reducer.paymentMethods.results).to.equal(payload);
  });

  it('should handle BILLING_FETCH_PAYMENT_METHODS_REQ', () => {
    const payload = {};
    const action = {
      type: `${BILLING_FETCH_PAYMENT_METHODS}_REQ`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.paymentMethods.isFetching).to.be.true;
  });

  //----------------------------------------------------------------------------
  // BILLING_CREATE_PAYMENT_METHOD
  //----------------------------------------------------------------------------
  it('should handle BILLING_CREATE_PAYMENT_METHOD_ACK', () => {
    const payload = [{
      '@type': 'paymentMethod',
      created: '2016-05-26T21:01:03Z',
      changedBy: '06948B82-98DE-4C46-8B84-C1E4E0C4F0DE',
      updated: '2016-05-26T21:01:03Z',
      id: 'PMD-6B92029C-EEAF-4B44-B9B2-6A903A2A',
      crmID: 'card_18FVb5CMIduy6ycT46kf7pNc',
      accountID: 'ACC-A4F26D67-F523-4833-A9E9-3BFABD02',
      organizationID: 'ORG-7B018331-80B8-44E8-B5E1-00407840',
      name: 'Visa (4242)',
      description: '############4242',
      cardHolderName: 'billingPaymentMethods billingPaymentMethods',
      expiryDate: '03/2019',
      cardType: 'Visa',
      country: 'US',
      lastFour: '4242',
      expiryYear: 19,
      expiryMonth: 3,
      linkID: 'DFE13B1C-272D-4A44-A0B6-74167ADC5E3A',
      gateway: 'stripe',
      state: 'Active',
      deleted: false,
      defaultPaymentMethod: false,
    }];
    const action = {
      type: `${BILLING_CREATE_PAYMENT_METHOD}_ACK`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.paymentMethods.isSubmitting).to.be.false;
  });

  it('should handle BILLING_CREATE_PAYMENT_METHOD_ERR', () => {
    const payload = [];
    const action = {
      type: `${BILLING_CREATE_PAYMENT_METHOD}_ERR`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.paymentMethods.isSubmitting).to.be.false;
    expect(reducer.paymentMethods.error).to.be.object;
    assert.property(reducer.paymentMethods.error, 'message');
  });

  it('should handle BILLING_CREATE_PAYMENT_TOKEN_REQ', () => {
    const payload = [];
    const action = {
      type: `${BILLING_CREATE_PAYMENT_TOKEN}_REQ`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.paymentMethods.isSubmitting).to.be.true;
  });

  it('should handle BILLING_CREATE_PAYMENT_TOKEN_ACK', () => {
    const payload = [];
    const action = {
      type: `${BILLING_CREATE_PAYMENT_TOKEN}_ACK`,
      payload,
    };
    const reducer = billingReducer(undefined, action);
    expect(reducer.paymentMethods.error).to.be.empty;
  });

  it('should handle BILLING_CREATE_PAYMENT_TOKEN_ERR', () => {
    let payload = {
      response: {
        body: {
          error: {
            type: 'card_error',
            message: 'The card number is not a valid credit card number.',
            code: 'invalid_number',
          },
        },
      },
    };
    let action = {
      type: `${BILLING_CREATE_PAYMENT_TOKEN}_ERR`,
      payload,
    };
    let reducer = billingReducer(undefined, action);
    expect(reducer.paymentMethods.isSubmitting).to.be.false;
    expect(reducer.paymentMethods.error).to.equal(payload.response.body.error);
    payload = {
      response: {
        body: {
          error: {
            type: 'authentication_error',
            message: 'The card number is not a valid credit card number.',
            code: 'invalid_number',
          },
        },
      },
    };
    action = {
      type: `${BILLING_CREATE_PAYMENT_TOKEN}_ERR`,
      payload,
    };
    reducer = billingReducer(undefined, action);
    expect(reducer.paymentMethods.isSubmitting).to.be.false;
    expect(reducer.paymentMethods.error).to.be.object;
    assert.property(reducer.paymentMethods.error, 'message');
  });
});
