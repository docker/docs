import React, { Component, PropTypes } from 'react';
const request = require('superagent-promise')(require('superagent'), Promise);
import md5 from 'md5';
import { SMALL, LARGE } from 'lib/constants/sizes';
import { FALLBACK_ELEMENT } from 'lib/constants/fallbacks';
import {
  Card,
  Checkbox,
  CloseIcon,
  Input,
  UploadIcon,
} from 'common';

import css from './styles.css';

const { array, object } = PropTypes;

export default class ProductDetailsForm extends Component {
  static propTypes = {
    categories: object.isRequired,
    categoryList: object.isRequired,
    description: object.isRequired,
    productIcon: object.isRequired,
    publisherDetails: object,
    screenshots: array.isRequired,
    searchTags: object.isRequired,
    supportLink: object.isRequired,
    tagline: object.isRequired,
  }

  onChangeCategory = (name) => (e) => {
    const {
      categoryList,
      categories: categoryField,
     } = this.props;
    let categoryValues = categoryField.value || [];
    if (e.target.checked) {
      categoryValues.push({ name, label: categoryList[name] });
    } else {
      categoryValues =
        categoryValues.filter(category => category.name !== name);
    }
    categoryField.onChange(categoryValues);
  }

  /*
    TASKS: nathan 08/04/16
    1 - Collect image file
    2 - Generate md5 encoded string
    3 - Send md5 to BE and receive s3 url
    4 - Put file to received url
    5 - Receive s3 uploaded image link
    6 - set value of screenFormField to url of image
      - screenFormField.url.onChange(data.url);
  */

  // NOTE: nathan 08/04/16 Not currently being used - reimplement for s3 upload
  onChangeScreenshotFile = (screenFormField) => (e) => {
    const file = e.target.files[0];
    if (!file) {
      return;
    }
    const md5file = md5(file);
    const extension = file.name.split('.').pop();
    if (extension !== 'png') {
      console.error('Need png extension');
      return;
    }
    screenFormField.url.onChange(this.getImageLink(md5file));
  }

  // NOTE: nathan 08/04/16 Not currently being used - reimplement for s3 upload
  onChangeProductIcon = (productIcon) => (e) => {
    const file = e.target.files[0];
    if (!file) {
      return;
    }
    const md5file = md5(file);
    const extension = file.name.split('.').pop();
    if (extension !== 'png') {
      console.error('Need png extension');
      return;
    }
    productIcon.onChange(this.getImageLink(md5file));
  }

  // NOTE: nathan 08/04/16 Not currently being used - reimplement for s3 upload
  onProductIconClick = (e) => {
    const fileElem = document.getElementById('productUploadIcon');
    if (fileElem) {
      fileElem.click();
    }
    e.preventDefault();
  }

  // NOTE: Not currently being used - reimplement for s3 upload
  getImageLink = (md5file) => {
    return request
      .post('')
      .accept('application/json')
      .send(md5file)
      .end()
      .then((res) => {
        return res.body;
        // screenFormField.url.onChange(res.value);
      });
  //     debugger;
  //     const data = JSON.parse(res.text);
  //     return request
  //       .put(data.signedRequest)
  //       .set('Content-Type', file.type)
  //       .send(file)
  //       .end()
  //   })
  //   .catch(err => {
  //     // TODO (jmorgan): show error text in form if creating a product failed
  //     console.log(err);
  //   });
  }

  renderScreenShot = (screenFormField, index) => {
    const {
      screenshots,
    } = this.props;
    const {
      label,
      url,
    } = screenFormField;

    const preview = screenFormField.url.value ? (
      <img
        alt=""
        className={css.screenshotpreview}
        src={screenFormField.url.value}
      />
    ) : null;
    return (
      <div key={index}>
        <div className={css.screenInfo}>
          <Input
            id={`screenshot-input-${index}`}
            type="text"
            style={{ marginBottom: '14px', width: '' }}
            placeholder="Label"
            errorText={label.touched && label.error || ''}
            { ...label }
          />
          <Input
            id={`screenshot-input-${index}`}
            type="text"
            style={{ marginBottom: '14px', width: '' }}
            placeholder="Screenshot URL"
            errorText={url.touched && url.error || ''}
            { ...url }
          />
          <div
            className={css.removeScreen}
            onClick={() => screenshots.removeField(index)}
          >
            <CloseIcon size={SMALL} />
          </div>
        </div>
        <div className={css.screenshotwrapper}>
          {preview}
        </div>

        {/* TODO: nathan 08/04/16 reimplement to get file upload to s3
          <input
            id={index}
            type="file"
            onChange={this.onChangeScreenshotFile(screenFormField)}
          />
        */}
      </div>
    );
  }

  render() {
    const {
      // searchTags,
      categories,
      categoryList,
      description,
      productIcon,
      screenshots,
      supportLink,
      tagline,
    } = this.props;

    // TODO: nathan 8/1/16 implement Search Tags input
    // const searchTags = (
    //   <div className={css.fields}>
    //     <div className={css.sectionTitle}>Search Tags</div>
    //     <div className={css.subText}>
    //       Have a search term that your product should show up for? Add it as
    //       a search tag. (Separate by comma's)
    //     </div>
    //     <Input
    //       id="product-names-input"
    //       errorText={searchTags.touched && searchTags.error || ''}
    //       { ...searchTags }
    //     />
    //   </div>
    // );

    const iconPreview = productIcon.value ? (
      <img
        alt="Product Icon"
        className={css.iconPreview}
        src={productIcon.value}
      />
    ) : FALLBACK_ELEMENT;

    return (
      <Card className={css.card}>
        <div className={css.title}>
          Product Details
        </div>
        <div className={css.productDetailsWrapper}>
          <div className={css.productIcon}>
            {/* TODO: nathan 08/04/16 - Add actual upload of image to s3
            <input
              type="file"
              id="productUploadIcon"
              accept="image/*"
              style={{ display: 'none' }}
              onChange={this.onChangeProductIcon(productIcon)}
            />
            <div
              className={css.uploadIcon}
              onClick={this.onProductIconClick}
            >
              <UploadIcon size={LARGE} variant="dull" />
            </div>
            <div>
              <span>Product Icon</span><br />
              <span>min 512x512px</span>
            </div>
            */}
            <div className={css.previewWrapper}>
              {iconPreview}
            </div>
          </div>
          <div>
            <div className={css.fields}>
              <div className={css.sectionTitle}>Product Tagline</div>
              <div className={css.subText}>
                A short description that appears with the product logo in search
                results. Max 140 characters.
              </div>
              <Input
                id="product-names-input"
                className={css.input}
                errorText={tagline.touched && tagline.error || ''}
                { ...tagline }
              />
            </div>
            <div className={css.fields}>
              <div className={css.sectionTitle}>Product Icon URL</div>
              <div className={css.subText}>
                Icon size must be a square image at least 512 X 512 pixels.
              </div>
              <Input
                id="product-icon-input"
                className={css.input}
                errorText={productIcon.touched && productIcon.error || ''}
                { ...productIcon }
              />
            </div>
            <div className={css.fields}>
              <div className={css.sectionTitle}>Categories</div>
              <div className={css.subText}>
                Select one or more Store categories that apply to your product.
                These will help customers find your product and similar
                services!
              </div>
              {
                Object.keys(categoryList).map(category =>
                  <Checkbox
                    key={category}
                    label={categoryList[category]}
                    onCheck={this.onChangeCategory(category)}
                    checked={categories.value ?
                      // eslint-disable-next-line
                      categories.value.filter(c => c.name === category).length > 0 :
                      false
                    }
                  />
                )
              }
            </div>
            <div className={css.fields}>
              <div className={css.sectionTitle}>Product Description</div>
              <div className={css.subText}>
                A longer product description that appears on the product page.
              </div>
              <textarea className={css.description} { ...description} />
              <div className={css.error}>
                {description.touched && description.error || ''}
              </div>
            </div>
            <div className={css.fields}>
              <div className={css.sectionTitle}>Support Link</div>
              <div className={css.subText}>
                the URL of your product's troubleshooting or support pages.
              </div>
              <Input
                id="product-names-input"
                className={css.input}
                errorText={supportLink.touched && supportLink.error || ''}
                { ...supportLink }
              />
            </div>
            <div className={css.fields}>
              <div className={css.sectionTitle}>Product Screenshots</div>
              <div className={css.subText}>
                Include at least 1 screenshot at 1920 by 1200 px or larger.
              </div>
              {screenshots.map(this.renderScreenShot)}
              <div className={css.addScreen}>
                <div
                  className={css.uploadIcon}
                  onClick={() => screenshots.addField()}
                >
                  <UploadIcon size={LARGE} variant="dull" />
                </div>
                <div>Add Screenshot</div>
              </div>
              <div className={css.screenErr}>{screenshots.error}</div>
            </div>
          </div>
        </div>
      </Card>
    );
  }
}
