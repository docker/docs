import React, { Component } from 'react';
import { find } from 'lodash';
import { DockerImage } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

const IMAGES = [
  'library/wordpress:latest',
  'tutum/wordpress:staging',
  'tutum/wordpress-stackable:v0.1.0',
  'tutum/centos:some-tag',
  'library/centos:other',
  'tutum/ubuntu:trusty',
  'tutum/jboss',
  'library/tomcat',
  'tutum/tomcat',
  'tutum/sync',
  'tutum/riak',
  'library/redis',
  'tutum/redis',
  'library/rabbitmq',
  'tutum/rabbitmq',
  'library/postgres',
  'tutum/postgresql',
  'tutum/newrelic-agent',
  'library/mysql',
  'tutum/myql',
  'library/mongo',
  'tutum/mongodb',
  'library/memcached',
  'tutum/memcached',
  'library/mariadb',
  'tutum/mariadb',
  'library/joomla',
  'tutum/joomla',
  'tutum/influxdb',
  'library/haproxy',
  'tutum/haproxy',
  'library/glassfish',
  'tutum/glassfish',
  'library/fedora',
  'tutum/fedora',
  'library/elasticsearch',
  'tutum/elasticsearch',
  'library/drupal',
  'tutum/drupal',
  'library/debian',
  'tutum/debian',
  'tutum/datadog-agent',
  'tutum/couchdb',
  'library/cassandra',
  'tutum/cassandra',
  'tutum/btsync',
  'tutum/wikimedia',
  'someother/image:latest',
];

const HANDPCIKED = IMAGES.filter((name) =>
  find([
    'tutum/ubuntu',
    'library/redis',
    'someother',
    'tutum/elasticsearch',
    'tutum/cassandra',
  ], (target) => name.indexOf(target) >= 0)
);

@asExample(mdHeader, mdApi)
export default class DockerImageDoc extends Component {
  render() {
    return (
      <div>
        <h3>Combinations</h3>
        <div className="clearfix">
          {IMAGES.map((name) => (
            <div style={{
              float: 'left',
              width: '33%',
              marginBottom: '8px',
              textOverflow: 'ellipsis',
              overflow: 'hidden',
              whiteSpace: 'nowrap',
            }} key={name.replace(/\//g, '_')}
            >
              <DockerImage imageName={name} />
            </div>
          ))}
        </div>
        <h3>Standalone</h3>
        <div className="clearfix">
          <span style={{ float: 'left', marginRight: '20px' }}>Regular:</span>
          {HANDPCIKED.map(
            (name) => <div
              key={name.replace(/\//g, '_')}
              style={{ float: 'left', minWidth: '80px' }}
            >
              <DockerImage standalone imageName={name} />
            </div>
            )}
        </div>
        <div className="clearfix">
          <span style={{ float: 'left', marginRight: '32px' }}>Large:</span>
          {HANDPCIKED.map(
            (name) => <div
              key={name.replace(/\//g, '_')}
              style={{ float: 'left', minWidth: '80px' }}
            >
              <DockerImage standalone size="large" imageName={name} />
            </div>
          )}
        </div>
        <div className="clearfix">
          <span style={{ float: 'left', marginRight: '22px' }}>Xlarge:</span>
          {HANDPCIKED.map(
            (name) => <div
              key={name.replace(/\//g, '_')}
              style={{ float: 'left', minWidth: '80px' }}
            >
              <DockerImage standalone size="xlarge" imageName={name} />
            </div>
          )}
        </div>
        <h3>With Text</h3>
        <DockerImage imageName="library/ubuntu:trusty" />
        <DockerImage size="large" imageName="library/ubuntu:trusty" />
        <DockerImage size="xlarge" imageName="library/ubuntu:trusty" />
      </div>
    );
  }
}
