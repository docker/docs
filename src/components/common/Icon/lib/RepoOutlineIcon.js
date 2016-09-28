import React, { Component, PropTypes } from 'react';
import svg from './_Svg';

const { bool } = PropTypes;

@svg(80)
class RepoOutlineIcon extends Component {
  static propTypes = {
    filled: bool,
  }
  render() {
    const { filled } = this.props;
    let fill1 = '#f7f8f9'; // $color-porcelain
    let fill2 = '#e0e4e7'; // $color-geyser

    if (filled) {
      fill1 = '#d1eefd';
      fill2 = '#1aaaf8';
    }
    return (
      <g>
        <path fill={fill1} d="M18.2 75.5c-3.7 0-6.8-3-6.8-6.7V11.2c0-3.7 3-6.7 6.8-6.7h43.5c3.7 0 6.8 3 6.8 6.7v57.6c0 3.7-3 6.7-6.8 6.7H18.2zm.3-7h43v-57h-43v57zm7-28v-7h17v7h-17zm0-15v-7h21v7h-21z" />
        <path fill={fill2} d="M47 18H25v8h22v-8zm-1 7H26v-6h20v6zM19 11h-1v58h44V11H19zm42 57H19V12h42v56zM43 33H25v8h18v-8zm-1 7H26v-6h16v6zM61.8 4H18.2c-4 0-7.2 3.2-7.2 7.2v57.6c0 4 3.2 7.2 7.2 7.2h43.5c4 0 7.2-3.2 7.2-7.2V11.2c.1-4-3.1-7.2-7.1-7.2zM68 68.8c0 3.4-2.8 6.2-6.2 6.2H18.2c-3.4 0-6.2-2.8-6.2-6.2V11.2C12 7.8 14.8 5 18.2 5h43.6c3.4 0 6.2 2.8 6.2 6.2v57.6z" />
      </g>
    );
  }
}

export default RepoOutlineIcon;
