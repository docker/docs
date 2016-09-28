import React, { Component } from 'react';
import { TagList } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class TagListDoc extends Component {
  render() {
    return (
      <div className="clearfix">
        <h3>Simple tag list</h3>
        <TagList list={['tag1', 'tag2', 'tag3', 'tag4']} size={380} />
        <h3>Truncated tag list</h3>
        <TagList
          list={[
            'tag1',
            'tag2',
            'tag3',
            'loooooooooong tag',
            'tag4',
            'tag5',
            'tag6',
            'tag7',
            'supertag',
            'tag8',
          ]}
          size={380}
        />
      </div>
    );
  }
}
