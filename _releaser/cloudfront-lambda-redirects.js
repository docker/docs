'use strict';

exports.handler = (event, context, callback) => {
    //console.log("event", JSON.stringify(event));
    const request = event.Records[0].cf.request;

    const redirects = JSON.parse(`{{.RedirectsJSON}}`);
    for (let key in redirects) {
        if (key !== request.uri) {
            continue;
        }
        //console.log(`redirect: ${request.uri} to ${redirects[key]}`);
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

    const redirectsPrefixes = JSON.parse(`{{.RedirectsPrefixesJSON}}`);
    for (let x in redirectsPrefixes) {
        const rp = redirectsPrefixes[x];
        if (!request.uri.startsWith(`/${rp['prefix']}`)) {
            continue;
        }
        let newlocation = "/";
        if (rp['strip']) {
            let re = new RegExp(`(^/${rp['prefix']})`, 'gi');
            newlocation = request.uri.replace(re,'/');
        }
        //console.log(`redirect: ${request.uri} to ${redirectsPrefixes[key]}`);
        const response = {
            status: '301',
            statusDescription: 'Moved Permanently',
            headers: {
                location: [{
                    key: 'Location',
                    value: newlocation,
                }],
            },
        }
        callback(null, response);
        return
    }

    callback(null, request);
};
