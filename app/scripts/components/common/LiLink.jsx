'use strict';

import React from 'react';
const debug = require('debug')('LiLink: ');
import { Link } from 'react-router';

export default class LiLink extends Link {
  isActive() {
    const { to, query, onlyActiveOnIndex } = this.props;
    return this.context.history.isActive(to, query, onlyActiveOnIndex);
  }
//In order to show active state on nav links, we need the active class on the li element
//Dependant on Foundation .active class
  render() {
    return <li className={this.isActive() ? 'active' : ''}>{super.render()}</li>;
  }
}
