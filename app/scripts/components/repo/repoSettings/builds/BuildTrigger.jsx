'use strict';

import React, {
  PropTypes,
  createClass
  } from 'react';
import { Link } from 'react-router';
import updateAutoBuildSettingsStore from 'actions/updateAutoBuildSettingsStore.js';
import toggleTriggerStatus from 'actions/toggleTriggerStatus.js';
import regenTriggerToken from 'actions/regenTriggerToken.js';
import _ from 'lodash';
import { FlexTable, FlexRow, FlexHeader, FlexItem } from '../../../common/FlexTable.jsx';
import Button from '@dux/element-button';
import moment from 'moment';
import FA from '../../../common/FontAwesome.jsx';
import Markdown from '@dux/element-markdown';
import Card, { Block } from '@dux/element-card';
import SimpleInput from 'common/SimpleInput';

import styles from './BuildTrigger.css';

var debug = require('debug')('BuildOptions');

var _mkTriggerLog = function(username, reponame) {
  return (log) => {
    let buildRequest = 'null';
    if (log.build_code) {
      buildRequest = (
        <Link to={`/r/${username}/${reponame}/builds/${log.build_code}/`}>
          {log.build_code}
        </Link>
      );
    }
    return (
      <FlexRow>
        <FlexItem>{moment.utc(log.created).format('MMMM Do, YYYY, h:mm a')}</FlexItem>
        <FlexItem>{log.ip_address}</FlexItem>
        <FlexItem>{log.result}</FlexItem>
        <FlexItem>{log.result_desc}</FlexItem>
        <FlexItem>{log.request_body}</FlexItem>
        <FlexItem>
          { buildRequest }
        </FlexItem>
      </FlexRow>
    );
  };
};

var TriggerStatus = createClass({
  propTypes: {
    triggerStatus: PropTypes.shape({
      token: PropTypes.string,
      trigger_url: PropTypes.string,
      active: PropTypes.bool.isRequired
    }),
    toggleActive: PropTypes.func.isRequired,
    regenToken: PropTypes.func.isRequired
  },
  getInitialState: function() {
    return {
      showExamples: false
    };
  },
  _toggleExamples(e) {
    this.setState({
      showExamples: !this.state.showExamples
    });
  },
  render() {
    if (!this.props.triggerStatus.active) {
      return (
        <Button variant='primary' size='tiny' onClick={this.props.toggleActive}>
          Activate Triggers
        </Button>
      );
    }
    const button = (
      <Button size='tiny' onClick={this.props.toggleActive}>
        Deactivate Triggers
      </Button>
    );
    const regen = (
      <Button variant='primary' size='tiny' onClick={this.props.regenToken}>
        <FA icon="fa-refresh" />
        &nbsp; Regenerate token
      </Button>
    );
    const { showExamples } = this.state;
    let examples;
    if (showExamples) {
      examples = (
        <div>
          <div className={styles.label}>Examples</div>
          <Examples triggerStatus={this.props.triggerStatus}/>
        </div>
      );
    }
    const examplesToggle = (
      <p>
        Trigger endpoints are activated. Use the trigger token or URL below in
        your requests.
        <span className={styles.examplesToggle} onClick={this._toggleExamples}>
          { showExamples ? ' Hide examples.' : ' Show examples.' }
        </span>
      </p>
    );

    let triggerToken = null;
    if (this.props.triggerStatus.token) {
      triggerToken = (
        <div>
          <div className={styles.label}>Trigger Token</div>
          <div className={'row ' + styles.status}>
            <div className='large-6 columns'>
              <SimpleInput value={this.props.triggerStatus.token} readOnly />
            </div>
            <div className="large-6 columns">
              {regen}
            </div>
          </div>
        </div>
      );
    }
    let triggerUrl = null;
    if (this.props.triggerStatus.trigger_url) {
      triggerUrl = (
        <div>
          <div className={styles.label}>Trigger URL</div>
          <div className={'row ' + styles.status}>
            <div className='large-9 columns'>
              <SimpleInput className={styles.triggerValue}
                     value={this.props.triggerStatus.trigger_url}
                     readOnly />
            </div>
          </div>
        </div>
      );
    }
    return (
      <div>
        <p>Note: Build requests are throttled so that they don't
        overload the system. If there is already a build request
        pending, the request will be ignored.</p>
        {examplesToggle}
        {button}
        <div className={'row ' + styles.triggerTitle}>
          <div className="columns large-12">
            {triggerToken}
            {triggerUrl}
            {examples}
          </div>
        </div>
      </div>
    );
  }
});

// Currently logs will be static
var TriggerLogs = createClass({
  PropTypes: {
    triggerLogs: PropTypes.array.isRequired
  },
  _makeLogs() {
    if (this.props.triggerLogs.length === 0) {
      return (
        <FlexRow>
          <FlexItem>
            No logs to show
          </FlexItem>
        </FlexRow>
      );
    } else {
      return (
        <div>{_.map(this.props.triggerLogs, _mkTriggerLog(this.props.username, this.props.reponame))}</div>
      );
    }
  },
  render() {
    return (
      <div className={'columns large-12 ' + styles.logs}>
        <br />
        <h5>Last 10 Trigger Logs</h5>
        <FlexTable>
          <FlexHeader>
            <FlexItem>Date/Time</FlexItem>
            <FlexItem>IP Address</FlexItem>
            <FlexItem>Status</FlexItem>
            <FlexItem>Status Description</FlexItem>
            <FlexItem>Request Body</FlexItem>
            <FlexItem>Build Request</FlexItem>
          </FlexHeader>
          {this._makeLogs()}
        </FlexTable>
      </div>
    );
  }
});

var Examples = createClass({
  propTypes: {
    triggerStatus: PropTypes.object.isRequired
  },
  render() {
    return (
      <Card>
        <Block>
          <div className="row">
            <div className={'columns large-12 ' + styles.content }>
              <Markdown>{`\`\`\`shell
# Trigger all tags/branchs for this automated build.
$ curl -H "Content-Type: application/json" --data '{"build": true}' -X POST ${this.props.triggerStatus.trigger_url}

# Trigger by docker tag name
$ curl -H "Content-Type: application/json" --data '{"docker_tag": "master"}' -X POST ${this.props.triggerStatus.trigger_url}

# Trigger by Source branch named staging
$ curl -H "Content-Type: application/json" --data '{"source_type": "Branch", "source_name": "staging"}' -X POST ${this.props.triggerStatus.trigger_url}

# Trigger by Source tag named v1.1
$ curl -H "Content-Type: application/json" --data '{"source_type": "Tag", "source_name": "v1.1"}' -X POST ${this.props.triggerStatus.trigger_url}
              \`\`\``}</Markdown>
            </div>
          </div>
        </Block>
      </Card>
    );
  }
});

var BuildTrigger = createClass({
  displayName: 'BuildTrigger',
  propTypes: {
    autoBuildStore: PropTypes.object.isRequired,
    triggerStatus: PropTypes.object.isRequired,
    triggerLogs: PropTypes.array.isRequired,
    namespace: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    JWT: PropTypes.string.isRequired
  },
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  _toggleActive(e) {
    e.preventDefault();
    this.context.executeAction(toggleTriggerStatus, {
      JWT: this.props.JWT,
      namespace: this.props.namespace,
      name: this.props.name,
      active: !this.props.triggerStatus.active
    });
  },
  _regenToken(e) {
    e.preventDefault();
    this.context.executeAction(regenTriggerToken, {
      JWT: this.props.JWT,
      namespace: this.props.namespace,
      name: this.props.name
    });
  },

  render() {
    const { name,
            namespace,
            triggerLogs,
            triggerStatus } = this.props;

    const subTitle = (
        <p>Trigger your Automated Build by sending a POST to a specific endpoint.</p>
    );
    return (
      <Card heading='Build Triggers'>
        <Block>
          {subTitle}
          <TriggerStatus triggerStatus={triggerStatus}
                         toggleActive={this._toggleActive}
                         regenToken={this._regenToken}/>
          <div className='row'>
            <TriggerLogs triggerLogs={triggerLogs}
                         username={namespace}
                         reponame={name}/>
          </div>
        </Block>
      </Card>
    );
  }
});

export default BuildTrigger;
