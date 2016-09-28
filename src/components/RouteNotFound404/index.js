import React, { Component } from 'react';
import css from './styles.css';

// Note: must include image like this so that webpack picks it up
import image404 from 'lib/images/404@2x.png';

/* eslint-disable max-len */

export default class RouteNotFound404 extends Component {
  render() {
    return (
      <div className={css.wrapper}>
        <div>
          <img
            className={css.img}
            alt="404 Route Not Found"
            src={image404}
          />
        </div>
      </div>
    );
  }
}

/* eslint-enable */
