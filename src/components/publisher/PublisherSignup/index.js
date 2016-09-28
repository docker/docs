import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { Card } from 'common';
import { CheckIcon, CloseIcon } from 'common/Icon';
import { publishSubscribe, publishGetSignup } from 'actions/publish';
import { XLARGE } from 'lib/constants/sizes';
import PublisherSignupForm from './PublisherSignupForm';

import css from './styles.css';

const { string, func, object } = PropTypes;

const PENDING = 'pending';

const TIPS = [
  {
    title: 'Reach More Customers',
    description: `Distribute software to Docker\'s fast-growing customer base.
      Customers discover, install and purchase your software directly from
      the Docker Store.`,
  },
  {
    title: 'Payments & Licensing Built-In',
    description: `We handle checkout, licensing and invoicing. Simply provide
      your content and we'll take it from there.`,
  },
  {
    title: 'Seamless & Secure Updates',
    description: `Deliver your software reliably and securely using containers.
      Customers are directly notified of updates and upgrades.`,
  },
];

const mapStateToProps = ({ account, publish }) => {
  const primaryEmails = account.userEmails &&
    account.userEmails.results &&
    account.userEmails.results.filter(r => r.primary);
  let primaryEmail;
  if (primaryEmails && primaryEmails.length > 0) {
    primaryEmail = primaryEmails[0].email;
  }

  let status = 'none';
  if (publish.signup && publish.signup.results) {
    status = publish.signup.results.status;
  }

  return {
    user: account.currentUser || {},
    email: primaryEmail,
    status,
  };
};

const mapDispatch = {
  publishGetSignup,
  publishSubscribe,
};

@connect(mapStateToProps, mapDispatch)
export default class PublisherSignup extends Component {
  static propTypes = {
    user: object,
    email: string,
    publishGetSignup: func,
    publishSubscribe: func,
    status: string,
  }

  state = {
    submitting: false,
    error: null,
  }

  onSubmit = values => {
    const {
      publishGetSignup: getSignup,
      publishSubscribe: subscribePublisher,
    } = this.props;
    this.setState({ submitting: false, error: null });
    const submitValues = {
      email: this.props.email,
      ...values,
    };

    subscribePublisher(submitValues).then(() => {
      this.setState({ submitting: false });
      getSignup();
      analytics.identify(this.props.user.id, {
        firstName: values.first_name,
        lastName: values.last_name,
        company: values.company,
      }, () => {
        analytics.track('Signed Up to Become a Publisher');
      });
    }).catch(err => {
      this.setState({ submitting: false, error: err.message });
    });
  };

  renderError() {
    return (
      <div>
        <div className={css.confirmed}>
          <CloseIcon size={XLARGE} className={css.error} />
        </div>
        <div className={css.confirmedMessage}>
          {'Could not fetch your Docker ID information.'}<br />
          {'Please check back soon.'}
        </div>
      </div>
    );
  }


  render() {
    const error = !this.props.user || !this.props.email ?
      this.renderError() : null;
    const form = this.props.status === PENDING ? (
      <div>
        <div className={css.confirmed}>
          <CheckIcon size={XLARGE} className={css.check} />
        </div>
        <div className={css.confirmedMessage}>
          {'You\'re on the list to become a Publisher.'}<br />
          {'We\'ll be in touch with you soon.'}
        </div>
      </div>
    ) : (
      <PublisherSignupForm
        onSubmit={this.onSubmit}
        submitting={this.state.submitting}
      />
    );

    return (
      <div>
        <div className={css.enroll}>
          <Card className={css.card}>
            <div className={css.enrollForm}>
              <h1>Become a Docker Store Publisher</h1>
              <h3>
                {`Enter your details and we\'ll get back to you at
                  ${this.props.email}`}
              </h3>
              {error || form}
            </div>
          </Card>
          <div className={css.tips}>
            {TIPS.map(tip => (
              <div key={tip.title} className={css.tip}>
                <h5>{tip.title}</h5>
                <p>{tip.description}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    );
  }
}
