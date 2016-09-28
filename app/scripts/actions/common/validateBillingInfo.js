'use strict';

export default function({storePrefix}){
  return (actionContext, fieldErrors) => {
    actionContext.dispatch(`${storePrefix}_ERRORS`, fieldErrors);
  };
}
