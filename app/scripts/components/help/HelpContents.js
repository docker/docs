'use strict';

import styles from './Help.css';

export const HelpContents = [{
  title: 'Getting started',
  icon: 'fa-power-off',
  body: [{
    href: 'https://docs.docker.com/engine/understanding-docker/',
    text: 'Learn the fundamentals of Docker'
  }, {
    href: 'https://docs.docker.com/engine/getstarted/step_one/',
    text: 'Install Docker'
  }, {
    href: 'https://docs.docker.com/engine/tutorials/dockerizing/',
    text: 'Hello world in a container'
  }]
}, {
  title: 'Docker Hub Docs',
  icon: 'fa-file-code-o',
  body: [{
    href: 'https://docs.docker.com/docker-hub/',
    text: 'Overview of Docker Hub'
  }, {
    href: 'https://docs.docker.com/docker-hub/repos/',
    text: 'Create repositories'
  }, {
    href: 'https://docs.docker.com/docker-hub/builds/',
    text: 'Automated builds'
  }, {
    href: 'https://docs.docker.com/docker-hub/official_repos/',
    text: 'Official repositories'
  }]
}, {
  title: 'Docker Hub Forum',
  icon: 'fa-comments',
  body: [{
    text: 'Ask and answer questions with the Docker Hub Community'
  }, {
    text: 'Share ideas and feature requests'
  }],
  subtext: {
    text: 'Community Support',
    btnStyle: styles.helpCommunitySupport,
    link: 'https://github.com/docker/hub-feedback/issues?q=is%3Aissue+is%3Aopen+sort%3Acomments-desc'
  }
}, {
  title: 'Docker Support',
  icon: 'fa-phone',
  body: [{
    text: 'Email support for Docker Hub subscribers'
  }],
  subtext: {
    text: 'Paid Support',
    btnStyle: styles.helpPaidSupport + ' button primary',
    link: 'https://support.docker.com/'
  }
}];
