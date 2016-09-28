import React, { Component, PropTypes } from 'react';
import { List } from 'common';
import UnscannedTagRow from '../UnscannedTagRow';

const { array, func } = PropTypes;

export default class RepositoryTagsList extends Component {
  static propTypes = {
    generateTagPath: func,
    tags: array,
  }

  renderTag = (tag) => {
    const { name } = tag;
    return <UnscannedTagRow key={name} tag={tag} />;
  }

  render() {
    const { tags } = this.props;
    if (!tags) {
      return null;
    }
    return (
      <div>
        <List hover>
          { tags.map(this.renderTag) }
        </List>
      </div>
    );
  }
}
