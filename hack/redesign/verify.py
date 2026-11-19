#!/usr/bin/env python3
"""Compare two built, flattened Hugo artifacts, including aliases and heading IDs."""
import json
import re
import sys
import xml.etree.ElementTree as ET
from html.parser import HTMLParser
from pathlib import Path
from urllib.parse import urlsplit

class Page(HTMLParser):
    def __init__(self, source):
        super().__init__()
        self.ids, self.headings, self.links = set(), [], []
        self.empty_links = 0
        self.feed(source)

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if attrs.get('id'):
            self.ids.add(attrs['id'])
        if tag in ('h2', 'h3', 'h4', 'h5', 'h6') and attrs.get('id'):
            self.headings.append(attrs['id'])
        if tag == 'a' and attrs.get('href') == '':
            self.empty_links += 1
        if tag == 'a' and attrs.get('href'):
            self.links.append(attrs['href'])

def pages(root):
    return {('/' + str(p.relative_to(root)).removesuffix('index.html')): p for p in root.rglob('*.html')}

def verify(before, after):
    old, current = pages(before), pages(after)
    redirects = json.loads((after / 'redirects.json').read_text())
    changes = {'/manuals/': '/c/', '/get-started/': '/learn/', '/get-started/get-docker/': '/get-docker/'}
    old_redirects = json.loads((before / 'redirects.json').read_text())
    retired = lambda route: (urlsplit(route).path.startswith(('/guides/', '/get-started/', '/reference/samples/')) or urlsplit(route).path == '/reference/glossary/')
    removed_headings = {'/testcontainers/': {'guides'}}
    alias_changes = {
        '/engine/get-started/': '/get-docker/',
        '/engine/tutorials/usingdocker/': '/get-docker/',
        '/guides/getting-started/get-docker-desktop/': '/get-docker/',
        **{f'/guides/{slug}/': None for slug in (
            'dhi-from-doi', 'dhi-from-wolfi', 'dhi-go-example',
            'dhi-nodejs-example', 'dhi-python-example', 'lab-docker-for-ai-redirect',
        )},
    }
    failures = []
    sitemap_urls = [node.text for node in ET.parse(after / 'sitemap.xml').iter()
                    if node.tag.endswith('}loc')]
    metadata_urls = [item.get('url') for item in json.loads((after / 'metadata.json').read_text())]
    old_indexed = {urlsplit(item['url']).path
                   for item in json.loads((before / 'metadata.json').read_text())
                   if item.get('url')}
    old_missing_indexed = old_indexed - old.keys()
    for name, urls in (('sitemap', sitemap_urls), ('metadata', metadata_urls)):
        if any(not url for url in urls):
            failures.append(f'Empty URL in {name}')
        for url in filter(None, urls):
            route = urlsplit(url).path
            if route not in current and route not in old_missing_indexed:
                failures.append(f'{name} points to a missing page: {route}')
            if route in redirects and old_redirects.get(route) != redirects[route]:
                failures.append(f'{name} includes a redirect: {route}')
    if set(sitemap_urls) != set(metadata_urls):
        failures.append('Sitemap and metadata inventories disagree')
    for route, path in old.items():
        if route in changes:
            if redirects.get(route) != changes[route]: failures.append(f'Wrong redirect: {route}')
            continue
        if retired(route) or retired(old_redirects.get(route, '')):
            continue
        if route not in current:
            failures.append(f'Lost route: {route}')
            continue
        original = Page(path.read_text())
        result = Page(current[route].read_text())
        if result.empty_links: failures.append(f'Empty links on {route}: {result.empty_links}')
        # Discovery page bodies intentionally change. Product heading anchors do not.
        if route not in ('/', '/reference/'):
            # Hugo embeds a build-dependent shortcode ordinal in eight existing
            # settings headings. Compare their semantic prefix and shortcode slot.
            normalize = lambda value: re.sub(r'hahahugoshortcode\d+s', 'hahahugoshortcode-s', value)
            missing = {normalize(value) for value in original.headings} - {normalize(value) for value in result.ids}
            missing -= removed_headings.get(route, set())
            if missing: failures.append(f'Lost headings on {route}: {sorted(missing)}')
    for route in current.keys() - old.keys():
        if Page(current[route].read_text()).empty_links: failures.append(f'Empty links on {route}')
    for route, target in old_redirects.items():
        expected = changes.get(target, target)
        if retired(target):
            # Deleted pages have no fallback; migrated lessons and aliases
            # must resolve directly to the same Learn destination.
            expected = redirects.get(target)
        if retired(route) and not retired(target):
            expected = changes.get(target, target)
        expected = alias_changes.get(route, expected)
        if route == expected and route in current and route not in redirects:
            continue
        if redirects.get(route) != expected:
            failures.append(f'Changed alias {route}: expected {expected}, got {redirects.get(route)}')
    for route in redirects:
        seen, target = {route}, redirects[route]
        while target in redirects:
            if target in seen:
                failures.append(f'Redirect loop: {route}')
                break
            seen.add(target)
            target = redirects[target]
        parsed = urlsplit(target)
        if not parsed.netloc and not parsed.path.startswith('/learn/') and parsed.path not in current and not (after / parsed.path.lstrip('/')).is_file():
            # Only fail added or changed redirects; existing upstream artifacts can contain historical targets.
            if old_redirects.get(route) != redirects[route]: failures.append(f'Missing redirect destination: {route} -> {target}')
    if failures:
        raise SystemExit('\n'.join(failures))
    print(f'Checked {len(old)} baseline routes, allowing retired learning, Samples, and Glossary pages; validated {len(redirects)} redirects.')

if __name__ == '__main__':
    if len(sys.argv) != 3: raise SystemExit('Usage: python3 hack/redesign/verify.py <baseline-artifact> <candidate-artifact>')
    verify(Path(sys.argv[1]), Path(sys.argv[2]))
