import React, { Component, PropTypes } from 'react';
import {
  publishFetchProductDetails,
  publishGetPublishers,
  publishUpdateProductDetails,
  publishUpdatePublisherInfo,
} from 'actions/publish';
import { reduxForm } from 'redux-form';
import classnames from 'classnames';
import find from 'lodash/find';

import {
  BackButtonArea,
  Button,
  ChevronIcon,
} from 'common';

import PublisherDetailsForm from './PublisherDetailsForm';
import ProductDetailsForm from './ProductDetailsForm';
import { STEP2 } from 'lib/constants/publisherSteps';
import {
  getProductStep,
} from 'lib/utils/publisherSteps';
import css from './styles.css';

const { func, object, bool, string } = PropTypes;

const validate = (values) => {
  const errors = {};

  if (!values.tagline) {
    errors.tagline = 'Required';
  } else if (values.tagline.length > 140) {
    errors.tagline = 'Tagline cannot be over 140 characters';
  }
  if (!values.description) {
    errors.description = 'Required';
  }
  if (!values.supportLink) {
    errors.supportLink = 'Required';
  }
  if (!values.company) {
    errors.company = 'Required';
  }
  if (!values.companyWebsite) {
    errors.companyWebsite = 'Required';
  }
  if (!values.productIcon) {
    errors.productIcon = 'Required';
  }

  const { screenshots } = values;
  if (screenshots.length < 1) {
    // eslint-disable-next-line no-underscore-dangle
    errors._error = 'You must include at least one screenshot!';
  } else {
    errors.screenshots = [];
    screenshots.forEach((screen, idx) => {
      errors.screenshots[idx] = {};
      if (!screen.label) {
        errors.screenshots[idx].label = 'Required';
      }
      if (!screen.url) {
        errors.screenshots[idx].url = 'Required';
      }
    });
  }

  return errors;
};

const mapStateToProps = ({ marketplace, publish }) => {
  const currentProductDetails = publish.currentProductDetails.results;
  const {
    company,
    links: publisherLinks,
  } = publish.publishers.results;
  const {
    categories,
    full_description,
    links: productLinks,
    logos,
    screenshots,
    short_description,
  } = currentProductDetails;
  const supportLink = find(productLinks, (link) => {
    return link.label === 'support';
  });
  const companyWebsite = find(publisherLinks, (link) => {
    return link.label === 'website';
  });
  const companyLogo = find(publisherLinks, (link) => {
    return link.label === 'logo';
  });
  const initialValues = {
    categories,
    company,
    companyWebsite: companyWebsite && companyWebsite.name || '',
    companyLogo: companyLogo && companyLogo.name || '',
    description: full_description,
    productIcon: logos && (logos.small || logos['small@2x']) || '',
    screenshots: screenshots && screenshots.map((screen) => {
      return { label: screen.label, url: screen.name };
    }) || [],
    supportLink: supportLink && supportLink.name,
    tagline: short_description,
  };

  return {
    initialValues,
    categoryList: marketplace.filters.categories,
    currentProductDetails,
  };
};

const mapDispatch = {
  publishFetchProductDetails,
  publishGetPublishers,
  publishUpdateProductDetails,
  publishUpdatePublisherInfo,
};

class PublisherProductInformation extends Component {
  static propTypes = {
    categoryList: object.isRequired,
    currentProductDetails: object.isRequired,
    onClickBack: func.isRequired,
    // functions
    onNextStep: func.isRequired,
    publishFetchProductDetails: func.isRequired,
    publishGetPublishers: func.isRequired,
    publishUpdateProductDetails: func.isRequired,
    publishUpdatePublisherInfo: func.isRequired,
    // redux form
    submitting: bool.isRequired,
    fields: object.isRequired,
    handleSubmit: func.isRequired,
    error: string,
  }

  onSubmit = values => {
    console.log(values);
    const {
      currentProductDetails,
      fields: {
        // searchTags,
        categories,
        company,
        companyLogo,
        companyWebsite,
        description,
        productIcon,
        screenshots: screenFormFields,
        supportLink,
        tagline,
      },
      onNextStep,
      publishFetchProductDetails: fetchProductDetails,
      publishGetPublishers: getPublisherInfo,
      publishUpdateProductDetails: updateProductDetails,
      publishUpdatePublisherInfo: updatePublisherInfo,
    } = this.props;

    const publisherInfo = {
      company: company.value,
      links: [
        {
          label: 'website',
          name: companyWebsite.value,
        },
        {
          label: 'logo',
          name: companyLogo.value,
        },
      ],
    };
    const screenshots = screenFormFields.filter((screen) => {
      return screen.url.value && screen.label.value;
    }).map((screen) => {
      return {
        name: screen.url.value,
        label: screen.label.value,
      };
    });
    const details = {
      name: currentProductDetails.name,
      status: STEP2,
      product_type: currentProductDetails.product_type,
      categories: categories.value,
      full_description: description.value,
      short_description: tagline.value,
      platforms: [],
      logos: {
        small: productIcon.value,
        'small@2x': productIcon.value,
      },
      links: [{
        name: supportLink.value,
        label: 'support',
      }],
      screenshots,
    };
    const product_id = currentProductDetails.id;
    return Promise.all([
      updateProductDetails({
        product_id,
        details,
      }).then(() => {
        return fetchProductDetails({ product_id });
      }).then((res) => {
        const { id, status: productStatus } = res.value;
        // NOTE: onNextStep sends you to the next step in the publisher flow
        onNextStep({ id, productStatus });
      }),
      updatePublisherInfo(publisherInfo).then(() => {
        return getPublisherInfo();
      }),
    ]);
  };

  render() {
    const {
      categoryList,
      currentProductDetails,
      fields: {
        categories,
        company,
        companyLogo,
        companyWebsite,
        description,
        productIcon,
        screenshots,
        searchTags,
        supportLink,
        tagline,
      },
      error,
      onClickBack,
      submitting,
      handleSubmit,
    } = this.props;

    const productStatus = currentProductDetails.status;
    const submit = getProductStep(productStatus) > 1 ?
    (
      <div className={css.submitWrapper}>
        <div className={css.globalError}>
          {error}
        </div>
        <Button
          className={css.submit}
          disabled={submitting}
        >
          Update
        </Button>
      </div>
    ) : (
      <div className={css.submitWrapper}>
        <div className={css.globalError}>
          {error}
        </div>
        <Button
          className={css.submit}
          disabled={submitting}
        >
          Save and continue
          <ChevronIcon className={css.chevron} />
        </Button>
      </div>
    );
    return (
      <form
        className={classnames({ [css.details]: true, wrapped: true })}
        onSubmit={handleSubmit(this.onSubmit)}
      >
        <div className={css.form}>
          <PublisherDetailsForm
            company={company}
            companyLogo={companyLogo}
            companyWebsite={companyWebsite}
            companyLogo={companyLogo}
          />
          <div className={css.tip}>
            <h5>Publisher Details</h5>
            <p>Let customers know who you are and what you do.</p>
          </div>
        </div>
        <div className={css.form}>
          <ProductDetailsForm
            categories={categories}
            categoryList={categoryList}
            description={description}
            productIcon={productIcon}
            screenshots={screenshots}
            searchTags={searchTags}
            supportLink={supportLink}
            tagline={tagline}
          />
          <div className={css.tip}>
            <h5>Create an appealing product profile</h5>
            <p>
              Briefly tell customers what your product is and what it does.
              What makes it stand out? This will be your first impression - make
              it count!
            </p>
          </div>
        </div>
        <div className={css.form}>
          <div className={css.buttonWrapper}>
            <BackButtonArea
              onClick={onClickBack}
              className={css.backButton}
              text="Previous Step"
            />
            {submit}
          </div>
        </div>
      </form>
    );
  }
}

export default reduxForm({
  form: 'publisherSubmitInformationForm',
  fields: [
    'categories',
    'company',
    'companyLogo',
    'companyWebsite',
    'description',
    'productIcon',
    'screenshots[].label',
    'screenshots[].url',
    'searchTags',
    'supportLink',
    'tagline',
  ],
  validate,
},
mapStateToProps,
mapDispatch,
)(PublisherProductInformation);
