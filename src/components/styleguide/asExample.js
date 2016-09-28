import React, { Component } from 'react';
import marked from 'marked';
import 'github-markdown-css/github-markdown.css';

function asExample(mdHeader, mdApi) {
  return function config(ComposedComponent) {
    return class Example extends Component {
      render() {
        return (
          <div>
            <div className="markdown-body">
              <div
                dangerouslySetInnerHTML={{
                  __html: marked(mdHeader, { sanitize: true }),
                }}
              ></div>
            </div>
            <div style={{ marginBottom: '42px' }}>
              <ComposedComponent {...this.props} />
            </div>
            <div className="markdown-body">
              <div
                dangerouslySetInnerHTML={{
                  __html: marked(mdApi, { sanitize: true }),
                }}
              ></div>
            </div>
          </div>
        );
      }
    };
  };
}

export default asExample;
