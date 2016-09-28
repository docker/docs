'use strict';
import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';

var AutobuildTagsStore = createStore({
  storeName: 'AutobuildTagsStore',
  handlers: {
    AUTOBUILD_TAGS_ERROR: '_autobuildTagsError'
  },
  initialize: function() {
    //tags is an array of tag
    //{ dockerfile_location, source_type['Tag' or 'Branch'], source_name[eg. master] }
    this.tags = [];
  },
  addTag: function(tag) {
    //tag
    //{
    // id: 'row-1'
    // sourceName: 'master'
    // fileLocation: '/'
    // buildTag: 'latest',
    // sourceType: 'Branch'
    //}
    this.tags.push(tag);
  },
  removeTag: function(id) {
    _.remove(this.tags, function(tag) {
      return (tag.id === id);
    });
  },
  setTagState: function(id, state) {
    var tagToUpdate = _.find(this.tags, function(tag) {
      return tag.id === id;
    });
    _.merge(tagToUpdate, state);
  },
  getState: function() {
    return {
      tags: this.tags
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.tags = state.tags;
  }
});

module.exports = AutobuildTagsStore;
