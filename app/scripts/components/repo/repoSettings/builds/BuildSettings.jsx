'use strict';

import React, {
  PropTypes,
  Component
  } from 'react';
const { string, object, func } = PropTypes;
import findIndex from 'lodash/array/findIndex';
import Card, { Block } from '@dux/element-card';
import FA from 'common/FontAwesome';
import updateAutoBuildSettings from 'actions/updateAutoBuildSettings.js';
import { FlexRow, FlexItem, FlexHeader, FlexTable } from 'common/FlexTable';
import SourceRepositoryCard from 'common/SourceRepositoryCard';
import styles from './BuildSettings.css';

export default class BuildSettings extends Component {

  static propTypes = {
    JWT: string.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    autoBuildStore: object.isRequired
  };

  static contextTypes = {
    executeAction: func.isRequired
  };

  state = {
    isContextualHelpOpen: false
  };

  _activeClick = (e) => {
    e.preventDefault();
    this.context.executeAction(updateAutoBuildSettings, {
      JWT: this.props.JWT,
      name: this.props.name,
      namespace: this.props.namespace,
      data: {active: !this.props.autoBuildStore.active}
    });
  }

  toggleContextualHelp = () => {
    if (this.state.isContextualHelpOpen) {
      this.setState({
        isContextualHelpOpen: false
      });
    } else {
      this.setState({
        isContextualHelpOpen: true
      });
    }
  }

  renderContextualHelp = () => {
    return (
      <div className={styles.lessLabel}>
        Here are a few examples:
        <br/>
        <FlexTable>
          <FlexHeader>
            <FlexItem>Scenario</FlexItem>
            <FlexItem>Type</FlexItem>
            <FlexItem>Name</FlexItem>
            <FlexItem>Docker Tag Name</FlexItem>
            <FlexItem>Matches</FlexItem>
            <FlexItem>Docker Tag Built</FlexItem>
          </FlexHeader>
          <FlexRow>
            <FlexItem>Exact match</FlexItem>
            <FlexItem>Branch</FlexItem>
            <FlexItem>master</FlexItem>
            <FlexItem>latest</FlexItem>
            <FlexItem>master</FlexItem>
            <FlexItem>latest</FlexItem>
          </FlexRow>
          <FlexRow>
            <FlexItem>Match versions</FlexItem>
            <FlexItem>Tag</FlexItem>
            <FlexItem>/^[0-9.]+$/</FlexItem>
            <FlexItem>{'release-{sourceref}'}</FlexItem>
            <FlexItem>1.2.0</FlexItem>
            <FlexItem>release-1.2.0</FlexItem>
          </FlexRow>
          <FlexRow>
            <FlexItem>Trailing modifiers</FlexItem>
            <FlexItem>Tag</FlexItem>
            <FlexItem>/^[0-9.]+/</FlexItem>
            <FlexItem>{'release-{sourceref}'}</FlexItem>
            <FlexItem>1.2.0-rc</FlexItem>
            <FlexItem>release-1.2.0-rc</FlexItem>
          </FlexRow>
        </FlexTable>
      </div>
    );
  }

  render() {
    const {
      autoBuildStore,
      JWT,
      name,
      namespace
    } = this.props;
    let {
      provider,
      repo_web_url: url
    } = autoBuildStore;

    let contextualHelp;
    if (this.state.isContextualHelpOpen) {
      contextualHelp = this.renderContextualHelp();
    }

    let linkText = `Source Project`;
    if ( url ) {
      const URLportions = url.split('/');
      const domainIndex = findIndex(URLportions, function (str) {
        return str === 'github.com' || str === 'bitbucket.org';
      });

      if (domainIndex > 0 && URLportions.length > (domainIndex + 2)) {
        const repoUsername = URLportions[domainIndex + 1];
        const repoName = URLportions[domainIndex + 2];
        linkText = `${repoUsername}/${repoName}`;
      }
    } else {
      url = '#';
      linkText = 'Not Available';
    }

    let icon = '';
    switch (provider.toLowerCase()) {
      case 'github':
        icon = 'fa-github';
        break;
      case 'bitbucket':
        icon = 'fa-bitbucket';
        break;
      default:
        icon = 'fa-link';
    }

    return (
      <div>
        <Card heading='Build Settings'>
          <Block>
            <div className={'row ' + styles.settingsWrapper}>
              <div className={'columns large-9 ' + styles.settingsBox}>
                <span className={styles.label}>
                  <label className={styles.hover}>
                    <input type="checkbox"
                           name="build-active"
                           checked={autoBuildStore.active}
                           onChange={this._activeClick} />  When active, builds will happen automatically on pushes.
                  </label>
                  {`The build rules below specify how to build your source into Docker images.
                  The name can be a string or a regex. The Docker Tag name may contain variables.
                  We currently support {sourceref}, which refers to the source branch/tag name.`}
                  {!this.state.isContextualHelpOpen ? (<a onClick={this.toggleContextualHelp}> Show more </a>) : (<a onClick={this.toggleContextualHelp}> Show less </a>)}
                </span>
              </div>
              <div className={'columns large-3 ' + styles.sourceBox}>
                <div className={styles.centerSource}>
                  <FA icon={icon} size='2x'/>
                  <div className={styles.source}>
                    <span>Source Repository</span><br/>
                    <a href={url}>{linkText}</a>
                  </div>
                </div>
              </div>
            </div>
          </Block>
        </Card>
        {contextualHelp}
      </div>
    );
  }
}
