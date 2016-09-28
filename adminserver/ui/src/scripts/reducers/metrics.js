'use strict';

import assign from 'object-assign';

import consts from 'consts';

const {
    metrics,
    loading
} = consts;

const actions = {
  [metrics.FETCH_METRICS]: (state, data) => {
    if (!data.ready) {
      return assign({}, state, {status: loading.PENDING});
    }

    if (data.error) {
      return assign({}, state, {status: loading.FAILURE, error: data.error});
    }

    return assign({}, state, {status: loading.SUCCESS, metrics: data.payload});

  },

  [metrics.OBSERVE_METRICS_DATA]: (state, data) => {
    return assign({}, state, {metrics: state.metrics.slice().concat(data.payload)});
  }
};

export default function metricStore(state = {}, data) {
  if (typeof actions[data.type] === 'function') {
    return actions[data.type](state, data);
  }
  return state;
}
