'use strict';

import React, { Component } from 'react';
import PageHeader from 'components/common/pageHeader';
import styles from './api.css';

// begin stackoverflow hack http://stackoverflow.com/questions/4003823/javascript-getcookie-functions/4004010#4004010
function getCookies() {
    var arrMap = function(arr, callback, thisArg) {
        for (var i = 0, n = arr.length, a = []; i < n; i++) {
            if (i in arr) {
                a[i] = callback.call(thisArg, arr[i]);
            }
        }
        return a;
    };
    var c = document.cookie, v = 0, cookies = {};
    if (document.cookie.match(/^\s*\$Version=(?:"1"|1);\s*(.*)/)) {
        c = RegExp.$1;
        v = 1;
    }
    if (v === 0) {
        arrMap(c.split(/[,;]/), function(cookie) {
            var parts = cookie.split(/=/, 2),
                name = decodeURIComponent(parts[0].replace(/^\s+/, '')),
                value = parts.length > 1 ? decodeURIComponent(parts[1].replace(/\s+$/, '')) : null;
            cookies[name] = value;
        });
    } else {
        arrMap(c.match(/(?:^|\s+)([!#$%&'*+\-.0-9A-Z^`a-z|~]+)=([!#$%&'*+\-.0-9A-Z^`a-z|~]*|"(?:[\x20-\x7E\x80\xFF]|\\[\x00-\x7F])*")(?=\s*[,;]|$)/g), function($0, $1) {
            var name = $0,
                value = $1.charAt(0) === '"'
                          ? $1.substr(1, -1).replace(/\\(.)/g, '$1')
                          : $1;
            cookies[name] = value;
        });
    }
    return cookies;
}
function getCookie(name) {
    return getCookies()[name];
}
// end stackoverflow hack

export default class API extends Component {
  pass_csrf(e) {
    const csrfToken = getCookie('csrftoken');
    e.target.contentWindow.postMessage(csrfToken, '*');
    e.target.style.height = (document.body.scrollHeight - 150) + 'px';
  }

  render() {
      return (
        <div>
          <PageHeader title='API' />
          <div className={ styles.row }>
              <iframe className={ styles.container } src='/api/docs/' onLoad={ this.pass_csrf } />
          </div>
        </div>
      );
  }
}
