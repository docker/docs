'use strict';

import React, { PropTypes, Component } from 'react';
const { string } = PropTypes;
import Card, { Block } from '@dux/element-card';
import findIndex from 'lodash/array/findIndex';
import FA from 'common/FontAwesome';

export default class SourceRepositoryCard extends Component {

  static propTypes = {
    url: string.isRequired,
    provider: string.isRequired
  }

  render() {
    let { provider, url } = this.props;
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
      <Card heading='Source Repository'>
        <Block>
          <span>
            <FA icon={icon} />
            <span>&nbsp;&nbsp;</span>
            <a href={url}>{linkText}</a>
          </span>
        </Block>
      </Card>
    );
  }
}
