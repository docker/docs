import cf from 'cloudfront';

const kvs = cf.kvs();

const redirectsPrefixes = {{.RedirectsPrefixesJSON}};

async function handler(event) {
    const request = event.request;
    const lookupKey = request.uri.replace(/\/$/, '');

    if (lookupKey !== '') {
        try {
            const target = await kvs.get(lookupKey);
            return {
                statusCode: 301,
                statusDescription: 'Moved Permanently',
                headers: { location: { value: target } },
            };
        } catch (err) {
            // not found in KVS — fall through
        }
    }

    for (let x in redirectsPrefixes) {
        const rp = redirectsPrefixes[x];
        if (!request.uri.startsWith(`/${rp['prefix']}`)) {
            continue;
        }
        let newlocation = '/';
        if (rp['strip']) {
            const re = new RegExp(`(^/${rp['prefix']})`, 'gi');
            newlocation = request.uri.replace(re, '/');
        }
        return {
            statusCode: 301,
            statusDescription: 'Moved Permanently',
            headers: { location: { value: newlocation } },
        };
    }

    const headers = request.headers;
    const acceptHeader = headers.accept ? headers.accept.value : '';
    const wantsMarkdown = acceptHeader.includes('text/markdown') ||
                         acceptHeader.includes('text/plain');

    let uri = request.uri;
    const hasFileExtension = /\.[^/]*$/.test(uri.split('/').pop());

    if (!hasFileExtension) {
        if (wantsMarkdown) {
            const stripped = uri.replace(/\/$/, '');
            uri = stripped === '' ? '/index.md' : stripped + '.md';
        } else {
            if (!uri.endsWith('/')) {
                uri += '/';
            }
            uri += 'index.html';
        }
        request.uri = uri;
    } else if (wantsMarkdown && uri.endsWith('/index.html')) {
        uri = uri === '/index.html' ? '/index.md' : uri.replace(/\/index\.html$/, '.md');
        request.uri = uri;
    }

    return request;
}
