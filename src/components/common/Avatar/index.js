import React, { Component, PropTypes } from 'react';
import AvatarMui from 'material-ui/Avatar';

const { string, element, number } = PropTypes;

export default class Avatar extends Component {
  // Look for the rest at: http://www.material-ui.com/#/components/avatar
  // size is in `px`
  static propTypes = {
    src: string,
    icon: element,
    size: number,
  }

  render() {
    return (
      <AvatarMui {...this.props} />
    );
  }
}
