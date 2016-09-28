import React, { Component, PropTypes } from 'react';
import {
  ListItem,
} from 'common';
import css from './styles.css';
import { formatBucketedNumber } from 'lib/utils/formatNumbers';
import routes from 'lib/constants/routes';
const {
  func,
  node,
  number,
  object,
  shape,
  string,
} = PropTypes;

export default class CommunityImageSearchResult extends Component {
  static propTypes = {
    categories: object,
    image: shape({
      popularity: number,
      name: string.isRequired,
      short_description: string,
    }),
    location: shape({
      pathname: string,
      query: object,
      state: node,
    }),
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  onClick = () => {
    const { name = '' } = this.props.image;
    const [namespace, reponame] = name.split('/');
    const pathname = routes.communityImageDetail({ namespace, reponame });
    this.context.router.push({ pathname });
  }

  renderPopularity(popularity) {
    if (!popularity) {
      return null;
    }
    const numPulls = formatBucketedNumber(popularity);
    const pullOrPulls = numPulls === '1' ? ' Pull' : ' Pulls';
    return (
      <div className={css.popularity}>
        {`${numPulls} ${pullOrPulls}`}
      </div>
    );
  }

  render() {
    const {
      popularity,
      name,
      short_description,
    } = this.props.image;
    return (
      <ListItem onClick={this.onClick}>
        <div className={css.communityRepo}>
          <div>
            <div className={css.name}>{name}</div>
            <div className={css.shortDescription}>{short_description}</div>
          </div>
          {this.renderPopularity(popularity)}
        </div>
      </ListItem>
    );
  }
}
