'use strict';
import createStore from 'fluxible/addons/createStore';
import last from 'lodash/array/last';
import {
  DEFAULT,
  ATTEMPTING,
  ERROR
} from './addwebhookformstore/Constants';
const debug = require('debug')('AddWebhookFormStore');

var WebhooksSettingsStore = createStore({
  storeName: 'AddWebhookFormStore',
  handlers: {
    RECEIVE_WEBHOOKS: '_receiveWebhooks',
    ADD_WEBHOOK_CLEAR: '_clear',
    ADD_WEBHOOK_START: '_start',
    ADD_WEBHOOK_RESET: '_reset',
    ADD_WEBHOOK_SUCCESS: '_success',
    ADD_WEBHOOK_NEW_HOOK: '_newHook',
    ADD_WEBHOOK_REMOVE_HOOK: '_removeHook',
    ADD_WEBHOOK_ERROR: '_error',
    ADD_WEBHOOK_MISSING_ARGS: '_handleMissingArgs',
    ADD_WEBHOOK_VALIDATION_ERRORS: '_handleValidationErrors'
  },
  initialize() {
    /**
     * hookFields represent each `input` pairing that is
     * rendered. They contain no data about the content of the input
     */
    this.hookFields = [1];
    this.STATUS = DEFAULT;
    this.serverErrors = {};
  },
  _error(args) {
    // TODO: handle generic error
    this.STATUS = ERROR;
    this.serverErrors = args;
    this.emitChange();
  },
  _handleMissingArgs(args) {
    debug('missing args: ', args);
    this.serverErrors = args;
  },
  _handleValidationErrors(args) {
    debug('validation errors: ', args);
    this.serverErrors = args;
  },
  _newHook() {
    const { hookFields: fields } = this;
    this.hookFields = fields.concat(last(fields) + 1);
    this.emitChange();
  },
  _reset() {
    this.initialize();
    this.emitChange();
  },
  _start() {
    this.STATUS = ATTEMPTING;
    this.emitChange();
  },
  _success() {},
  _receiveWebhooks(payload) {
    debug(payload);
    this.pipelines = payload.results;
    this.emitChange();
  },
  _receiveAddWebhookErrors(error) {
  },
  getState() {
    return {
      STATUS: this.STATUS,
      pipelines: this.pipelines,
      hookFields: this.hookFields,
      serverErrors: this.serverErrors
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.pipelines = state.pipelines;
    this.hookFields = state.hookFields;
  }
});

module.exports = WebhooksSettingsStore;
