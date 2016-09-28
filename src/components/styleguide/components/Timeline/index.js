import React, { Component } from 'react';
import { Timeline } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class TimelineDoc extends Component {
  render() {
    /* eslint-disable quotes, quote-props, max-len, comma-dangle */
    const actions = [
      {
        "action": "Service Redeploy",
        "can_be_canceled": false,
        "can_be_retried": false,
        "created": "Mon, 29 Feb 2016 10:36:25 +0000",
        "end_date": "Mon, 29 Feb 2016 10:36:37 +0000",
        "ip": "172.17.0.1",
        "location": "unknown",
        "method": "POST",
        "object": "/api/app/v1/service/573e735c-8f0c-4dc2-9836-1d14ed602cb2/",
        "path": "/api/app/v1/service/573e735c-8f0c-4dc2-9836-1d14ed602cb2/redeploy/",
        "resource_uri": "/api/audit/v1/action/1381b3f4-258f-49ef-9aaf-659b488f151c/",
        "start_date": "Mon, 29 Feb 2016 10:36:25 +0000",
        "state": "Success",
        "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/601.4.4 (KHTML, like Gecko) Version/9.0.3 Safari/601.4.4",
        "uuid": "1381b3f4-258f-49ef-9aaf-659b488f151c"
      },
      {
        "action": "Service Start",
        "can_be_canceled": false,
        "can_be_retried": false,
        "created": "Mon, 29 Feb 2016 09:38:34 +0000",
        "end_date": "Mon, 29 Feb 2016 09:38:46 +0000",
        "ip": "172.17.0.1",
        "location": "unknown",
        "method": "POST",
        "object": "/api/app/v1/service/573e735c-8f0c-4dc2-9836-1d14ed602cb2/",
        "path": "/api/app/v1/service/573e735c-8f0c-4dc2-9836-1d14ed602cb2/start/",
        "resource_uri": "/api/audit/v1/action/2b072082-b943-4dc0-9e18-3242bd40978a/",
        "start_date": "Mon, 29 Feb 2016 09:38:34 +0000",
        "state": "Success",
        "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/601.4.4 (KHTML, like Gecko) Version/9.0.3 Safari/601.4.4",
        "uuid": "2b072082-b943-4dc0-9e18-3242bd40978a"
      },
      {
        "action": "Service Create",
        "can_be_canceled": false,
        "can_be_retried": false,
        "created": "Mon, 29 Feb 2016 09:38:28 +0000",
        "end_date": "Mon, 29 Feb 2016 09:38:28 +0000",
        "ip": "172.17.0.1",
        "location": "unknown",
        "method": "POST",
        "object": "/api/app/v1/service/573e735c-8f0c-4dc2-9836-1d14ed602cb2/",
        "path": "/api/app/v1/service/",
        "resource_uri": "/api/audit/v1/action/d10b0ac9-6d07-48ed-b961-911d2d78c7e2/",
        "start_date": "Mon, 29 Feb 2016 09:38:28 +0000",
        "state": "Success",
        "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/601.4.4 (KHTML, like Gecko) Version/9.0.3 Safari/601.4.4",
        "uuid": "d10b0ac9-6d07-48ed-b961-911d2d78c7e2"
      }
    ];
    /* eslint-enable quotes, quote-props, max-len, comma-dangle */
    return (
      <div className="clearfix" style={{ margin: '40px 0' }}>
        <Timeline actions={actions} />
      </div>
    );
  }
}
