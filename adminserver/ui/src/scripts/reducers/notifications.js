/*global notifications*/
'use strict';

import assign from 'object-assign';
import consts from 'consts';

const actions = {
  [consts.ADD_GROWL_NOTIFICATION]: (state, data) => {
    // Remove any existing notification with the same ID.
    // This may happen when we create a notification for saving yaml settings
    // twice in succession.
    let copy = actions[consts.REMOVE_GROWL_NOTIFICATION](state, data);
    copy.growl.push(data.payload);
    return copy;
  },
  [consts.UPDATE_GROWL_NOTIFICATION]: (state, data) => {
    return assign({}, state, {
      growl: state.growl.slice().map( item => {
        if (item.id !== data.payload.id) {
          return item;
        }
        return assign({}, item, data.payload);
      })
    });
  },
  [consts.REMOVE_GROWL_NOTIFICATION]: (state, data) => {
    return assign({}, state, {
        growl: state.growl.slice().filter( item => { return item.id !== data.payload.id; })
    });
  },
  [consts.ADD_BANNER_NOTIFICATION]: (state, data) => {
    // Remove any existing notification with the same ID.
    // twice in succession.
    const copy = actions[consts.REMOVE_BANNER_NOTIFICATION](state, data);
    copy.banner.push(data.payload);
    return copy;
  },
  [consts.UPDATE_BANNER_NOTIFICATION]: (state, data) => {
    return assign({}, state, {
      banner: state.banner.slice().map( item => {
        if (item.id !== data.payload.id) {
          return item;
        }
        return assign({}, item, data.payload);
      })
    });
  },
  [consts.REMOVE_BANNER_NOTIFICATION]: (state, data) => {
    return assign({}, state, {
        banner: state.banner.slice().filter( item => { return item.id !== data.payload.id; })
    });
  }
};

function serverNotifications() {
  if (notifications) {
    const bannerNotifications = notifications.map((item, idx) => {
        if (!item.id) {
            return assign(item, {id: `id${idx}`});
        }
        return item;
    });
    return {
        growl: [],
        banner: bannerNotifications
    };
  }
  return {
      growl: [],
      banner: []
  };
}

export default function(state = serverNotifications(), data) {
  if (typeof actions[data.type] === 'function') {
    return actions[data.type](state, data);
  }
  return state;
}
