'use strict';

import React, { Component } from 'react';
import Card, { Block } from '@dux/element-card';
import Markdown from '@dux/element-markdown';
import eusa from './eusa.js';

export default class extends Component {
  render() {
    return (
      <div className='row'>
        <div className='small-8 columns small-centered'>
          <Card>
            <Block>
              <Markdown>
                {eusa.md}
              </Markdown>
            </Block>
          </Card>
        </div>
      </div>
    );
  }
}
