import React, { Component, PropTypes } from 'react';
import {
  publishAcceptVendorAgreement,
  publishFetchProductDetails,
  publishCreateProduct,
  publishUpdateProductRepos,
  publishUpdateProductDetails,
} from 'actions/publish';
import PublisherSubmitSourceForm from './PublisherSubmitSourceForm';
import { connect } from 'react-redux';
import classnames from 'classnames';

import { STEP1 } from 'lib/constants/publisherSteps';
import css from './styles.css';

const TIPS = [
  {
    title: 'Use a private repository',
    description: `Make sure your product is available in a private repository in
    Docker Cloud or Docker Hub. We need to know the repository name, namespace,
    and tag.`,
  },
  {
    title: 'Submit the product for review',
    description: `Once you select a repository and tag, and provide information
      about the product, submit it for review. We'll contact you if you need to
      make any changes. Once a product is approved, you can decide when to make
      it available on the Docker Store by publishing it.`,
  },
  {
    title: 'Docker Security Scanning',
    description: `All product submissions to the Docker Store are scanned for
      security vulnerabilities using Docker Security Scanning. Only products
      without vulnerabilities are eligible to be published.`,
  },
];

const { func, object } = PropTypes;

const mapStateToProps = ({ publish }) => {
  const currentProductDetails = publish.currentProductDetails.results;
  return {
    currentProductDetails,
  };
};

const mapDispatch = {
  publishAcceptVendorAgreement,
  publishFetchProductDetails,
  publishCreateProduct,
  publishUpdateProductDetails,
  publishUpdateProductRepos,
};

@connect(mapStateToProps, mapDispatch)
export default class PublisherProductSource extends Component {
  static propTypes = {
    currentProductDetails: object.isRequired,
    onNextStep: func.isRequired,
    publishAcceptVendorAgreement: func.isRequired,
    publishCreateProduct: func.isRequired,
    publishFetchProductDetails: func.isRequired,
    publishUpdateProductDetails: func.isRequired,
    publishUpdateProductRepos: func.isRequired,
  }

  onSubmit = values => {
    const {
      currentProductDetails,
      onNextStep,
      publishAcceptVendorAgreement: acceptVendorAgreement,
      publishCreateProduct: createProduct,
      publishFetchProductDetails: fetchProductDetails,
      publishUpdateProductDetails: updateProductDetails,
      publishUpdateProductRepos: updateProductRepos,
    } = this.props;
    const {
      display_name,
      repoSources,
    } = values;
    // NOTE: If the product exists and we're going back to edit it.
    /*
    NOTE: A product could have been created in an offline manner on the
    Publisher portal, in which case the vendor agreement might not have gotten
    accepted. Ensure vendor agreement gets accepted (even if that happens
    redundantly in the initial source submission page).
    */
    if (currentProductDetails.id) {
      const {
        categories,
        full_description,
        id: product_id,
        links,
        logos,
        platforms,
        product_type,
        screenshots,
        short_description,
      } = currentProductDetails;
      const promises = [
        updateProductRepos({
          repoSources,
          product_id,
        }),
        updateProductDetails({
          product_id,
          details: {
            name: display_name,
            status: STEP1,
            product_type,
            categories,
            platforms,
            links,
            screenshots,
            full_description,
            short_description,
            logos,
          },
        }),
      ];
      // TODO: nathan 08/17 - refactor to only accept if not accepted;
      return acceptVendorAgreement().then(() => {
        return Promise.all(promises).then(() => {
          return fetchProductDetails({ product_id });
        }).then((res) => {
          const { id, status: productStatus } = res.value;
          // NOTE: onNextStep sends you to the next step in the publisher flow
          onNextStep({ id, productStatus });
        });
      });
    }
    return acceptVendorAgreement().then(() => {
      return createProduct({
        name: display_name,
        status: STEP1,
        repositories: repoSources,
      }).then((res) => {
        const { product_id } = res.value;
        return fetchProductDetails({ product_id });
      }).then((res) => {
        const { id, status: productStatus } = res.value;
        onNextStep({ id, productStatus });
      });
    });
  };

  render() {
    return (
      <div className={classnames({ [css.details]: true, wrapped: true })}>
        <PublisherSubmitSourceForm onSubmit={this.onSubmit} />
        <div className={css.tips}>
          {TIPS.map(tip => (
            <div key={tip.title} className={css.tip}>
              <h5>{tip.title}</h5>
              <p>{tip.description}</p>
            </div>
          ))}
        </div>
      </div>
    );
  }
}
