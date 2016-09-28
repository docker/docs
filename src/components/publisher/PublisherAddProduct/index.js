import React, { Component, PropTypes } from 'react';
import { ChevronIcon } from 'common';
import PublisherProductSource from './PublisherProductSource';
import PublisherProductInformation from './PublisherProductInformation';
import PublisherProductTiers from './PublisherProductTiers';
import PublisherProductSubmitted from './PublisherProductSubmitted';
import { LARGE } from 'lib/constants/sizes';
import {
  getProductStep,
  getCurrentStep,
} from 'lib/utils/publisherSteps';
import { connect } from 'react-redux';
import get from 'lodash/get';
import classnames from 'classnames';

import css from './styles.css';

const { string, object, number, shape, func } = PropTypes;

const mapStateToProps = ({ form, publish }, { location }) => {
  const currentProductDetails = publish.currentProductDetails.results;
  const productStep = getProductStep(currentProductDetails.status);
  const currentStep =
    getCurrentStep(parseInt(location.query.step, 10), productStep);
  const display_name_value = get(form, ['publisherSubmitSourceForm',
    'display_name', 'value']) || '[New Product]';
  const display_name = currentProductDetails.id && currentProductDetails.name ?
    currentProductDetails.name : display_name_value;
  return {
    display_name,
    currentStep,
    productStep,
  };
};

@connect(mapStateToProps)
export default class PublisherAddProduct extends Component {
  static propTypes = {
    currentStep: number.isRequired, // integer between 0 - 4
    display_name: string.isRequired,
    location: object.isRequired,
    productStep: number.isRequired, // integer between 0 - 4
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  onClickBack = () => {
    const step = this.props.currentStep;
    if (step > 0) {
      const { query } = this.props.location;
      const product_id = query.product_id;
      this.context.router.replace({
        pathname: location.pathname,
        query: {
          step: step - 1,
          product_id,
        },
      });
    }
  }

  onStepClick = (step, productStep) => () => {
    if (step <= productStep) {
      const { query } = this.props.location;
      const product_id = query.product_id;
      this.context.router.replace({
        pathname: location.pathname,
        query: {
          step,
          product_id,
        },
      });
    }
  }

  onNextStep = ({ productStatus, id }) => {
    // NOTE: Sets the query param & path to next steps page
    const step = getProductStep(productStatus);
    this.context.router.replace({
      pathname: location.pathname,
      query: {
        step,
        product_id: id,
      },
    });
  }

  render() {
    const {
      currentStep,
      productStep,
    } = this.props;

    let sourceStatus = 'Submit source repositories';
    let infoStatus = 'Not submitted';
    let pricingStatus = 'Not submitted';
    let publishStatus = 'Not approved';
    switch (productStep) {
      case 1:
      // publisher has submitted source repo's but has not submitted product
      // details or information yet.
        sourceStatus = 'Source submitted';
        infoStatus = 'Add product details';
        break;
      case 2:
      // publisher has submitted source repo's and product details, but has not
      // submitted product rate plans yet.
        sourceStatus = 'Source submitted';
        infoStatus = 'Product information submitted';
        pricingStatus = 'Add product tiers';
        break;
      case 3:
      // publisher has submitted source repo's, details, and rate plans but has
      // not submitted product for approval
        sourceStatus = 'Source submitted';
        infoStatus = 'Product information submitted';
        pricingStatus = 'Pricing product tiers submitted';
        publishStatus = 'Product waiting approval';
        break;
      case 4:
      // publisher has submitted product but has not been approved yet.
        sourceStatus = 'Source submitted';
        infoStatus = 'Product information submitted';
        pricingStatus = 'Pricing product tiers submitted';
        publishStatus = 'Moderator feedback received';
        break;
      case 5:
      // publisher product has been approved and is waiting to be published
        sourceStatus = 'Source submitted';
        infoStatus = 'Product information submitted';
        pricingStatus = 'Pricing product tiers submitted';
        publishStatus = 'Product published to the Store!';
        break;
      default:
        sourceStatus = 'Submit source repositories';
        infoStatus = 'Not submitted';
        pricingStatus = 'Not submitted';
        publishStatus = 'Not approved';
    }
    let childComponent = (
      <div></div>
    );
    switch (currentStep) {
      case 1:
        childComponent = (
          <PublisherProductInformation
            onClickBack={this.onClickBack}
            onNextStep={this.onNextStep}
          />
        );
        break;
      case 2:
        childComponent = (
          <PublisherProductTiers
            onClickBack={this.onClickBack}
            onNextStep={this.onNextStep}
          />
        );
        break;
      case 3:
        childComponent = (
          <PublisherProductSubmitted
            onClickBack={this.onClickBack}
          />
        );
        break;
      // TODO: Add back in cases once we have these flows up
      // case 4:
      // case 5:
      default:
        childComponent = (
          <PublisherProductSource
            onNextStep={this.onNextStep}
          />
        );
        break;
    }

    return (
      <div>
        <div className={css.statusbar}>
          <div className="wrapped">
            <div className={css.statusbarcontents}>
              <h3>{this.props.display_name}</h3>
            </div>
          </div>
        </div>
        <div className="wrapped">
          <ul className={css.steps}>
            <li
              className={classnames({
                [css.step]: true,
                [css.chevStep]: true,
                [css.active]: currentStep === 0,
              })}
              onClick={this.onStepClick(0, productStep)}
            >
              <div>
                <h4>PRODUCT SOURCE</h4>
                <h5>{sourceStatus}</h5>
              </div>
              <ChevronIcon className={css.chevron} size={LARGE} />
            </li>
            <li
              className={classnames({
                [css.step]: true,
                [css.chevStep]: true,
                [css.disabled]: productStep < 1,
                [css.active]: currentStep === 1,
              })}
              onClick={this.onStepClick(1, productStep)}
            >
              <div>
                <h4>PRODUCT INFORMATION</h4>
                <h5>{infoStatus}</h5>
              </div>
              <ChevronIcon className={css.chevron} size={LARGE} />
            </li>
            <li
              className={classnames({
                [css.step]: true,
                [css.chevStep]: true,
                [css.disabled]: productStep < 2,
                [css.active]: currentStep === 2,
              })}
              onClick={this.onStepClick(2, productStep)}
            >
              <div>
                <h4>PRODUCT PRICING</h4>
                <h5>{pricingStatus}</h5>
              </div>
              <ChevronIcon className={css.chevron} size={LARGE} />
            </li>
            <li
              className={classnames({
                [css.step]: true,
                [css.disabled]: productStep < 3,
                [css.active]: currentStep > 2,
              })}
              onClick={this.onStepClick(3, productStep)}
            >
              <div>
                <h4>PUBLISH TO STORE</h4>
                <h5>{publishStatus}</h5>
              </div>
            </li>
          </ul>
        </div>
        {childComponent}
      </div>
    );
  }
}
