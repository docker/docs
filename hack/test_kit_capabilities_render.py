"""Render the real capability reference in a small HTML/Markdown Hugo site.

Run with HUGO=/path/to/hugo python3 -m unittest discover -s hack
-p test_kit_capabilities_render.py. Set KIT_RENDER_OUTPUT to retain a fixture.
Support data uses format_version=1 and a capabilities array of unique
{id, local, cloud, restrictions} records. IDs join the released catalog.
"""

import copy
from html.parser import HTMLParser
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import tempfile
import unittest


ROOT = Path(__file__).resolve().parent.parent
PAGE = Path("content/manuals/ai/sandboxes/customize/author/capabilities.md")
CATALOG = "sandbox_kit_capabilities"
SUPPORT = "sandbox_kit_support"
HUGO = os.environ.get("HUGO", "hugo")


def unique_keys(pairs):
    result = {}
    for key, value in pairs:
        if key in result:
            raise ValueError(f"Duplicate JSON key: {key}")
        result[key] = value
    return result


def strict_json(text):
    def invalid(value):
        raise ValueError(f"Invalid JSON constant: {value}")
    return json.loads(text, object_pairs_hook=unique_keys, parse_constant=invalid)


class Text(HTMLParser):
    def __init__(self):
        super().__init__()
        self.parts = []

    def handle_data(self, data):
        self.parts.append(data)


def plain_html(source):
    parser = Text()
    parser.feed(source)
    return "".join(parser.parts)


def plain_markdown(source):
    return re.sub(r"\[([^\]]+)\]\([^)]+\)", r"\1", source).replace("`", "")


class RenderTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        # Hugo's data decoder accepts duplicate object keys. Enforce strict JSON
        # here before any rendering, including for both actual vendored files.
        cls.catalog = strict_json((ROOT / "data" / f"{CATALOG}.json").read_text())
        cls.support = strict_json((ROOT / "data" / f"{SUPPORT}.json").read_text())

    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory(prefix="kit-capabilities-render-")
        self.addCleanup(self.temporary.cleanup)
        self.site = Path(self.temporary.name)
        for directory in ("data", "layouts/_markup", "layouts/shortcodes", str(PAGE.parent)):
            (self.site / directory).mkdir(parents=True, exist_ok=True)
        for parent in PAGE.parent.parents:
            if parent == Path("content"):
                break
            (self.site / parent / "_index.md").write_text("---\ntitle: Section\n---\n")
        (self.site / PAGE.parent / "_index.md").write_text("---\ntitle: Author kits\n---\n")
        shutil.copyfile(ROOT / PAGE, self.site / PAGE)
        for layout in ("single.markdown.md", "_markup/render-link.html"):
            shutil.copyfile(ROOT / "layouts" / layout, self.site / "layouts" / layout)
        (self.site / "layouts/single.html").write_text("{{ .Content }}")
        (self.site / "layouts/shortcodes/summary-bar.html").write_text("")
        # Link targets needed by the real render-link hook. Other site content
        # and navigation are unrelated to this isolated reference rendering.
        (self.site / PAGE.parent.parent / "use-kits.md").write_text(
            '---\ntitle: Use kits\n---\n'
            '## Choose directories for kit host mounts\n\n'
            '## Keep kit data when recreating a sandbox\n'
        )
        config = {
            "baseURL": "https://example.test/", "enableInlineShortcodes": True,
            "disableKinds": ["home", "section", "taxonomy", "term", "RSS", "sitemap"],
            "outputs": {"page": ["HTML", "markdown"]},
            "outputFormats": {"markdown": {"mediaType": "text/markdown", "baseName": "index", "isPlainText": True}},
            "permalinks": {"page": {"manuals": "/:sections[1:]/:slugorcontentbasename/"}},
            "markup": {"goldmark": {"extensions": {"typographer": False}, "parser": {"attribute": {"block": True}}}},
        }
        (self.site / "hugo.json").write_text(json.dumps(config))
        for name in (CATALOG, SUPPORT):
            shutil.copyfile(ROOT / "data" / f"{name}.json", self.site / "data" / f"{name}.json")

    def write_data(self, catalog, support):
        for name, data in ((CATALOG, catalog), (SUPPORT, support)):
            (self.site / "data" / f"{name}.json").write_text(json.dumps(data))

    def render(self, success=True):
        result = subprocess.run([HUGO, "--source", str(self.site)], capture_output=True, text=True)
        output = result.stdout + result.stderr
        if not success:
            self.assertNotEqual(result.returncode, 0, output)
            self.assertIn("kit-capabilities:", output)
            return
        self.assertEqual(result.returncode, 0, output)
        directory = self.site / "public/ai/sandboxes/customize/author/capabilities"
        return (directory / "index.html").read_text(), (directory / "index.md").read_text()

    def test_real_page_renders_each_capability_once_in_both_formats(self):
        html, markdown = self.render()
        entries = {entry["id"]: entry for entry in self.support["capabilities"]}
        self.assertEqual(len(re.findall(r"^### ", markdown, re.M)), len(self.catalog["capabilities"]))
        for capability in self.catalog["capabilities"]:
            short = capability["id"].removeprefix("com.docker.sandbox/")
            anchor = short.replace("@", "")
            heading = f'<h3 id="{anchor}">{short}</h3>'
            self.assertEqual(html.count(heading), 1)
            self.assertEqual(markdown.count(f"### {short}\n"), 1)
            html_section = html.split(heading, 1)[1].split("<h3", 1)[0].split('<h2', 1)[0]
            markdown_section = markdown.split(f"### {short}\n", 1)[1].split("\n### ", 1)[0].split("\n## ", 1)[0]
            self.assertIn(capability["description"], markdown_section)
            self.assertIn(plain_markdown(capability["description"]), plain_html(html_section))
            entry = entries.get(capability["id"], {"local": "Not verified", "cloud": "Not verified", "restrictions": "Docker Sandboxes support has not been verified."})
            self.assertIn(f'| {entry["local"]} | {entry["cloud"]} |', markdown_section)
            for value in (entry["local"], entry["cloud"], plain_markdown(entry["restrictions"])):
                self.assertIn(value, plain_html(html_section))
            self.assertIn(entry["restrictions"], markdown_section)
            cells = re.findall(r"<td[^>]*>(.*?)</td>", html_section, re.S)
            self.assertEqual([plain_html(cell) for cell in cells], [entry["local"], entry["cloud"]])
            self.assertIn(f'Specification reference: [{short}]({capability["spec_url"]}).', markdown_section)
            self.assertIn(f'href="{capability["spec_url"]}"', html_section)
        for output in (html, markdown):
            self.assertNotIn("{{", output)
            self.assertNotIn("capability-catalog:descriptions", output)
            self.assertEqual(output.count(f'Docker Sandbox Kit Specification {self.catalog["version"]}'), 1)
        self.assertIn('href="https://example.test/ai/sandboxes/customize/use-kits/#choose-directories-for-kit-host-mounts"', html)
        self.assertIn('href="#port-address-families"', html)
        self.assertIn('../use-kits.md#keep-kit-data-when-recreating-a-sandbox', markdown)
        for capability in self.catalog["capabilities"]:
            for url in re.findall(r"\]\(([^)]+)\)", capability["description"]):
                self.assertTrue(url.startswith("https://"), url)
                self.assertIn(f'href="{url}"', html)
        if destination := os.environ.get("KIT_RENDER_OUTPUT"):
            shutil.copytree(self.site, destination, dirs_exist_ok=True)

    def test_catalog_only_update_refreshes_provenance_and_pinned_links(self):
        catalog = copy.deepcopy(self.catalog)
        old_version = catalog["version"]
        catalog["version"] = "v3.0.0-fixture.1"
        for entry in catalog["capabilities"]:
            for field in ("spec_url", "schema_url", "description"):
                entry[field] = entry[field].replace(old_version, catalog["version"])
        self.write_data(catalog, self.support)
        html, markdown = self.render()
        for output in (html, markdown):
            self.assertIn("Docker Sandbox Kit Specification v3.0.0-fixture.1", output)
            self.assertNotIn(f"/blob/{old_version}/", output)
            self.assertNotIn(f"/releases/tag/{old_version}", output)

    def test_unknown_catalog_capability_is_unverified(self):
        catalog = copy.deepcopy(self.catalog)
        original = catalog["capabilities"][-1]
        added = dict(original, id="com.docker.sandbox/zz-fixture@1", description="A fixture capability.")
        for field in ("schema_url", "spec_url"):
            added[field] = added[field].replace(original["id"], added["id"])
        catalog["capabilities"].append(added)
        self.write_data(catalog, self.support)
        html, markdown = self.render()
        self.assertIn("### zz-fixture@1\n", markdown)
        self.assertIn('| Not verified | Not verified |', markdown.split('### zz-fixture@1', 1)[1])
        self.assertIn('Docker Sandboxes support has not been verified.', html.split('id="zz-fixture1"', 1)[1])

    def test_orphaned_or_duplicate_support_requires_review(self):
        for change in ("orphan", "duplicate"):
            with self.subTest(change=change):
                support = copy.deepcopy(self.support)
                record = dict(support["capabilities"][0])
                if change == "orphan":
                    record["id"] = "com.docker.sandbox/removed@1"
                support["capabilities"].append(record)
                self.write_data(self.catalog, support)
                self.render(success=False)

    def test_malformed_support_is_rejected(self):
        for field, value in (("local", "Supported"), ("cloud", 12), ("restrictions", ""), ("restrictions", "two\nparagraphs"), ("unexpected", "field")):
            with self.subTest(field=field, value=value):
                support = copy.deepcopy(self.support)
                support["capabilities"][0][field] = value
                self.write_data(self.catalog, support)
                self.render(success=False)

    def test_duplicate_or_malformed_catalog_is_rejected(self):
        for change in ("duplicate", "id", "description", "spec_url", "schema_url", "format", "schema"):
            with self.subTest(change=change):
                catalog = copy.deepcopy(self.catalog)
                if change == "duplicate":
                    catalog["capabilities"].append(catalog["capabilities"][0])
                elif change == "format":
                    catalog["format_version"] = 2
                elif change == "schema":
                    catalog["schema_version"] = "4"
                else:
                    catalog["capabilities"][0][change] = ""
                self.write_data(catalog, self.support)
                self.render(success=False)

    def test_strict_json_rejects_duplicate_keys_and_non_json_constants(self):
        for text in ('{"format_version":1,"format_version":2}', '{"capabilities":NaN}'):
            with self.subTest(text=text), self.assertRaises(ValueError):
                strict_json(text)


if __name__ == "__main__":
    unittest.main()
