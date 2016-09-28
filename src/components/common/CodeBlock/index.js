import React, { Component, PropTypes } from 'react';
import css from './styles.css';

export default class CodeBlock extends Component {
  static propTypes = {
    className: PropTypes.string,
    children: PropTypes.node.isRequired,
  }

  render() {
    const { className = '' } = this.props;

    return (
      <div className={`${css.codeblock} ${className}`}>
        {this.props.children}
      </div>
    );
  }
}
