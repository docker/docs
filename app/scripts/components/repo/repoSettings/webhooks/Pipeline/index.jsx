'use strict';
import React, { Component, PropTypes } from 'react';
// Some issues getting InlineSVG loaded
const InlineSVG = require('svg-inline-react');

import connectToStores from 'fluxible-addons-react/connectToStores';
import deletePipeline from 'actions/deletePipeline';
import PipelineHistoryStore from 'stores/PipelineHistoryStore';
import JWTStore from 'stores/JWTStore';
import { VelocityTransitionGroup } from 'velocity-react';

import Options from './Options.jsx';
import History from './History.jsx';

import styles from './index.css';

class Pipeline extends Component {
  static propTypes = {
    isActive: PropTypes.bool.isRequired,
    name: PropTypes.string.isRequired,
    namespace: PropTypes.string.isRequired,
    results: PropTypes.object,
    pipeline: PropTypes.shape({
      created: PropTypes.string.isRequired,
      expectFinalCallback: PropTypes.bool.isRequired,
      lastUpdated: PropTypes.string.isRequired,
      slug: PropTypes.string.isRequired,
      webhooks: PropTypes.arrayOf(PropTypes.shape({
        created: PropTypes.string,
        hookUrl: PropTypes.string.isRequired,
        lastUpdated: PropTypes.string,
        name: PropTypes.string
      }))
    })
  }

  static contextTypes = {
    executeAction: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    this.state = this._getInitialState(props);
    this.toggleHistory = this._toggleHistory.bind(this);
    this.deletePipeline = this._deletePipeline.bind(this);
  }

  _getInitialState(props) {
    return {
      deleting: false,
      isActive: this.props.isActive
    };
  }

  _toggleHistory() {
    this.setState({isActive: !this.state.isActive});
  }

  _deletePipeline() {
    this.setState({deleting: true});
    const { JWT: jwt, namespace, name, pipeline: { slug } } = this.props;
    this.context.executeAction(deletePipeline, {
      jwt, namespace, name, slug
    });
  }

  possiblyRenderHistory() {
    if (!this.state.isActive) {
      return undefined;
    }

    const attempts = this.props.results[this.props.pipeline.slug];
    console.log(attempts); // eslint-disable-line
    return (
      <History
        key='pipelineHistory'
        JWT={this.props.JWT}
        name={this.props.name}
        namespace={this.props.namespace}
        slug={this.props.pipeline.slug}
        attempts={(attempts && attempts.results) || undefined}
      />
    );
  }

  render() {
    const pipelineClasses = [styles.pipeline];

    if (this.state.deleting) {
      pipelineClasses.push(styles.pipeline_deleting);
    }

    const latestAttempt = this.props.results.length && this.props.results[0];

    return (
      <div
        className={pipelineClasses.join(' ')}
      >
        <div className={styles.overview}>
          <div className={styles.icon}>
            <InlineSVG src={require('../webhook.svg')} />
          </div>
          <div className={styles.summaryDetails}>
            <div className={styles.name}>
              {this.props.pipeline.name}
            </div>
            <div className={styles.url}>
              {this.props.pipeline.webhooks[0].hookUrl}
            </div>
          </div>
          <Options
            isActive={this.state.isActive}
            toggleHistory={this.toggleHistory}
            delete={this.deletePipeline}
          />
        </div>
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
          { this.possiblyRenderHistory() }
        </VelocityTransitionGroup>
      </div>
    );
  }
}

export default connectToStores(
  Pipeline,
  [
   PipelineHistoryStore,
   JWTStore
  ],
  ({ getStore }, props) => {
   const { results } = getStore(PipelineHistoryStore).getState();
   const { jwt } = getStore(JWTStore).getState();
   return {
     results,
     JWT: jwt
   };
 });
