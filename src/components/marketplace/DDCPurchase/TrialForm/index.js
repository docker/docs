import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import merge from 'lodash/merge';
import {
  Button,
  Card,
  Checkbox,
  Input,
  Select,
} from 'common';
import {
  countryOptions,
} from 'lib/constants/countries';
import { SUBSCRIPTIONS } from 'lib/constants/overlays';
import { EVAL_LINK } from 'lib/constants/eusa';
import {
  accountFetchUser,
  accountToggleMagicCarpet,
  accountSelectNamespace,
} from 'actions/account';
import {
  billingFetchProfile,
  billingFetchProfileSubscriptions,
  billingCreateProfile,
  billingCreateSubscription,
} from 'actions/billing';
import AccountSelect from 'components/marketplace/Purchase/AccountSelect';
import EmailSelect from 'components/marketplace/Purchase/EmailSelect';
import validate from './validations';
import css from './styles.css';

const { array, bool, object, string, func } = PropTypes;

class TrialForm extends Component {
  static propTypes = {
    accountEmails: array.isRequired,
    accountOptions: array.isRequired,
    billingProductDetails: object,
    billingProfile: object,
    currentUser: object.isRequired,
    DDCTrialExists: bool.isRequired,
    fields: object.isRequired,
    initialized: object.isRequired,
    location: object,
    fetchedNamespaces: object.isRequired,
    productDetails: object,
    profileError: string,
    // actions
    accountFetchUser: func.isRequired,
    accountSelectNamespace: func.isRequired,
    accountToggleMagicCarpet: func.isRequired,
    billingCreateProfile: func.isRequired,
    billingCreateSubscription: func.isRequired,
    billingFetchProfile: func.isRequired,
    billingFetchProfileSubscriptions: func.isRequired,
    handleSubmit: func.isRequired,
  }

  static contextTypes = {
    router: object.isRequired,
  }

  onCheck = () => {
    const { accepted } = this.props.fields;
    const checked = accepted.value === 'checked' ? 'unchecked' : 'checked';
    accepted.onChange(checked);
  }

  onSelectChange = (field) => (data) => {
    if (field === 'account') {
      const {
        accountFetchUser: fetchUser,
        accountSelectNamespace: selectNamespace,
        billingFetchProfile: fetchBillingProfiles,
        billingFetchProfileSubscriptions: fetchProfileSubscriptions,
        currentUser,
        fetchedNamespaces,
      } = this.props;
      if (fetchedNamespaces[data.value]) {
        const userOrg = fetchedNamespaces[data.value];
        const { type, id: docker_id } = userOrg;
        const isOrg = type === 'Organization';
        selectNamespace({ namespace: data.value });
        fetchBillingProfiles({ docker_id, isOrg });
        fetchProfileSubscriptions({ docker_id });
      } else {
        const isOrg = data.value !== currentUser.username;
        // fetching user info vs. org info locally hits different api's which
        // requires knowing which endpoint to hit.
        // unnecessary in production since v2/users will redirect to v2/orgs
        fetchUser({ namespace: data.value, isOrg }).then((userRes) => {
          const { id: docker_id } = userRes.value;
          selectNamespace({ namespace: data.value });
          fetchBillingProfiles({ docker_id, isOrg });
          fetchProfileSubscriptions({ docker_id });
        });
      }
    }
    const fieldObject = this.props.fields[field];
    fieldObject.onChange(data.value);
  }

  submitForm(values) {
    // IF SUBSCRIPTION EXISTS > DONT SHOW THE FORM :|
    const {
      billingProductDetails,
      billingProfile,
      DDCTrialExists,
      fetchedNamespaces,
      location,
      accountToggleMagicCarpet: toggleMagicCarpet,
      billingCreateProfile: createProfile,
      billingCreateSubscription: createSubscription,
    } = this.props;
    if (DDCTrialExists) {
      // If trial already exists prevent submission
      return;
    }
    const userObject = fetchedNamespaces[values.account];
    const auth = { docker_id: userObject.id };
    const profile = {
      first_name: values.firstName.trim(),
      last_name: values.lastName.trim(),
      email: values.email.trim(),
      phone_primary: values.phone.trim(),
      company_name: values.company.trim(),
      job_function: values.job.trim(),
    };
    const plan = location.query.plan;
    const date = new Date();
    const prodName = billingProductDetails.label || '';
    const subData = {
      name: `${prodName} Subscription ${date.toDateString()}`,
      product_rate_plan: plan,
      pricing_components: [],
      product_id: billingProductDetails.id,
      eusa: {
        accepted: values.accepted === 'checked',
      },
    };
    const router = this.context.router;
    // grabbing path from location since we cannot tell whether this is a bundle
    // or an image otherwise
    const detailRoute = location.pathname.split('/purchase')[0];
    if (!billingProfile || (billingProfile.error)) {
      const data = { ...auth, profile };
      createProfile(data).then(() => {
        createSubscription(merge({}, subData, auth)).then(
          () => {
            router.push(`${detailRoute}?overlay=subscriptions`);
            // redirect to details page on success with subscription carpet open
            toggleMagicCarpet({ magicCarpet: SUBSCRIPTIONS });
            analytics.track('ddc_trial_signup');
          }
        );
        const keys = userObject.type === 'User' ? {
          selectedUser: userObject.id,
          Docker_Hub_User_Name__c: userObject.username,
        } : {
          dockerUUID: userObject.id,
          Docker_Hub_Organization_Name__c: userObject.orgname,
        };

        analytics.identify(userObject.id, {
          ...keys,
          firstName: values.firstName.trim(),
          lastName: values.lastName.trim(),
          company: values.company.trim(),
          title: values.job.trim(),
          phone: values.phone.trim(),
          name: `${values.firstName.trim()} ${values.lastName.trim()}`,
          address: {
            country: values.country,
            state: values.state.trim(),
          },
        });
      });
    } else {
      createSubscription(merge({}, subData, auth)).then(
        () => {
          router.push(`${detailRoute}?overlay=subscriptions`);
          // redirect to details page on success with subscription carpet open
          toggleMagicCarpet({ magicCarpet: SUBSCRIPTIONS });
          analytics.track('ddc_trial_signup');
        }
      );
    }
  }

  renderForm() {
    const {
      accountEmails,
      DDCTrialExists,
      fields: propFields,
      handleSubmit,
      initialized,
      profileError,
    } = this.props;
    const {
      accepted,
      company,
      country,
      firstName,
      lastName,
      job,
      phone,
      state,
    } = propFields;
    return (
      <div className={css.contact}>
        <div className={css.sectionTitle}>
          Contact Information
        </div>
        <form
          onSubmit={handleSubmit(::this.submitForm)}
          className={css.form}
        >
          <div className={css.soft}>All fields required</div>
          <div className={css.row}>
            <Input
              {...firstName}
              readOnly={!!initialized.firstName}
              id={'first'}
              placeholder="First Name"
              errorText={firstName.touched && firstName.error}
            />
            <Input
              {...lastName}
              readOnly={!!initialized.lastName}
              id={'last'}
              placeholder="Last Name"
              errorText={lastName.touched && lastName.error}
            />
          </div>
          <div className={css.row}>
            <Input
              {...company}
              readOnly={!!initialized.company}
              id={'company'}
              placeholder="Company"
              errorText={company.touched && company.error}
            />
            <Input
              {...job}
              readOnly={!!initialized.job}
              id={'job'}
              placeholder="Job Title"
              errorText={job.touched && job.error}
            />
          </div>
          <div className={css.row}>
            <EmailSelect
              fields={propFields}
              initialized={initialized}
              onSelectChange={this.onSelectChange('email')}
              accountEmails={accountEmails}
            />
            <Input
              {...phone}
              readOnly={!!initialized.phone}
              id={'phone'}
              placeholder="Phone Number"
              errorText={phone.touched && phone.error}
            />
          </div>
          <div className={css.row}>
            <Input
              {...state}
              readOnly={!!initialized.state}
              id={'state'}
              placeholder="State/Province"
              errorText={state.touched && state.error}
            />
            <Select
              {...country}
              onBlur={() => {}}
              disabled={!!initialized.country}
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
          <div className={css.euta}>
            <Checkbox
              {...accepted}
              checked={accepted.value === 'checked'}
              onBlur={() => {}}
              onCheck={this.onCheck}
              style={{ width: '' }}
            />
            <div>
              I have read and agree to Docker's&nbsp;
              <a href={EVAL_LINK} >
                Software Evaluation Agreement
              </a>
            </div>
          </div>
          <div>
            <div className={css.error}>
              {(accepted.touched && accepted.error) || profileError}
            </div>
            <Button
              id="submit"
              disabled={DDCTrialExists}
              className={css.submit}
              type="submit"
              onClick={handleSubmit(::this.submitForm)}
            >
              Start your evaluation!
            </Button>
          </div>
        </form>
      </div>
    );
  }

  render() {
    const {
      accountOptions,
      DDCTrialExists,
      fields: propFields,
    } = this.props;

    const trialTitle =
      'Sign up for a free 30-day evaluation of Docker Datacenter';
    const trialSubtext = `
      This 30-day evaluation gives you complete access to our Docker Datacenter
      solution for a full fledged hands-on experience.
    `;
    let form = this.renderForm();
    if (DDCTrialExists) {
      form = (
        <Card className={css.subscriptionExists}>
          <div>
            Subscription already exists.
          </div>
        </Card>
      );
    }
    return (
      <Card
        className={css.card}
        title={trialTitle}
      >
        <div className={css.subtitle}>
          {trialSubtext}
        </div>
        <AccountSelect
          options={accountOptions}
          fields={propFields}
          onSelectChange={this.onSelectChange('account')}
        />
        {form}
      </Card>
    );
  }
}

const fields = [
  'accepted',
  'account',
  'company',
  'country',
  'email',
  'firstName',
  'job',
  'lastName',
  'phone',
  'state',
];

const mapStateToProps = ({ account, billing }, props) => {
  const {
    profiles,
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
  if (!!billingProfile) {
    const address =
      billingProfile.addresses && billingProfile.addresses[0] || {};
    initialized = {
      company: billingProfile.company_name,
      country: address.country,
      email: billingProfile.email,
      firstName: billingProfile.first_name,
      job: billingProfile.job_function,
      lastName: billingProfile.last_name,
      phone: billingProfile.phone_primary,
      state: address.province,
      accepted: 'unchecked',
    };
    merge(initialData, initialized);
  }
  return {
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
  accountToggleMagicCarpet,
  billingCreateProfile,
  billingCreateSubscription,
  billingFetchProfile,
  billingFetchProfileSubscriptions,
};

export default reduxForm({
  form: 'trialForm',
  fields,
  validate,
},
mapStateToProps,
dispatcher,
)(TrialForm);
