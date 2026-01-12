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

    // Check Accept header for markdown/text requests
    const headers = request.headers;
    const acceptHeader = headers.accept ? headers.accept[0].value : '';
    const wantsMarkdown = acceptHeader.includes('text/markdown') ||
                          acceptHeader.includes('text/plain');

    // Handle directory requests by appending index.html or index.md for requests without file extensions
    let uri = request.uri;

    // Check if the URI has a dot after the last slash (indicating a filename)
    // This is more accurate than just checking the end of the URI
    const hasFileExtension = /\.[^/]*$/.test(uri.split('/').pop());

    // If it's not a file, treat it as a directory
    if (!hasFileExtension) {
        // Ensure the URI ends with a slash before appending index file
        if (!uri.endsWith("/")) {
            uri += "/";
        }
        // Serve markdown if Accept header requests it, otherwise serve HTML
        uri += wantsMarkdown ? "index.md" : "index.html";
        request.uri = uri;
    } else if (wantsMarkdown && uri.endsWith('.html')) {
        // If requesting a specific HTML file but wants markdown, try the .md version
        uri = uri.replace(/\.html$/, '.md');
        request.uri = uri;
    }

    callback(null, request);
};
