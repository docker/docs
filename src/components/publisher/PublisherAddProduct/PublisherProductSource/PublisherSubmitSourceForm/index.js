import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import {
  Button,
  Card,
  Checkbox,
  ChevronIcon,
  CloseIcon,
  Input,
  Modal,
  Select,
} from 'components/common';
import {
  repositoryFetchRepositoriesForNamespace,
  repositoryFetchImageTags,
} from 'actions/repository';
import { SMALL } from 'lib/constants/sizes';

import css from './styles.css';

const { array, func, bool, object, string } = PropTypes;
const VENDOR_AGREEMENT_URL = '#';

const mapStateToProps = ({ account, publish }) => {
  const currentProductDetails = publish.currentProductDetails.results;
  const {
    agreement,
    repositories,
    tags,
  } = publish.submit;
  const repoSources = currentProductDetails.repositories || [];
  const display_name = currentProductDetails.name;
  const publisher = publish.publishers.results;
  const accepted =
    publisher.vendor_agreement_status === 'accepted_online' ||
    publisher.vendor_agreement_status === 'accepted_offline';
  const initialValues = {
    accepted,
    display_name,
    repoSources,
  };
  return {
    initialValues,
    namespaces: account.ownedNamespaces || [],
    userHasAccepted: accepted,

    // Hub repo meta to populate source repo dropdowns
    agreement: agreement.results,
    agreementError: agreement.error,
    repositories,
    tags,
  };
};

const mapDispatch = {
  repositoryFetchRepositoriesForNamespace,
  repositoryFetchImageTags,
};

class PublisherSubmitSourceForm extends Component {
  static propTypes = {
    namespaces: array.isRequired,
    userHasAccepted: bool.isRequired,
    submitFailed: bool.isRequired,
    // Hub repo meta to populate source repo dropdowns
    agreement: string,
    agreementError: string,
    repositories: object.isRequired,
    tags: object.isRequired,
    // Redux Form methods/params
    submitting: bool.isRequired,
    fields: object.isRequired,
    untouch: func.isRequired,
    // Actions
    handleSubmit: func.isRequired,
    repositoryFetchImageTags: func.isRequired,
    repositoryFetchRepositoriesForNamespace: func.isRequired,
  }

  state = {
    showAgreementModal: false,
  }

  onShowAgreementModal = () => {
    const iframe = this.refs.agreementFrame;
    if (!iframe || !iframe.contentWindow) {
      return;
    }
    const frameDoc = iframe.contentWindow.document;
    frameDoc.write(this.props.agreement);
  }

  onChangeNamespace = (index) => (data) => {
    const { repoSources } = this.props.fields;
    const namespace = data.value;
    repoSources[index].namespace.onChange(namespace);
    repoSources[index].reponame.onChange('');
    repoSources[index].tag.onChange('');
    this.props.untouch(
      `repoSources[${index}].reponame`,
      `repoSources[${index}].tag`
    );
    this.props.repositoryFetchRepositoriesForNamespace({ namespace });
  }

  onChangeReponame = (index) => (data) => {
    const { repoSources } = this.props.fields;
    const reponame = data.value;
    const namespace = repoSources[index].namespace.value;
    repoSources[index].reponame.onChange(reponame);
    repoSources[index].tag.onChange('');
    this.props.untouch(`repoSources[${index}].tag`);
    this.props.repositoryFetchImageTags({
      namespace,
      reponame,
      page_size: 0,
    });
  }

  onCheck = () => {
    const field = this.props.fields.accepted;
    field.onChange(!field.value);
  }

  renderSourceRow = ({ namespace, reponame, tag }, index) => {
    const {
      fields,
      namespaces,
      repositories,
      tags,
    } = this.props;
    const { repoSources } = fields;
    const option = (value) => {
      return { label: value, value };
    };
    const namespaceOptions = namespaces.map(n => option(n));
    const selectedNamespace = namespace.value;
    const namespaceRepos = repositories[selectedNamespace];
    let reponameOptions = [];
    let repoFetching = namespaceRepos && namespaceRepos.isFetching;
    if (namespaceRepos && namespaceRepos.results) {
      reponameOptions = namespaceRepos.results.map(r => option(r));
    }
    const selectedRepo = reponame.value;
    let tagOptions = [];
    let tagFetching = false;
    if (tags[selectedNamespace] && tags[selectedNamespace][selectedRepo]) {
      const repoTags = tags[selectedNamespace][selectedRepo];
      if (repoTags.isFetching) {
        tagFetching = repoTags.isFetching;
      }
      if (repoTags.results) {
        tagOptions = repoTags.results.map(t => option(t));
      }
    }
    return (
      <div key={index} className={css.sourceRow}>
        <Select
          { ...namespace }
          className={css.select}
          placeholder="Select Namespace"
          options={namespaceOptions}
          onChange={this.onChangeNamespace(index)}
          onBlur={() => {}}
          ignoreCase
          backspaceRemoves={false}
          clearable={false}
          searchable={false}
          errorText={namespace.touched && namespace.error || ''}
        />
        <Select
          { ...reponame }
          className={css.select}
          placeholder="Select Repository"
          options={reponameOptions}
          onBlur={() => {}}
          onChange={this.onChangeReponame(index)}
          isLoading={repoFetching}
          ignoreCase
          backspaceRemoves={false}
          clearable={false}
          searchable={false}
          disabled={!namespace.value}
          errorText={reponame.touched && reponame.error || ''}
        />
        <Select
          { ...tag }
          className={css.select}
          placeholder="Select Tag"
          options={tagOptions}
          onBlur={() => {}}
          isLoading={tagFetching}
          ignoreCase
          backspaceRemoves={false}
          clearable={false}
          searchable={false}
          disabled={!reponame.value}
          errorText={tag.touched && tag.error || ''}
        />
        <div
          className={css.removeField}
          onClick={() => repoSources.removeField(index)}
        >
          <CloseIcon size={SMALL} />
        </div>
      </div>
    );
  }

  render() {
    const {
      fields: {
        display_name,
        repoSources,
        accepted,
      },
      userHasAccepted,
      submitting,
      handleSubmit,
      agreementError,
    } = this.props;
    return (
      <div>
        <form className={css.form} onSubmit={handleSubmit}>
          <Card className={css.card}>
            <div className={css.fields}>
              <div className={css.sectionTitle}>Product Name</div>
              <div className={css.subText}>
                Your product name as it will appear in the Docker Store
              </div>
              <Input
                className={css.displaynameinput}
                id="product-names-input"
                errorText={display_name.touched && display_name.error || ''}
                { ...display_name }
              />
            </div>
            <div className={css.fields}>
              <div className={css.sectionTitle}>
                Product source repository and tag
              </div>
              <div className={css.subText}>
                You can submit multiple private Docker Cloud or Hub repositories
                as sources for different tiers of your product. <br />
                Don't have any Docker images yet? <a href="https://success.docker.com/Store#Create_Great_Content">Click here.</a>
              </div>
              {
                !repoSources.length &&
                  <div className={css.noRep}>No Repositories</div>
              }
              {repoSources.map(this.renderSourceRow)}
              <div className={css.addRepo}>
                <a
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    repoSources.addField();
                  }}
                >Add another repository</a>
              </div>
            </div>
          </Card>
          <div className={css.submit}>
            <div className={css.fields}>
              <div className={css.agreement}>
                <Checkbox
                  className={css.agreementcheckbox}
                  checked={!!accepted.value}
                  onCheck={this.onCheck}
                  style={{ width: 'auto' }}
                  disabled={userHasAccepted}
                  onBlur={() => {}}
                />
                <div className={css.agreementlabel} onClick={this.onCheck}>
                  I have read and agree to the&nbsp;
                  <a
                    href={VENDOR_AGREEMENT_URL}
                    onClick={(e) => {
                      e.stopPropagation();
                      e.preventDefault();
                      this.setState({ showAgreementModal: true });
                      return false;
                    }}
                  >
                    vendor agreement
                  </a>
                </div>
                <Modal
                  onAfterOpen={this.onShowAgreementModal}
                  isOpen={this.state.showAgreementModal}
                  className={css.modal}
                  onRequestClose={() => {
                    this.setState({ showAgreementModal: false });
                  }}
                  style={{
                    content: {
                      overflow: 'auto',
                      left: '120px',
                      right: '120px',
                      top: '80px',
                      bottom: '80px',
                      position: 'absolute',
                      maxHeight: 'auto',
                      width: 'auto',
                    },
                  }}
                >
                  <div
                    className={css.closeModal}
                    onClick={() => {
                      this.setState({ showAgreementModal: false });
                    }}
                  >
                    <CloseIcon size={SMALL} />
                  </div>
                  {agreementError && 'Could not fetch vendor agreement'}
                  <iframe
                    className={css.agreementFrame}
                    ref="agreementFrame"
                  >
                  </iframe>
                </Modal>
              </div>
              <div className={css.error}>
                {accepted.touched && accepted.error}
              </div>
            </div>
            <div>
              <Button
                className={css.button}
                disabled={submitting}
              >
                Save and continue <ChevronIcon className={css.chevron} />
              </Button>
            </div>
          </div>
        </form>
      </div>
    );
  }
}

export default reduxForm({
  form: 'publisherSubmitSourceForm',
  fields: [
    'display_name',
    'repoSources[].namespace',
    'repoSources[].reponame',
    'repoSources[].tag',
    'accepted',
  ],
  validate: values => {
    const errors = {};

    if (!values.display_name) {
      errors.display_name = 'Required';
    }
    // validations for deep form
    const repoSources = values.repoSources;
    if (repoSources.length < 1) {
      errors.display_name = 'You must include at least one repository!';
    } else {
      errors.repoSources = [];
      repoSources.forEach((repo, idx) => {
        errors.repoSources[idx] = {};
        if (!repo.namespace) {
          errors.repoSources[idx].namespace = 'Required';
        }
        if (!repo.reponame) {
          errors.repoSources[idx].reponame = 'Required';
        }
        if (!repo.tag) {
          errors.repoSources[idx].reponame = 'Required';
        }
      });
    }

    if (!values.accepted) {
      errors.accepted =
        'You must accept the vendor agreement to continue';
    }

    return errors;
  },
},
mapStateToProps,
mapDispatch,
)(PublisherSubmitSourceForm);
