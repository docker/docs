'use strict';
import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import _ from 'lodash';
import moment from 'moment';
import { FullSection } from '../../common/Sections.jsx';
import { FlexTable, FlexRow, FlexHeader, FlexItem } from '../../common/FlexTable.jsx';
import downloadInvoice from '../../../actions/downloadInvoice.js';
const debug = require('debug')('COMPONENT:INVOICE_TABLE');

let mkInvoiceTable = function(invoice) {
  var subtotal = '$' + (parseInt(invoice.subtotal_in_cents, 10) / 100);
  var total = '$' + (parseInt(invoice.total_in_cents, 10) / 100);
  var date = moment.utc(invoice.created_at).format('MMM Do YYYY');
  debug('CREATED_AT: ', invoice.created_at);
  debug('CREATED_AT_MOMENT: ', date);
  return (
    <FlexRow key={invoice.invoice_number}>
      <FlexItem>{date}</FlexItem>
      <FlexItem>{invoice.invoice_number}</FlexItem>
      <FlexItem>{invoice.state}</FlexItem>
      <FlexItem>{total}</FlexItem>
      <FlexItem>
        <a href="#"
           onClick={this.downloadInvoice(invoice.invoice_number)}
           data-invoice={invoice.id}>Download Invoice</a>
      </FlexItem>
    </FlexRow>
  );
};

var InvoiceTables = React.createClass({
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  propTypes: {
    invoices: PropTypes.array,
    username: PropTypes.string,
    JWT: PropTypes.string
  },
  downloadInvoice(id) {
    return (e) => {
      e.preventDefault();
      this.context.executeAction(downloadInvoice, {JWT: this.props.JWT, username: this.props.username, invoiceId: id});
    };
  },
  render: function() {
    if (!this.props.JWT || _.isEmpty(this.props.invoices)) {
      return (
        <div></div>
      );
    } else {
      return (
        <FullSection title='Invoice'>
          <div className="columns large-12">
            <FlexTable>
              <FlexHeader>
                <FlexItem>Date</FlexItem>
                <FlexItem>Invoice #</FlexItem>
                <FlexItem>State</FlexItem>
                <FlexItem>Total</FlexItem>
                <FlexItem>Download</FlexItem>
              </FlexHeader>
              {this.props.invoices.map(mkInvoiceTable, this)}
            </FlexTable>
          </div>
        </FullSection>
      );
    }
  }
});

module.exports = InvoiceTables;
