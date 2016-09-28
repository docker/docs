import React, { Component, PropTypes } from 'react';
// import './styles.css';
import marked from 'marked';
const { string } = PropTypes;
export default class Markdown extends Component {
  static propTypes = {
    className: string,
    rawMarkdown: string,
  }

  render() {
    const { className = '', rawMarkdown } = this.props;
    const HTML = {
      __html: marked(rawMarkdown, { sanitize: true }),
    };
    return (
      <div className="dMarkdown">
        <div className={className} dangerouslySetInnerHTML={HTML}>
        </div>
      </div>
    );
  }

}
