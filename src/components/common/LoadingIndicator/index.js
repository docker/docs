import React, { Component } from 'react';
import CircularProgress from 'material-ui/CircularProgress';

// http://www.material-ui.com/#/components/circular-progress for proptypes
export default class LoadingIndicator extends Component {
  render() {
    // $color-malibu
    return <CircularProgress {...this.props} color="#75ccfa" />;
  }
}
