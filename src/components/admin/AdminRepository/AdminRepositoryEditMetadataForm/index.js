const request = require('superagent-promise')(require('superagent'), Promise);
import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import { DOWNLOAD_ATTRIBUTES_READABLE } from 'lib/constants/eusa';
import { Input, Button, LabelField, Checkbox } from 'components/common';
import uuid from 'uuid-v4';
import css from './styles.css';
import { jwt } from 'lib/utils/authHeaders';

class AdminRepositoryEditMetadataForm extends Component {
  static propTypes = {
    categories: PropTypes.object.isRequired,
    platforms: PropTypes.object.isRequired,
    fields: PropTypes.object.isRequired,
    submitting: PropTypes.bool.isRequired,
    handleSubmit: PropTypes.func.isRequired,
  }

  onChangeCategory = (name) => (e) => {
    let categories = this.props.fields.categories.value || [];
    if (e.target.checked) {
      categories.push({ name, label: this.props.categories[name] });
    } else {
      categories = categories.filter(c => c.name !== name);
    }
    this.props.fields.categories.onChange(categories);
  }

  onChangePlatform = (name) => (e) => {
    let platforms = this.props.fields.platforms.value || [];
    if (e.target.checked) {
      platforms.push({ name, label: this.props.platforms[name] });
    } else {
      platforms = platforms.filter(c => c.name !== name);
    }
    this.props.fields.platforms.onChange(platforms);
  }

  onChangePreviewFile = (e) => {
    const file = e.target.files[0];
    if (!file) {
      return;
    }

    const extension = file.name.split('.').pop();
    if (extension !== 'png') {
      console.err('Need png extension');
      return;
    }

    request
      .get('/api/ui/admin/signed-s3-put')
      .set(jwt())
      .accept('application/json')
      .query({
        'file-name':
          `${this.props.fields.id.value}-${uuid()}-logo_large.${extension}`,
        'file-type': file.type,
      })
      .end()
      .then((res) => {
        const data = JSON.parse(res.text);
        return request
          .put(data.signedRequest)
          .set('Content-Type', file.type)
          .send(file)
          .end().then(() => {
            this.props.fields.logo_url.large.onChange(data.url);
          });
      }).catch(err => {
        // TODO (jmorgan): show error text in form if creating a product failed
        console.log(err);
      });
  }

  onChangeScreenshotFile = (screenshot) => (e) => {
    const file = e.target.files[0];
    if (!file) {
      return;
    }

    const extension = file.name.split('.').pop();
    if (extension !== 'png') {
      console.err('Need png extension');
      return;
    }

    request
      .get('/api/ui/admin/signed-s3-put')
      .accept('application/json')
      .query({
        'file-name':
          // eslint-disable-next-line
          `${this.props.fields.id.value}-${uuid()}-screenshot_large.${extension}`,
        'file-type': file.type,
      })
      .end()
      .then((res) => {
        const data = JSON.parse(res.text);
        return request
          .put(data.signedRequest)
          .set('Content-Type', file.type)
          .send(file)
          .end().then(() => {
            screenshot.url.onChange(data.url);
          });
      }).catch(err => {
        // TODO (jmorgan): show error text in form if creating a product failed
        console.log(err);
      });
  }

  render() {
    const {
      fields: {
        display_name,
        namespace,
        reponame,
        publisher,
        short_description,
        full_description,
        categories,
        platforms,
        logo_url,
        screenshots,
        links,
        eusa,
        download_attribute,
        instructions,
      },
      submitting,
      handleSubmit,
    } = this.props;

    return (
      <form onSubmit={handleSubmit} className={css.repositoryform}>
        <LabelField label="Display Name">
          <Input
            disabled
            id={'display_name_input'}
            placeholder="Display Name"
            { ...display_name}
          />
        </LabelField>
        <LabelField label="Store Repository">
          <div className={css.monospace}>
            store&nbsp;/&nbsp;
            <Input
              disabled
              id={'namespace_input'}
              placeholder="namespace"
              { ...namespace }
            />
            &nbsp;/&nbsp;
            <Input
              disabled
              id={'reponame_input'}
              placeholder="reponame"
              { ...reponame }
            />
          </div>
        </LabelField>
        <LabelField label="Publisher Docker ID">
          <Input
            disabled
            id={'publisher_id_input'}
            placeholder="Publisher Docker ID"
            { ...publisher.id }
          />
        </LabelField>
        <LabelField label="Publisher Name">
          <Input
            disabled
            id={'publisher_name_input'}
            placeholder="Publisher Name"
            { ...publisher.name }
          />
        </LabelField>
        <LabelField label="Entitlement">
          <select
            { ...download_attribute }
            value={download_attribute.value || ''}
          >
          <option></option>
          {Object.keys(DOWNLOAD_ATTRIBUTES_READABLE).map(k => (
            <option key={k} value={k}>
              {DOWNLOAD_ATTRIBUTES_READABLE[k]}
            </option>
          ))}
          </select>
        </LabelField>
        <LabelField label="Logo">
          <div className={css.logo}>
            <img
              alt=""
              className={css.preview}
              src={logo_url.large.value}
            />
            <input
              id={'logo-input'}
              type="file"
              onChange={this.onChangePreviewFile}
            />
            <label htmlFor="logo-input">{'Size must be 256x256 pixels'}</label>
          </div>
        </LabelField>
        <LabelField label="Screenshots">
          <div>
            {!screenshots.length && <div>No Screenshots</div>}
            {screenshots.map((screenshot, index) =>
              <div key={index}>
                <h3>Screenshot {index + 1}</h3>
                <div>
                  <div className={css.screenshotwrapper}>
                    <img
                      alt=""
                      className={css.screenshotpreview}
                      src={screenshot.url.value}
                    />
                  </div>
                  <Input
                    id={`screenshot-input-${index}`}
                    type="text"
                    placeholder="Label"
                    { ...screenshot.label }
                  />&nbsp;&nbsp;&nbsp;<input
                    id={index}
                    type="file"
                    onChange={this.onChangeScreenshotFile(screenshot)}
                  />
                  <label htmlFor="index">
                    Size must be 1920 by 1200 pixels or larger
                  </label>
                </div>
              </div>
            )}
            <div>
              <Button type="button" onClick={() => screenshots.addField()}>
                Add Screenshot
              </Button>
            </div>
          </div>
        </LabelField>
        <LabelField label="Short Description">
          <textarea { ...short_description} />
        </LabelField>
        <LabelField className={css.fulldescription} label="Full Description">
          <textarea { ...full_description} />
        </LabelField>
        <LabelField label="Instructions">
          <textarea placeholder="" { ...instructions } />
        </LabelField>
        <LabelField label="End-User License Agreement">
          <textarea placeholder="" { ...eusa } />
        </LabelField>
        <LabelField label="Links">
          <div>
            {!links.length && <div>No Links</div>}
            {links.map((link, index) =>
              <div key={index}>
                <div>
                  <Input
                    id={`links-label-input-${index}`}
                    type="text"
                    placeholder="Label"
                    { ...link.label }
                  />
                  &nbsp;&nbsp;&nbsp;
                  <Input
                    id={`links-url-input-${index}`}
                    type="text"
                    placeholder="https://..."
                    { ...link.url }
                  />
                </div>
              </div>
            )}
            <div>
              <Button type="button" onClick={() => links.addField()}>
                Add Link
              </Button>
            </div>
          </div>
        </LabelField>
        <LabelField label="Categories">
          <div>
            {Object.keys(this.props.categories).map(k =>
              <Checkbox
                key={k}
                label={this.props.categories[k]}
                onCheck={this.onChangeCategory(k)}
                checked={categories.value ?
                  categories.value.filter(c => c.name === k).length > 0 :
                  false
                }
              />
            )}
          </div>
        </LabelField>
        <LabelField label="Platforms">
          <div>
            {Object.keys(this.props.platforms).map(k =>
              <Checkbox
                key={k}
                label={this.props.platforms[k]}
                onCheck={this.onChangePlatform(k)}
                checked={platforms.value ?
                  !!platforms.value.filter(p => p.name === k).length :
                  false
                }
              />
            )}
          </div>
        </LabelField>
        <div className={css.savebutton}>
          <Button disabled={submitting} className={css.nomargin}>
            Save
          </Button>
        </div>
      </form>
    );
  }
}

export default reduxForm({
  form: 'adminRepositoryEditMetadataForm',
  fields: [
    'id',
    'display_name',
    'namespace',
    'reponame',
    'publisher.id',
    'publisher.name',
    'short_description',
    'full_description',
    'logo_url.large',
    'screenshots[].label',
    'screenshots[].url',
    'links[].label',
    'links[].url',
    'eusa',
    'download_attribute',
    'instructions',

    // TODO (jmorgan): change this to use x[].y redux-form syntax
    'categories',
    'platforms',
  ],
  validate: () => {
    const errors = {};
    return errors;
  },
})(AdminRepositoryEditMetadataForm);
