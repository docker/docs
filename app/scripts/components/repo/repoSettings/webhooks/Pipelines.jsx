'use strict';
import React, { Component, PropTypes } from 'react';
import ReactDOM from 'react-dom';
import { VelocityTransitionGroup } from 'velocity-react';

const debug = require('debug')('Webhooks');

import Pipeline from './Pipeline/index.jsx';

export default class Pipelines extends Component {
  static propTypes = {
    activeSlug: PropTypes.string,
    name: PropTypes.string.isRequired,
    namespace: PropTypes.string.isRequired,
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

  renderPipeline(pipeline) {
    return (
      <Pipeline
        key={pipeline.slug}
        ref={pipeline.slug}
        isActive={this.props.activeSlug === pipeline.slug}
        name={this.props.name}
        namespace={this.props.namespace}
        pipeline={pipeline}
      />
    );
  }

  render() {
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
          { this.props.pipelines.map(this.renderPipeline.bind(this)) }
        </VelocityTransitionGroup>
      </div>
    );
  }
}
