'use strict';

var createStore = require('fluxible/addons/createStore');
import {
  ATTEMPTING,
  DEFAULT,
  FACEPALM,
  SUCCESSFUL
} from './deletepipelinestore/Constants';

var debug = require('debug')('SignupStore');

export default createStore({
  storeName: 'DeletePipelineStore',
  handlers: {
    DELETE_PIPELINE_ATTEMPTING: '_start',
    DELETE_PIPELINE_FACEPALM: '_facepalm',
    DELETE_PIPELINE_SUCCESS: '_success'
  },
  initialize() {
    this.STATUS = DEFAULT;
  },
  _start() {
    this.STATUS = ATTEMPTING;
    this.emitChange();
  },
  _facepalm() {
    this.STATUS = FACEPALM;
    this.emitChange();
  },
  _success() {
    this.STATUS = SUCCESSFUL;
    this.emitChange();
  },
  getState() {
    return {
      STATUS: this.STATUS
    };
  },
  dehydrate() {
    return {};
  },
  rehydrate(state) {
    this.state = state;
  }
});
