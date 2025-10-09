'use strict';

exports.handler = (event, context, callback) => {
    //console.log("event", JSON.stringify(event));
    const request = event.Records[0].cf.request;
    const requestUrl = request.uri.replace(/\/$/, "")

    const redirects = JSON.parse(`{{.RedirectsJSON}}`);
    for (let key in redirects) {
        const redirectTarget = key.replace(/\/$/, "")
        if (redirectTarget !== requestUrl) {
            continue;
        }
        //console.log(`redirect: ${requestUrl} to ${redirects[key]}`);
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

    // Handle directory requests by appending index.html for requests without file extensions
    let uri = request.uri;

    // Check if the URI has a dot after the last slash (indicating a filename)
    // This is more accurate than just checking the end of the URI
    const hasFileExtension = /\.[^/]*$/.test(uri.split('/').pop());

    // If it's not a file, treat it as a directory and append index.html
    if (!hasFileExtension) {
        // Ensure the URI ends with a slash before appending index.html
        if (!uri.endsWith("/")) {
            uri += "/";
        }
        uri += "index.html";
        request.uri = uri;
    }

    callback(null, request);
};
