import React, { Component } from 'react';
import { NodeProvider } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

const PROVIDERS = [
  'digitalocean',
  'azure',
  'aws',
  'packet',
  'softlayer',
];

@asExample(mdHeader, mdApi)
export default class NodeProviderDoc extends Component {
  render() {
    return (
      <div>
        <h3>Combinations</h3>
        <div className="clearfix">
          {PROVIDERS.map((name) => (
            <div style={{
              float: 'left',
              width: '33%',
              marginBottom: '8px',
              textOverflow: 'ellipsis',
              overflow: 'hidden',
              whiteSpace: 'nowrap',
            }} key={name}
            >
              <NodeProvider providerName={name} />
            </div>
          ))}
        </div>
        <h3>Standalone</h3>
        <div className="clearfix">
          <span style={{ float: 'left', marginRight: '20px' }}>Regular:</span>
          {PROVIDERS.map(
            (name) => <div
              key={name}
              style={{ float: 'left', minWidth: '80px' }}
            >
              <NodeProvider standalone providerName={name} />
            </div>
            )}
        </div>
        <div className="clearfix">
          <span style={{ float: 'left', marginRight: '32px' }}>Large:</span>
          {PROVIDERS.map(
            (name) => <div
              key={name.replace(/\//g, '_')}
              style={{ float: 'left', minWidth: '80px' }}
            >
              <NodeProvider standalone size="large" providerName={name} />
            </div>
          )}
        </div>
        <div className="clearfix">
          <span style={{ float: 'left', marginRight: '22px' }}>Xlarge:</span>
          {PROVIDERS.map(
            (name) => <div
              key={name.replace(/\//g, '_')}
              style={{ float: 'left', minWidth: '80px' }}
            >
              <NodeProvider standalone size="xlarge" providerName={name} />
            </div>
          )}
        </div>
      </div>
    );
  }
}
