import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import map from 'lodash/map';
// Blob is a polyfill to support older browsers
require('lib/utils/blob');
import { saveAs } from 'lib/utils/file-saver';
import classnames from 'classnames';

import {
  Card,
  DownloadIcon,
  FetchingError,
  FullscreenLoading,
} from '../lib/common';
import { billingFetchInvoicePDF } from 'actions/billing';
import css from './styles.css';
const PENDING = 'Pending';

const { array, bool, func, object, shape, string } = PropTypes;

const mapStateToProps = ({ billing }) => {
  const { invoices } = billing;
  return {
    invoices,
  };
};

const dispatcher = {
  fetchInvoicePDF: billingFetchInvoicePDF,
};

@connect(mapStateToProps, dispatcher)
export default class InvoicesView extends Component {
  static propTypes = {
    invoices: shape({
      error: string,
      isFetching: bool.isRequired,
      results: array.isRequired,
    }).isRequired,
    fetchInvoicePDF: func.isRequired,
    selectedUser: object.isRequired,
  }
  // When an invoice is being downloaded, there will be a key for invoice id
  // with status information: [invoiceId]: { isDownloading: bool, error: bool }
  state = {}

  downloadInvoice = ({ invoice_id, issuedDate }) => () => {
    const { fetchInvoicePDF, selectedUser } = this.props;
    const { id: docker_id } = selectedUser;
    // Make sure to not override any existing invoice download statuses
    this.setState({
      [invoice_id]: {
        isDownloading: true,
        error: false,
      },
    }, () => {
      fetchInvoicePDF({ docker_id, invoice_id })
        .then((res) => {
          // value contains the API response
          const downloadContent = res.value;
          this.setState({
            [invoice_id]: {
              isDownloading: false,
              error: false,
            },
          });
          const blob = new Blob([downloadContent], { type: 'application/pdf' });
          saveAs(blob, `Docker Invoice ${issuedDate}.pdf`);
        })
        .catch(() => {
          this.setState({
            [invoice_id]: {
              isDownloading: false,
              error: true,
            },
          });
        });
    });
  };

  renderInvoiceRow = (invoice) => {
    const {
      bf_invoice_id: invoice_id,
      cost,
      currency,
      issued,
      payment_received,
      state,
      subscriptions,
    } = invoice;
    // Do not render invoices in the 'Pending' state
    if (state === PENDING) {
      return null;
    }
    const { isDownloading, error } = this.state[invoice_id] || {};
    const issuedDate = new Date(issued).toDateString();
    const paid = new Date(payment_received).toDateString();
    const costAndCurrency = `$${cost} ${currency}`;
    const subscriptionsList = map(subscriptions, (sub) => {
      return (<div key={sub.subscription_id}>{sub.subscription_name}</div>);
    });
    let downloadText = 'PDF';
    if (isDownloading && !error) {
      downloadText = 'Loading...';
    } else if (isDownloading && error) {
      downloadText = 'Error';
    }
    const classes = classnames({
      [css.downloadLicense]: true,
      [css.downloadLicenseError]: isDownloading && error,
    });
    return (
      <div key={invoice_id} className={css.row}>
        <div>{issuedDate}</div>
        <div>{subscriptionsList}</div>
        <div>{costAndCurrency}</div>
        <div>{paid}</div>
        <div
          className={classes}
          onClick={this.downloadInvoice({ invoice_id, issuedDate })}
        >
          <DownloadIcon />
          {downloadText}
        </div>
      </div>
    );
  }

  render() {
    const { isFetching, error, results } = this.props.invoices;
    if (isFetching) {
      return <FullscreenLoading />;
    } else if (error) {
      return (
        <div className={css.fetchingError}>
          <FetchingError resource="your invoices" />
        </div>
      );
    }
    let content;
    if (results.length < 1) {
      content = <div className={css.noInvoices}>No invoices</div>;
    } else {
      content = map(results, this.renderInvoiceRow);
    }
    return (
      <Card title="Invoices" shadow>
        <div className={css.table}>
          <div className={css.head}>
            <div>Date</div>
            <div>Plans & Subscriptions</div>
            <div>Amount</div>
            <div>Paid</div>
            <div></div>
          </div>
          {content}
        </div>
      </Card>
    );
  }
}
