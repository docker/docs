import React, { Component, PropTypes } from 'react';
import {
  Card,
  Input,
} from 'components/common';
import { FALLBACK_ELEMENT } from 'lib/constants/fallbacks';

import css from './styles.css';

const { object } = PropTypes;

export default class PublisherDetailsForm extends Component {
  static propTypes = {
    company: object.isRequired,
    companyLogo: object.isRequired,
    companyWebsite: object.isRequired,
  }

  render() {
    const {
      company,
      companyLogo: logo,
      companyWebsite: website,
    } = this.props;

    const iconPreview = logo.value ? (
      <img
        alt="Product Icon"
        className={css.iconPreview}
        src={logo.value}
      />
    ) : FALLBACK_ELEMENT;

    return (
      <Card className={css.card}>
        <div className={css.title}>
          Publisher Details
        </div>
        <div className={css.publisherDetailsWrapper}>
          <div className={css.publisherLogo}>
            <div className={css.companyLogo}>
              {iconPreview}
            </div>
          </div>
          <div>
            <div className={css.fields}>
              <div className={css.sectionTitle}>Company Name</div>
              <div className={css.subText}>
                Your company name as it will appear in the Docker Store
              </div>
              <Input
                className={css.input}
                id="product-names-input"
                errorText={company.touched && company.error || ''}
                { ...company }
              />
            </div>
            <div className={css.fields}>
              <div className={css.sectionTitle}>Company Logo URL</div>
              <Input
                className={css.input}
                id="product-names-input"
                errorText={logo.touched && logo.error || ''}
                { ...logo }
              />
            </div>
            <div className={css.fields}>
              <div className={css.sectionTitle}>Company Website</div>
              <Input
                className={css.input}
                id="product-names-input"
                errorText={website.touched && website.error || ''}
                { ...website }
              />
            </div>
          </div>
        </div>
      </Card>
    );
  }
}
