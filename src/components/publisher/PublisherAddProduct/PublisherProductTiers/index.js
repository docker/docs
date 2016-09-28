import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import classnames from 'classnames';
import find from 'lodash/find';
import reject from 'lodash/reject';
import partition from 'lodash/partition';
import merge from 'lodash/merge';
import cloneDeep from 'lodash/cloneDeep';

import {
  publishCreateProductTiers,
  publishDeleteProductTiers,
  publishFetchProductDetails,
  publishFetchProductTiers,
  publishUpdateProductDetails,
  publishUpdateProductTiers,
} from 'actions/publish';
import ProductTier from './ProductTier';
import {
  BackButtonArea,
  Button,
  ChevronIcon,
} from 'components/common';
import { DOWNLOAD_ATTRIBUTES } from 'lib/constants/eusa';

import { STEP2, STEP3 } from 'lib/constants/publisherSteps';
import css from './styles.css';

const { array, func, object, bool, string } = PropTypes;
const MAX_DESCRIPTION_LENGTH = 140;
const MAX_NAME_LENGTH = 100;
const FLAT_CHARGE_LABEL = 'FlatCharge';

const validate = (values) => {
  const errors = {};
  if (!values.productTiers[values.defaultTier]) {
    // eslint-disable-next-line no-underscore-dangle
    errors._error = 'You must include at least 1 product tier.';
  }
  const productTiers = values.productTiers;
  if (productTiers.length < 1) {
    // eslint-disable-next-line no-underscore-dangle
    errors._error = 'You must include at least 1 product tier.';
  } else {
    errors.productTiers = [];
    productTiers.forEach((tier, idx) => {
      errors.productTiers[idx] = {};
      if (
        !tier.price ||
        !tier.price.trim() ||
        isNaN(Number(tier.price)) ||
        Number(tier.price) < 0
      ) {
        // Allow for free tiers
        errors.productTiers[idx].price = 'Invalid price.';
      }
      if (!tier.description) {
        errors.productTiers[idx].description = 'Required.';
      } else if (tier.description.length > MAX_DESCRIPTION_LENGTH) {
        errors.productTiers[idx].description = 'Description is too long.';
      }
      if (!tier.tierName) {
        errors.productTiers[idx].tierName = 'Required.';
      } else if (tier.tierName.length > MAX_NAME_LENGTH) {
        errors.productTiers[idx].tierName =
          `Tier name can only be ${MAX_NAME_LENGTH} characters long.`;
      }
      if (!tier.reposField || tier.reposField.length < 1) {
        errors.productTiers[idx].reposField = 'Required.';
      }
      if (!tier.instructions) {
        errors.productTiers[idx].instructions = 'Required.';
      }
    });
  }
  return errors;
};

// NOTE - Helper to format BE response to fit in redux form
const formatTierForm = (tier) => {
  const {
    id,
    name: tierName,
    repositories,
    description,
    instructions,
    is_default: isDefault,
    eusa: vendorAgreement,
    price,
  } = tier;

  const reposField = repositories.map((repo) => {
    const publisherRepo = repo.publishers_repo_name;
    const tag = publisherRepo.tag ? `/${publisherRepo.tag}` : '';
    return `${publisherRepo.namespace}/${publisherRepo.reponame}${tag}`;
  });
  const {
    trial,
    pricing_components,
  } = price;
  const flatCharge = find(pricing_components, (pc) => {
    return pc.label === FLAT_CHARGE_LABEL;
  });
  return {
    description,
    hasFreeTrial: !!trial,
    id,
    instructions,
    isDefault,
    price: flatCharge.tiers[0].price,
    reposField,
    tierName,
    vendorAgreement,
  };
};

const mapStateToProps = ({ publish }) => {
  const {
    currentProductDetails: productDetails,
    currentProductTiers,
  } = publish;
  const currentProductDetails = productDetails.results;
  const productTiers =
    currentProductTiers.results && currentProductTiers.results.length ?
    currentProductTiers.results.map(formatTierForm) :
    /*
    NOTE: initialize dynamic form with at least 1 tier
    */
    [{
      description: '',
      hasFreeTrial: false,
      instructions: '',
      isDefault: true,
      price: '',
      reposField: [],
      tierName: '',
      vendorAgreement: '',
    }];
  const initialValues = {
    productTiers,
    defaultTier: 0,
  };
  return {
    currentProductDetails,
    currentProductTiers: currentProductTiers.results,
    initialValues,
  };
};

const mapDispatch = {
  publishCreateProductTiers,
  publishDeleteProductTiers,
  publishFetchProductDetails,
  publishFetchProductTiers,
  publishUpdateProductDetails,
  publishUpdateProductTiers,
};

class PublisherProductTiers extends Component {
  static propTypes = {
    currentProductDetails: object.isRequired,
    currentProductTiers: array.isRequired,
    onClickBack: func.isRequired,

    // redux form
    submitting: bool.isRequired,
    fields: object.isRequired,
    handleSubmit: func.isRequired,
    error: string,

    // actions
    onNextStep: func.isRequired,
    publishCreateProductTiers: func.isRequired,
    publishDeleteProductTiers: func.isRequired,
    publishFetchProductDetails: func.isRequired,
    publishFetchProductTiers: func.isRequired,
    publishUpdateProductDetails: func.isRequired,
    publishUpdateProductTiers: func.isRequired,
  }

  onSubmit = (submitStep) => (values) => {
    // NOTE: if a user is saving - STATUS will be STEP2. Submitting is STEP3
    const {
      currentProductDetails,
      currentProductTiers,
      onNextStep,
      publishCreateProductTiers: createProductTiers,
      publishDeleteProductTiers: deleteProductTier,
      publishFetchProductDetails: fetchProductDetails,
      publishFetchProductTiers: fetchProductTiers,
      publishUpdateProductDetails: updateProductDetails,
      publishUpdateProductTiers: updateProductTiers,
    } = this.props;
    const { id: product_id } = currentProductDetails;
    // separate new tiers from the already submitted tiers
    // If a tier has an id then it has already been previously created
    const [updateTiers, newTiers] =
      partition(values.productTiers, (tier) => { return tier.id; });
    const promises = [];
    if (updateTiers.length) {
      const tiersObject = {};
      updateTiers.forEach((tier) => {
        tiersObject[tier.id] = this.formatTiersSubmit(tier);
      });
      promises.push(updateProductTiers({ product_id, tiersObject }));
    }
    if (newTiers.length) {
      const tiersList = newTiers.map(this.formatTiersSubmit);
      promises.push(createProductTiers({ product_id, tiersList }));
    }
    // find tiers that are no longer in the form and delete them
    /*
    NOTE: if a user removes a tier, then adds it back - this will be a NEW tier
    and the old one will be deleted.
    */
    const deleteTiers = reject(currentProductTiers, (existingTier) => {
      return find(updateTiers, (updateTier) => {
        return updateTier.id === existingTier.id;
      });
    });
    if (deleteTiers.length) {
      deleteTiers.forEach((tier) => {
        promises.push(deleteProductTier({ product_id, tier_id: tier.id }));
      });
    }
    // TODO: 08/17 - nathan Add server validations if calls fail
    return Promise.all(promises).then(() => {
      return fetchProductTiers({ product_id });
    })
    .then(() => {
      // Update details PURELY to bump the product status.
      const details = cloneDeep(currentProductDetails);
      details.status = submitStep;
      return updateProductDetails({ product_id, details });
    })
    .then(() => {
      return fetchProductDetails({ product_id });
    })
    .then((res) => {
      const { id, status: productStatus } = res.value;
      // NOTE: onNextStep sends you to the next step in the publisher flow
      onNextStep({ id, productStatus });
    });
  };

  onRemoveTier = (index) => () => {
    const { productTiers } = this.props.fields;
    productTiers.removeField(index);
  };

  formatTiersSubmit = (formTier) => {
    /*
    Formats collected form data to backend requirements
    */
    const {
      currentProductTiers,
    } = this.props;
    const {
      description,
      hasFreeTrial,
      id,
      instructions,
      isDefault: is_default,
      price,
      reposField,
      tierName: name,
      vendorAgreement: eusa,
    } = formTier;
    /*
    NOTE:
    Each tier can have multiple different pricing components with different
    component tiers. Pricing component AND component tier 'name' fields are
    being hardcoded in the backend - so FOR NOW we only need to worry about the
    'trial' field, and the nested 'price' field of the tier.
    */
    const priceComponent = {
      trial: hasFreeTrial ? 1 : 0,
      pricing_components: [{
        tiers: [{
          price,
        }],
      }],
    };
    const newData = {
      name,
      description,
      instructions,
      is_default,
      eusa,
      download_attribute: DOWNLOAD_ATTRIBUTES.POST_EUSA,
      price: priceComponent,
    };
    /*
    IF tier already exists in fetched data - then we know to update
    */
    const currentTier = find(currentProductTiers, (tier) => {
      return tier.id === id;
    });
    const currentRepos = currentTier ? currentTier.repositories : [];
    const repositories = [];
    reposField.forEach((repo) => {
      // reposField from form will always be a string as 'namespace/repo/tag'
      // from onSelectRepo of ProductTier
      const [namespace, reponame, tag] = repo.split('/');
      // If repo already exists - push the repo object into the array since this
      // includes more pre-filled repo data that will be overwritten otherwise
      const repoObj =
        find(currentRepos, ({ publishers_repo_name: prn }) => {
          return prn.namespace === namespace &&
                 prn.reponame === reponame &&
                 prn.tag === tag;
        }) ||
      // Else if repo does not already exist, create a new repo object to push
        {
          publishers_repo_name: {
            namespace,
            reponame,
            tag,
          },
        };
      repositories.push(repoObj);
    });

    return merge(currentTier, newData, { repositories });
  }

  renderProductTiers = (tier, idx) => {
    const {
      currentProductDetails: {
        repositories,
      },
      fields: {
        defaultTier,
      },
    } = this.props;
    return (
      <ProductTier
        key={idx}
        index={idx}
        defaultTier={defaultTier}
        productTier={tier}
        removeTier={this.onRemoveTier(idx)}
        repositories={repositories || []}
        tierLength={MAX_DESCRIPTION_LENGTH}
        onSelectDefault={this.onSelectDefault}
      />
    );
  }

  render() {
    const {
      fields: {
        productTiers,
      },
      onClickBack,
      submitting,
      handleSubmit,
      error,
    } = this.props;

    // currenlty limiting users to only 2 product tiers
    // if you have 2 + tiers - don't show the add another product tier link
    const addTier = productTiers.length < 2 ? (
      <div className={css.addRatePlan}>
        <a
          href="#"
          onClick={(e) => {
            e.preventDefault();
            productTiers.addField();
          }}
        >Add another product tier</a>
      </div>
    ) : null;
    return (
      <div className={classnames({ [css.details]: true, wrapped: true })}>
        <form>
          {productTiers.map(this.renderProductTiers)}
          {addTier}
          <div className={css.buttonWrapper}>
            <BackButtonArea
              onClick={onClickBack}
              className={css.backButton}
              text="Previous Step"
            />
            <div className={css.submitWrapper}>
              <div className={css.globalError}>
                {error}
              </div>
              <Button
                className={css.submit}
                disabled={submitting}
                onClick={handleSubmit(this.onSubmit(STEP2))}
              >
                Save
              </Button>
              <Button
                className={css.submit}
                disabled={submitting}
                onClick={handleSubmit(this.onSubmit(STEP3))}
              >
                Submit
                <ChevronIcon className={css.chevron} />
              </Button>
            </div>
          </div>
        </form>
        <div className={css.tips}>
          <div className={css.tip}>
            <div className={css.tipTitle}>Set up product tier</div>
            <p className={css.tipBody}>
              Specify how much you will charge for each version of the
              product.
              This price includes all taxes.
            </p>
          </div>
          <div className={css.tip}>
            <div className={css.tipTitle}>Failed payment behavior</div>
            <p className={css.tipBody}>
              The product version that you set as "Default" is the one
              automatically selected for customers when they visit the
              product page.
              Use this to highlight one of the product versions.
            </p>
          </div>
          <div className={css.tip}>
            <div className={css.tipTitle}>Set up free trials</div>
            <p className={css.tipBody}>
              Select the free trial box to offer your customers a 30 day trial.
              We&#39;ll skip the charge for this plan for the first 30 days.
              Once 30 days have elapsed, your customer will be charged as usual.
            </p>
          </div>
          <div className={css.tip}>
            <div className={css.tipTitle}>Vendor license agreement</div>
            <p className={css.tipBody}>
              Add a URL for your license agreement. This will be presented as a
              link that customers can click to read more, and a checkbox that
              indicates their agreement.
            </p>
          </div>
          <div className={css.tip}>
            <div className={css.tipTitle}>Installation per version</div>
            <p className={css.tipBody}>
              Different product versions may require different installation
              instructions.
              Provide those here.
            </p>
          </div>
        </div>
      </div>
    );
  }
}

export default reduxForm({
  form: 'publisherSubmitInformationForm',
  fields: [
    'defaultTier',
    'productTiers[].description',
    'productTiers[].hasFreeTrial',
    'productTiers[].id', // field only filled if tier pre-exists in catalogue db
    'productTiers[].instructions',
    'productTiers[].isDefault',
    'productTiers[].price',
    'productTiers[].reposField',
    'productTiers[].tierName',
    'productTiers[].vendorAgreement',
  ],
  validate,
},
mapStateToProps,
mapDispatch,
)(PublisherProductTiers);
