import React, { Component } from 'react';
import svg from './_Svg';

@svg
class AddIcon extends Component {
  render() {
    return (
      <g>
        <path fill="none" d="M12 4c-1.104 0-2 .896-2 2v3h4V6c0-1.104-.896-2-2-2z" />
        <path fill="none" d="M12 4c-1.104 0-2 .896-2 2v3h4V6c0-1.104-.896-2-2-2z" />
        <path d="M17 9h-1V6c0-2.21-1.79-4-4-4S8 3.79 8 6v3H7c-1.104 0-2 .896-2 2v8c0 1.104.896 2 2 2h10c1.104 0 2-.896 2-2v-8c0-1.104-.896-2-2-2zm-5 8c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm2-8h-4V6c0-1.104.896-2 2-2s2 .896 2 2v3z" />
      </g>
    );
  }
}

export default AddIcon;
