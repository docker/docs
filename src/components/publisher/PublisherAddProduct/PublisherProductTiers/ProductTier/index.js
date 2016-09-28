import React, { Component, PropTypes } from 'react';
import {
  Card,
  Checkbox,
  CloseIcon,
  Input,
} from 'components/common';
import { SMALL } from 'lib/constants/sizes';

import css from './styles.css';
const { array, shape, number, object, func } = PropTypes;

export default class ProductTier extends Component {

  static propTypes = {
    index: number.isRequired,
    productTier: shape({
      tierName: object.isRequired,
      price: object.isRequired,
      reposField: object.isRequired,
    }).isRequired,
    defaultTier: object.isRequired,
    repositories: array.isRequired,
    removeTier: func.isRequired,
    tierLength: number.isRequired,
  }

  onCheck = (field) => () => {
    field.onChange(!field.value);
  }

  onSelectDefault = () => {
    const { defaultTier } = this.props;
    defaultTier.onChange(this.props.index);
  }

  onSelectRepo = (name) => (e) => {
    const {
      reposField,
    } = this.props.productTier;
    let repoValues = reposField.value || [];
    if (e.target.checked) {
      // repoValues.push(name); <-- allows multiple repos per product tier
      // NOTE: This effectively makes repo's radio buttons - allowing only 1
      // repo per product tier.
      repoValues = [name]; // <-- replaces array with only a single repo
    } else {
      repoValues =
        repoValues.filter(repo => repo !== name);
    }
    reposField.onChange(repoValues);
  }

  renderRepoSelect = (repo, idx) => {
    const {
      reposField,
    } = this.props.productTier;
    const tag = repo.tag ? (`/${repo.tag}`) : '';
    const fullRepository = `${repo.namespace}/${repo.reponame}${tag}`;
    /*
    NOTE:
    reposField form is an array where on check we push the fullname into the arr
    */
    const checked = reposField.value ?
      reposField.value.indexOf(fullRepository) >= 0 : false;
    return (
      <Checkbox
        className={css.repoSelect}
        key={idx}
        label={fullRepository}
        onCheck={this.onSelectRepo(fullRepository)}
        checked={checked}
      />
    );
  }

  renderBillingRow() {
    const {
      productTier: {
        tierName,
        price,
        hasFreeTrial,
      },
      defaultTier,
      index,
    } = this.props;
    /*
    NOTE: the Checkbox component doesn't work well with redux forms
    onBlur and value are required AFTER deconstructing the field object
    otherwise the deconstructed values will overwrite the passed values
    */
    let addTrial;
    if (price.value > 0) {
      addTrial = (
        <Checkbox
          label="Add a 30 day free trial that will transition to a paid plan"
          {...hasFreeTrial}
          checked={!!hasFreeTrial.value}
          onBlur={() => {}}
          onCheck={this.onCheck(hasFreeTrial)}
          value={'freeTrial'}
        />
      );
    }
    const removeTier = index > 0 ? (
      <div
        className={css.removeField}
        onClick={this.props.removeTier}
      >
        <CloseIcon size={SMALL} />
      </div>
    ) : (<div></div>);
    return (
      <div className={css.billingRow}>
        <div>
          <div className={css.sectionTitle}>Tier Name</div>
          <Input
            className={css.input}
            id="product-names-input"
            style={{ marginBottom: '14px' }}
            errorText={tierName.touched && tierName.error || ''}
            { ...tierName }
          />
        </div>
        <div>
          <div className={css.sectionTitle}>Price/Month</div>
          <div className={css.priceDelim}>
            <span>$</span>
            <div>
              <Input
                className={css.price}
                id="product-names-input"
                style={{ marginBottom: '14px' }}
                errorText={price.touched && price.error || ''}
                { ...price }
              />
            </div>
          </div>
        </div>
        <div className={css.options}>
          <Checkbox
            {...defaultTier}
            label="Default"
            onCheck={this.onSelectDefault}
            onBlur={() => {}}
            checked={defaultTier.value === index}
            value={`index${index}`}
          />
          {addTrial}
        </div>
        {removeTier}
      </div>
    );
  }

  render() {
    const {
      productTier: {
        description,
        vendorAgreement,
        instructions,
        reposField,
      },
      repositories,
      tierLength,
    } = this.props;

    return (
      <Card
        className={css.card}
      >
        {this.renderBillingRow()}

        <div className={css.fields}>
          <div className={css.sectionTitle}>Select Repository</div>
          {repositories.map(this.renderRepoSelect)}
          <div className={css.error}>
            {reposField.touched && reposField.error || ''}
          </div>
        </div>

        <div className={css.fields}>
          <div className={css.sectionTitle}>Tier Description</div>
          <div className={css.subText}>
            What features does this specific tier include?
          </div>
          <textarea
            className={css.textArea}
            placeholder={`Max ${tierLength} characters`}
            { ...description}
          />
          <div className={css.error}>
            {description.touched && description.error || ''}
          </div>
        </div>

        <div className={css.fields}>
          <div className={css.sectionTitle}>Vendor Agreement</div>
          <div className={css.subText}>
            Enter url of the page where your customers can read the license
            agreement for this plan.
          </div>
          <Input
            className={css.input}
            id="product-names-input"
            errorText={vendorAgreement.touched && vendorAgreement.error || ''}
            { ...vendorAgreement }
          />
        </div>

        <div className={css.fields}>
          <div className={css.sectionTitle}>Installation Instructions</div>
          <div className={css.subText}>
            Add instructions on how to install and configure the product. if
            possible, include a link to installation documentation.
          </div>
          <textarea
            className={css.textArea}
            placeholder="Step by step instructions"
            { ...instructions}
          />
          <div className={css.error}>
            {instructions.touched && instructions.error || ''}
          </div>
        </div>
      </Card>
    );
  }
}
