// TODO: @camacho 2/9/16 remove in 1 week once new client has been loaded everywhere
'use strict';

import React, { PropTypes } from 'react';

module.exports = React.createClass({
  componentWillMount: (props) => window.location = 'https://www.docker.com/pricing',
  // Required for creating a component
  render: () => null
});
