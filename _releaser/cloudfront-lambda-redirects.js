'use strict';

exports.handler = (event, context, callback) => {
    //console.log("event", JSON.stringify(event));
    const request = event.Records[0].cf.request;
    const redirects = JSON.parse(`{{.RedirectsJSON}}`);
    for (let key in redirects) {
        if (key !== request.uri) {
            continue;
        }
        //console.log(`redirect: ${key} to ${redirects[key]}`);
        const response = {
            status: '301',
            statusDescription: 'Moved Permanently',
            headers: {
                location: [{
                    key: 'Location',
                    value: redirects[key],
                }],
            },
        }
        callback(null, response);
        return
    }
    callback(null, request);
};
