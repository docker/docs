'use strict';

import React from 'react';
import provideContext from 'fluxible-addons-react/provideContext';
const debug = require('debug')('hub:bootstrapCreateElement');
/**
 * We create a react element for the Router using this function, passing it to
 * the `createElement` (func) property in <Router />. This will enable us to
 * provide context for the Fluxible app and adds some context functions like
 * `executeAction` and `getStore` to all components in our Routes.
 *
 * context is a Fluxible Context.
 */
export default (context) => {
  return (component, props) => {
    debug('setting context');
    props.context = context.getComponentContext();
    props.JWT = props.JWT || (props.cookies && props.cookies.JWT) || false;
    return React.createElement(provideContext(component), props);
  };
};
