'use strict';
import {
  Autobuilds
  } from 'hub-js-sdk';
import has from 'lodash/object/has';
import merge from 'lodash/object/merge';
import sortBy from 'lodash/collection/sortBy';
import async from 'async';
var debug = require('debug')('hub:actions:savePushAutoBuildSettings');

function saveTag({JWT, namespace, name}) {
  return (tag, callback) => {
    if (tag.isNew && !tag.toDelete) { // toDelete unecessary now as tags with isNew && toDelete SHOULD be removed
      var ab = merge({}, tag, {repoName: name, namespace: namespace});
      Autobuilds.createAutomatedBuildTags(JWT, ab, function(err, res) {
        if (err) {
          debug(err);
          tag.error = true;
          callback(null, { tag, err });
        } else {
          callback(null, res.body);
        }
      });
    } else if (tag.toDelete && !!tag.id) {
      Autobuilds.deleteAutomatedBuildTags(JWT, namespace, name, tag.id, function(err, res){
        if (err) {
          debug(err);
          tag.error = true;
          callback(null, { tag, err });
        } else {
          callback(null, null);
        }
      });
    } else {
      Autobuilds.updateAutomatedBuildTags(JWT, namespace, name, tag.id, tag, function(err, res){
        if (err) {
          debug(err);
          tag.error = true;
          callback(null, { tag, err });
        } else {
          callback(null, res.body);
        }
      });
    }
  };
}

export default function(actionContext, {JWT, namespace, name, tags}, done) {
  let encounteredError = false;
  async.map(tags, saveTag({ JWT, namespace, name }), function(error, results) {
    results.forEach( (resultTag) => {
      if (has(resultTag, 'tag')) {
        const { tag, err } = resultTag;
        //We are keeping track of only `update` docker tag errors and at the moment we get them in `non_field_errors`
        const nonFieldErrors = err.response.body.non_field_errors;
        if (has(tag, 'error')) {
          actionContext.dispatch('SAVE_BUILD_TAGS_ERROR', {name: tag.name, error: nonFieldErrors});
          encounteredError = true;
        }
      }
    });
    if (!encounteredError) {
      //Gets here only if there are no errors
      actionContext.dispatch('SAVE_BUILD_TAGS_SUCCESS');
      const sorted = sortBy(results, 'id');
      actionContext.dispatch('UPDATE_AUTO_BUILD_SETTINGS', {
        field: 'autoBuildStore',
        key: 'build_tags',
        value: sorted
      });
    }
    done();
  });
}
