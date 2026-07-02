'use strict';

// Tests for cloudfront-lambda-redirects.js.
//
// The function source is a Go text/template (see aws.go) with two
// placeholders that are filled with JSON at build time. These tests render
// the template the same way aws.go does, then load and exercise the handler.
//
// Run with: node --test hack/releaser/cloudfront-lambda-redirects.test.js

const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const Module = require('node:module');

const SRC = path.join(__dirname, 'cloudfront-lambda-redirects.js');

const REDIRECTS = {
    '/old/page/': '/new/page/',
    '/target-with-query/': '/dest/?ref=docs',
};

const REDIRECTS_PREFIXES = [
    { prefix: 'keep/', strip: false },
    { prefix: 'strip/', strip: true },
];

// Render the template the same way getLambdaFunctionZip in aws.go does, then
// evaluate it as a CommonJS module so we can call handler() directly.
function loadHandler() {
    const rendered = fs
        .readFileSync(SRC, 'utf8')
        // Function replacers avoid String.replace's special handling of `$`
        // sequences in the replacement, in case a redirect target contains one.
        .replace('{{.RedirectsJSON}}', () => JSON.stringify(REDIRECTS))
        .replace('{{.RedirectsPrefixesJSON}}', () => JSON.stringify(REDIRECTS_PREFIXES));

    const m = new Module(SRC);
    m._compile(rendered, SRC);
    return m.exports.handler;
}

const handler = loadHandler();

// invoke wraps the callback-style handler in a promise.
function invoke({ uri, querystring = '', accept = 'text/html' }) {
    const request = {
        uri,
        querystring,
        headers: { accept: [{ key: 'Accept', value: accept }] },
    };
    const event = { Records: [{ cf: { request } }] };
    return new Promise((resolve, reject) => {
        handler(event, {}, (err, result) => {
            if (err) return reject(err);
            resolve({ result, request });
        });
    });
}

function locationOf(result) {
    return result.headers.location[0].value;
}

test('exact redirect preserves the query string', async () => {
    const { result } = await invoke({
        uri: '/old/page/',
        querystring: 'utm_source=newsletter&utm_campaign=launch',
    });
    assert.equal(result.status, '301');
    assert.equal(locationOf(result), '/new/page/?utm_source=newsletter&utm_campaign=launch');
});

test('exact redirect without a query string is unchanged', async () => {
    const { result } = await invoke({ uri: '/old/page/' });
    assert.equal(result.status, '301');
    assert.equal(locationOf(result), '/new/page/');
});

test('exact redirect appends with & when target already has a query string', async () => {
    const { result } = await invoke({
        uri: '/target-with-query/',
        querystring: 'utm_medium=email',
    });
    assert.equal(locationOf(result), '/dest/?ref=docs&utm_medium=email');
});

test('prefix redirect (strip) preserves the query string', async () => {
    const { result } = await invoke({
        uri: '/strip/some/path',
        querystring: 'utm_source=x',
    });
    assert.equal(result.status, '301');
    assert.equal(locationOf(result), '/some/path?utm_source=x');
});

test('prefix redirect (no strip) preserves the query string', async () => {
    const { result } = await invoke({
        uri: '/keep/anything',
        querystring: 'utm_source=x',
    });
    assert.equal(result.status, '301');
    assert.equal(locationOf(result), '/?utm_source=x');
});

test('directory rewrite passes the request through with query string intact', async () => {
    const { result, request } = await invoke({
        uri: '/some/page',
        querystring: 'utm_source=x',
    });
    // No redirect response: the handler forwards the (mutated) request.
    assert.equal(result, request);
    assert.equal(request.uri, '/some/page/index.html');
    assert.equal(request.querystring, 'utm_source=x');
});
