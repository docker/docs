import React, { Component } from 'react';
import css from './styles.css';

export default function asStaticBox(ComposedComponent) {
  return class StaticBox extends Component {
    render() {
      return (
        <div className={css.main}>
          <div className={css.container}>
            <div className={css.wrapper}>
              <ComposedComponent {...this.props} />
            </div>
          </div>
        </div>
      );
    }
  };
}
