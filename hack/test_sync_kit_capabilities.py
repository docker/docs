"""Test the release catalog import without contacting GitHub."""
import copy
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location("sync", Path(__file__).with_name("sync-kit-capabilities.py"))
sync = importlib.util.module_from_spec(spec)
spec.loader.exec_module(sync)
VERSION = "v3.0.0-m.14"
REVISION = "a" * 40


class SyncTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        (self.root / "data").mkdir()
        self.catalog_path = self.root / sync.CATALOG
        self.support_path = self.root / sync.SUPPORT
        self.catalog_path.write_bytes(b"old catalog")
        self.page = self.root / "page.md"
        self.page.write_bytes(b"unchanged page")
        self.name = "com.docker.sandbox/alpha@1"
        base = f"https://github.com/{sync.REPOSITORY}/blob/{VERSION}/"
        self.entry = {"id": self.name, "description": "A capability.",
                      "spec_url": base + "docs/spec/capabilities/" + self.name + ".md",
                      "schema_url": base + "schema/capabilities/" + self.name + ".schema.json",
                      **{k: "b" * 64 for k in ("spec_sha256", "schema_sha256", "schema_contract_sha256")}}
        self.catalog = {"format_version": 1, "schema_version": "3", "version": VERSION,
                        "revision": REVISION, "shared_spec_sha256": "c" * 64,
                        "conformance_sha256": "d" * 64, "capabilities": [self.entry]}
        self.support = {"format_version": 1, "capabilities": [{"id": self.name,
                        "local": "v0.1.0", "cloud": "Unsupported", "restrictions": "None"}]}
        self.support_path.write_text(json.dumps(self.support))
        self.support_bytes = self.support_path.read_bytes()
        self.payload = json.dumps(self.catalog, indent=2).encode() + b"\n"
        self.checksum = None
        self.calls = []
        self.release = {"tagName": VERSION, "isDraft": False, "publishedAt": "2026-01-01", "assets": [{"name": n} for n in ("capabilities.json", "checksums.txt")]}
        self.root_patch = patch.object(sync, "ROOT", self.root); self.root_patch.start(); self.addCleanup(self.root_patch.stop)
        self.run_patch = patch.object(sync, "run", side_effect=self.fake_run); self.run_patch.start(); self.addCleanup(self.run_patch.stop)
        self.report_patch = patch.object(sync, "report"); self.report_patch.start(); self.addCleanup(self.report_patch.stop)
        self.env = patch.dict(os.environ, {"KIT_SPEC_VERSION": ""}); self.env.start(); self.addCleanup(self.env.stop)

    def fake_run(self, *args, **kwargs):
        self.calls.append(args)
        if args[:3] == ("gh", "release", "list"):
            result = [{"tagName": VERSION}]
        elif args[:3] == ("gh", "release", "view"):
            result = self.release
        elif args[:2] == ("gh", "api"):
            result = {"sha": REVISION}
        elif args[:3] == ("gh", "release", "download"):
            directory = Path(args[args.index("--dir") + 1])
            (directory / "capabilities.json").write_bytes(self.payload)
            (directory / "checksums.txt").write_text(self.checksum if self.checksum is not None else hashlib.sha256(self.payload).hexdigest() + "  capabilities.json\n")
            return
        else:
            raise AssertionError(args)
        return subprocess.CompletedProcess(args, 0, json.dumps(result))

    def fail_unchanged(self):
        with self.assertRaises((ValueError, TypeError)):
            sync.main()
        self.assertEqual(self.catalog_path.read_bytes(), b"old catalog")
        self.assertEqual(self.support_path.read_bytes(), self.support_bytes)
        self.assertEqual(self.page.read_bytes(), b"unchanged page")

    def test_success_preserves_bytes_and_is_idempotent(self):
        sync.main()
        self.assertEqual(self.catalog_path.read_bytes(), self.payload)
        mtime = self.catalog_path.stat().st_mtime_ns
        sync.main()
        self.assertEqual(self.catalog_path.stat().st_mtime_ns, mtime)
        self.assertEqual(self.support_path.read_bytes(), self.support_bytes)
        self.assertEqual(self.page.read_bytes(), b"unchanged page")
        self.assertTrue(all(c[0] == "gh" for c in self.calls))
        self.assertNotIn("--exclude-pre-releases", self.calls[0])

    def test_manual_and_old_release(self):
        os.environ["KIT_SPEC_VERSION"] = VERSION
        self.release["assets"] = []
        sync.main()
        self.assertEqual(self.calls[0][:3], ("gh", "release", "view"))
        self.assertEqual(self.catalog_path.read_bytes(), b"old catalog")

    def test_invalid_tags(self):
        for tag in ("main", "v2.0.0", "v3.01.0", "v3.0.0-m.01", "v3.0.0+build"):
            with self.subTest(tag=tag):
                os.environ["KIT_SPEC_VERSION"] = tag
                self.fail_unchanged()
                self.assertEqual(self.calls, [])

    def test_checksum_failures(self):
        good = hashlib.sha256(self.payload).hexdigest() + "  capabilities.json\n"
        for value in ("", good + good, "a" * 64 + "  capabilities.json\n", "bad  capabilities.json\n", good + "bad capabilities.json extra\n"):
            with self.subTest(value=value):
                self.checksum = value
                self.fail_unchanged()

    def test_bad_json_and_metadata(self):
        for value in (b"invalid", b'{}', b'{"x":NaN}', b'{"x":1,"x":2}', json.dumps({**self.catalog, "revision": "e" * 40}).encode(), json.dumps({**self.catalog, "version": "v3.1.0"}).encode()):
            with self.subTest(value=value):
                self.payload = value
                self.fail_unchanged()

    def test_catalog_validation(self):
        cases = []
        for key, value in (("format_version", 2), ("schema_version", "2"), ("shared_spec_sha256", "bad"), ("capabilities", [])):
            item = copy.deepcopy(self.catalog); item[key] = value; cases.append(item)
        for key, value in (("id", "invalid"), ("description", ""), ("spec_url", "https://example.com"), ("schema_url", "https://example.com"), ("schema_sha256", "bad"), ("schema_contract_sha256", "bad"), ("spec_sha256", "bad")):
            item = copy.deepcopy(self.catalog); item["capabilities"][0][key] = value; cases.append(item)
        item = copy.deepcopy(self.catalog); item["capabilities"].append(copy.deepcopy(item["capabilities"][0])); cases.append(item)
        for item in cases:
            with self.subTest(item=item):
                self.payload = json.dumps(item).encode(); self.fail_unchanged()

    def test_support_orphan_and_duplicate_fail(self):
        for entries in ([{**self.support["capabilities"][0], "id": "com.docker.sandbox/missing@1"}], self.support["capabilities"] * 2):
            self.support_path.write_text(json.dumps({"format_version": 1, "capabilities": entries}))
            self.support_bytes = self.support_path.read_bytes()
            self.fail_unchanged()

    def test_invalid_support_fields_preserve_files(self):
        original = self.support["capabilities"][0]
        cases = [{**original, "extra": "x"}, {k: v for k, v in original.items() if k != "cloud"},
                 {**original, "local": "v0.1"}, {**original, "restrictions": "one\ntwo"}]
        for entry in cases:
            with self.subTest(entry=entry):
                self.support_path.write_text(json.dumps({"format_version": 1, "capabilities": [entry]}))
                self.support_bytes = self.support_path.read_bytes()
                self.fail_unchanged()

    def test_atomic_replace_failure_preserves_files_and_cleans_scratch(self):
        with patch.object(sync.os, "replace", side_effect=OSError("disk failure")):
            with self.assertRaises(OSError):
                sync.main()
        self.assertEqual(self.catalog_path.read_bytes(), b"old catalog")
        self.assertEqual(self.support_path.read_bytes(), self.support_bytes)
        self.assertEqual(self.page.read_bytes(), b"unchanged page")
        self.assertEqual(list(self.catalog_path.parent.glob(".kit-catalog-*")), [])

    def test_latest_list_flags_and_duplicate_assets(self):
        self.release["assets"].append({"name": "checksums.txt"})
        self.fail_unchanged()
        self.assertIn("--exclude-drafts", self.calls[0])
        self.assertEqual(self.calls[0][self.calls[0].index("--limit") + 1], "1")
        self.release["assets"] = [{"name": "capabilities.json"}, {"name": "capabilities.json"}, {"name": "checksums.txt"}]
        self.fail_unchanged()

    def test_unknown_catalog_entry_is_allowed(self):
        self.support_path.write_text(json.dumps({"format_version": 1, "capabilities": []}))
        sync.main()
        self.assertEqual(self.catalog_path.read_bytes(), self.payload)

    def test_unpublished_and_bad_assets_fail(self):
        for change in ({"isDraft": True}, {"publishedAt": None}, {"tagName": "v3.0.0"}, {"assets": [{"name": "capabilities.json"}]}):
            with self.subTest(change=change), patch.dict(self.release, change):
                self.fail_unchanged()


if __name__ == "__main__":
    unittest.main()
