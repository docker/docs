'use strict';

import React, {
  PropTypes,
  createClass,
  Component,
  cloneElement
} from 'react';

import find from 'lodash/collection/find';
import connectToStores from 'fluxible-addons-react/connectToStores';
import { VelocityTransitionGroup } from 'velocity-react';
// Some issues getting InlineSVG loaded
const InlineSVG = require('svg-inline-react');

import WebhooksSettingsStore from 'stores/WebhooksSettingsStore';

import AddWebhookForm from './AddWebhookForm';
import WebhooksTutorial from './WebhooksTutorial';
import Pipelines from './Pipelines';

const debug = require('debug')('webhooks');
import styles from './index.css';

class WebhooksSettings extends Component {
  static propTypes = {
    history: PropTypes.object.isRequired,
    pipelines: PropTypes.arrayOf(PropTypes.shape({
      created: PropTypes.string.isRequired,
      expectFinalCallback: PropTypes.bool.isRequired,
      lastUpdated: PropTypes.string.isRequired,
      name: PropTypes.string.isRequired,
      slug: PropTypes.string.isRequired,
      webhooks: PropTypes.arrayOf(PropTypes.shape({
        created: PropTypes.string,
        hookUrl: PropTypes.string.isRequired,
        lastUpdated: PropTypes.string,
        name: PropTypes.string
      }))
    }))
  }

  state = {
    showForm: false
  }

  toggleAddWebhookForm() {
    this.setState({showForm: !this.state.showForm});
  }

  hideAddWebhookForm() {
    this.setState({showForm: false});
  }

  showAddWebhookForm() {
    this.setState({showForm: true});
  }

  renderToggleFormButton() {
      const klasses = [styles.addWebHook];
      let title;

      if (this.state.showForm) {
        klasses.push(styles.addWebHook_rotated);
        title = 'Hide form';
      } else {
        title = 'Add new webhook';
      }

      return (
        <button
          className={klasses.join(' ')}
          title={title}
          onClick={this.toggleAddWebhookForm.bind(this)}
        >
          <InlineSVG src={require('./plus.svg')}/>
        </button>
      );
  }

  possiblyRenderTutorial() {
    if (this.state.showForm ||
    (this.props.pipelines && this.props.pipelines.length)) {
      return undefined;
    }

    return (
      <WebhooksTutorial
        key='tutorial'
        addWebhook={this.showAddWebhookForm.bind(this)}
      />
    );
  }

  possiblyRenderTutorialOrForm() {
    let form, tutorial;

    if (this.state.showForm) {
      const { namespace, name } = this.props;
      form = (
        <AddWebhookForm
          key='addForm'
          cancel={this.hideAddWebhookForm.bind(this)}
          namespace={namespace}
          name={name}
        />
      );
    } else if (!(this.props.pipelines && this.props.pipelines.length)) {
      tutorial = (
        <WebhooksTutorial
          key='tutorial'
          addWebhook={this.showAddWebhookForm.bind(this)}
        />
      );
    }

    // We only want to transition the form if there are pipelines
    // It looks weird otherwise to swap the form and tutorial with slides
    if (this.props.pipelines && this.props.pipelines.length) {
      return (
        <div>
          <VelocityTransitionGroup
            enter={{
              animation: 'slideDown',
              duration: '140ms'
            }}
            leave={{
              animation: 'slideUp',
              duration: '140ms'
            }}
          >
            { form }
          </VelocityTransitionGroup>
          { tutorial }
        </div>
      );
    } else {
      return (
        <div>
          { form }
          { tutorial }
        </div>
      );
    }
  }

  possiblyRenderPipelines() {
    if (!this.props.pipelines || !this.props.pipelines.length) {
      return undefined;
    }

    return (
      <Pipelines
        name={this.props.name}
        namespace={this.props.namespace}
        activeSlug={this.props.params.pipeline}
        pipelines={this.props.pipelines}
      />
    );
  }

  render() {
    return (
      <div className={styles.wrap}>
        <header className={styles.header}>Workflows</header>
        <div className={styles.leftColumn}>
          <header className={styles.segmentHeader}>Trigger Event</header>
          <segment className={styles.event}>
            <div className={styles.eventTitle}>
              <InlineSVG src={require('./imagePush.svg')}/>
              Image Pushed
            </div>
            <p className={styles.eventDescription}>
              When an image is pushed to this repo,
              your workflows will kick off based on
              your specified webhooks.&nbsp;
              <a
                href="https://docs.docker.com/docker-hub/webhooks/"
                title="Learn more about webhooks and events"
                target="_blank"
              >
                Learn More
              </a>
            </p>
          </segment>
        </div>
        <div className={styles.rightColumn}>
          <header className={styles.segmentHeader}>
            Web Hooks
            { this.renderToggleFormButton() }
          </header>
          { this.possiblyRenderTutorialOrForm() }
          { this.possiblyRenderPipelines() }
        </div>
      </div>
    );
  }
}

export default connectToStores(
  WebhooksSettings,
  [WebhooksSettingsStore],
  ({ getStore }, props) => getStore(WebhooksSettingsStore).getState()
);
