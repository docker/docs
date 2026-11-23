"""Run with python3 -m unittest discover -s hack -p 'test_sync_kit_capabilities.py'."""

import hashlib
import importlib.util
import io
import json
import os
from pathlib import Path
import subprocess
import tarfile
import tempfile
import unittest
from unittest.mock import patch


spec = importlib.util.spec_from_file_location(
    "sync_kit_capabilities", Path(__file__).with_name("sync-kit-capabilities.py")
)
sync = importlib.util.module_from_spec(spec)
spec.loader.exec_module(sync)
VERSION = "v3.0.0-m.14"
REVISION = "a" * 40
GO_MOD = f"module {sync.MODULE}\n\ngo 1.23\nreplace golang.org/x/mod => golang.org/x/mod v0.43.0\n"


def source_archive(extra=()):
    output = io.BytesIO()
    with tarfile.open(fileobj=output, mode="w:gz") as archive:
        for name, body in (
            ("source/go.mod", GO_MOD),
            ("source/cmd/capability-catalog/main.go", "package main\nfunc main() {}\n"),
        ):
            data = body.encode()
            member = tarfile.TarInfo(name)
            member.size = len(data)
            archive.addfile(member, io.BytesIO(data))
        for member in extra:
            archive.addfile(member)
    return output.getvalue()


class SyncTests(unittest.TestCase):
    def setUp(self):
        self.environment = patch.dict(os.environ, {
            "KIT_SPEC_VERSION": "", "GH_TOKEN": "secret", "GITHUB_TOKEN": "secret",
        })
        self.environment.start()
        self.addCleanup(self.environment.stop)
        self.release = {
            "tagName": VERSION, "isDraft": False, "publishedAt": "2026-11-23T00:00:00Z",
            "assets": [{"name": name} for name in ("capabilities.json", "checksums.txt")],
        }
        self.calls = []
        self.payload = json.dumps({"version": VERSION, "revision": REVISION}).encode()
        self.checksum = None
        self.updater_fails = False
        self.archive = source_archive()
        self.updater_module = None
        self.command = patch.object(sync, "run", side_effect=self.fake_run)
        self.command.start()
        self.addCleanup(self.command.stop)
        self.report = patch.object(sync, "report").start()
        self.addCleanup(patch.stopall)

    def fake_run(self, *args, **kwargs):
        self.calls.append((args, kwargs))
        if args[:3] == ("gh", "release", "list"):
            result = [{"tagName": VERSION}]
        elif args[:3] == ("gh", "release", "view"):
            result = self.release
        elif args[:2] == ("gh", "api") and "/tarball/" in args[2]:
            kwargs["stdout"].write(self.archive)
            return
        elif args[:2] == ("gh", "api"):
            result = {"sha": REVISION}
        elif args[:3] == ("gh", "release", "download"):
            directory = Path(args[args.index("--dir") + 1])
            (directory / "capabilities.json").write_bytes(self.payload)
            checksum = self.checksum
            if checksum is None:
                checksum = hashlib.sha256(self.payload).hexdigest() + "  capabilities.json\n"
            (directory / "checksums.txt").write_text(checksum)
            return
        elif args[:2] == ("go", "run"):
            self.updater_module = (Path(kwargs["cwd"]) / "go.mod").read_text()
            if self.updater_fails:
                raise subprocess.CalledProcessError(1, args)
            return
        else:
            raise AssertionError(args)
        return subprocess.CompletedProcess(args, 0, json.dumps(result))

    def assert_no_updater(self):
        self.assertFalse(any(args[0] == "go" for args, _ in self.calls))

    def test_latest_includes_prereleases_and_pins_updater_without_tokens(self):
        sync.main()
        listing = self.calls[0][0]
        self.assertIn("--exclude-drafts", listing)
        self.assertNotIn("--exclude-pre-releases", listing)
        self.assertEqual(listing[listing.index("--limit") + 1], "1")
        args, kwargs = self.calls[-1]
        self.assertEqual(args[:4], ("go", "run", "./cmd/capability-catalog", "docs"))
        self.assertEqual(Path(args[-1]), sync.ROOT / sync.PAGE)
        self.assertEqual(self.updater_module, GO_MOD)
        self.assertTrue(any(
            args == ("gh", "api", f"repos/{sync.REPOSITORY}/tarball/{REVISION}")
            for args, _ in self.calls
        ))
        self.assertNotIn("GH_TOKEN", kwargs["env"])
        self.assertNotIn("GITHUB_TOKEN", kwargs["env"])
        self.assertEqual(kwargs["env"]["GOWORK"], "off")
        self.assertNotEqual(Path(kwargs["cwd"]), sync.ROOT)
        self.assertFalse(Path(kwargs["cwd"]).exists())  # Temporary assets cleaned up.
        self.assertEqual(self.report.call_args.kwargs, {"available": "true", "version": VERSION})

    def test_explicit_version_resolves_published_release_without_listing(self):
        os.environ["KIT_SPEC_VERSION"] = VERSION
        sync.main()
        self.assertEqual(self.calls[0][0][:4], ("gh", "release", "view", VERSION))

    def test_old_release_skips_cleanly(self):
        self.release["assets"] = []
        sync.main()
        self.assert_no_updater()
        self.assertEqual(len(self.calls), 2)
        self.assertEqual(self.report.call_args.kwargs, {"available": "false"})

    def test_invalid_tags_never_reach_upstream(self):
        for tag in ("main", "v2.0.0", "v3.01.0", "v3.0.0-m.01", "v3.0.0+build", "v3.0.0\n", "--help", "v3.0.0;echo bad"):
            with self.subTest(tag=tag):
                os.environ["KIT_SPEC_VERSION"] = tag
                with self.assertRaises(ValueError):
                    sync.main()
                self.assertEqual(self.calls, [])

    def test_draft_unpublished_or_mismatched_release_fails(self):
        for field, value in (("isDraft", True), ("publishedAt", None), ("tagName", "v3.0.0")):
            with self.subTest(field=field), patch.dict(self.release, {field: value}):
                with self.assertRaises(ValueError):
                    sync.main()
                self.assert_no_updater()

    def test_missing_or_duplicate_assets_fail(self):
        for names in (("capabilities.json",), ("capabilities.json", "checksums.txt", "checksums.txt"), ("capabilities.json", "capabilities.json", "checksums.txt")):
            with self.subTest(names=names):
                self.release["assets"] = [{"name": name} for name in names]
                with self.assertRaises(ValueError):
                    sync.main()
                self.assert_no_updater()

    def test_bad_missing_duplicate_or_malformed_checksum_fails(self):
        correct = hashlib.sha256(self.payload).hexdigest() + "  capabilities.json\n"
        for checksum in ("", correct + correct, correct + "bad capabilities.json extra\n", "a" * 64 + "  capabilities.json\n", "bad  capabilities.json\n"):
            with self.subTest(checksum=checksum):
                self.checksum = checksum
                with self.assertRaises(ValueError):
                    sync.main()
                self.assert_no_updater()

    def test_invalid_json_or_catalog_metadata_fails(self):
        for payload in (b"not JSON", b"[]", b"{}", b'{"invalid": NaN}', b'{"version": 1, "version": 2}', json.dumps({"version": VERSION, "revision": "b" * 40}).encode(), json.dumps({"version": "v3.0.0", "revision": REVISION}).encode()):
            with self.subTest(payload=payload):
                self.payload = payload
                with self.assertRaises(ValueError):
                    sync.main()
                self.assert_no_updater()

    def test_unsafe_or_ambiguous_source_archive_never_executes(self):
        entries = []
        for name in ("../outside", "/absolute", "source/../../outside", "other/go.mod", "source/go.mod"):
            entries.append(tarfile.TarInfo(name))
        for kind in (tarfile.SYMTYPE, tarfile.LNKTYPE, tarfile.FIFOTYPE):
            member = tarfile.TarInfo("source/link")
            member.type = kind
            member.linkname = "/outside"
            entries.append(member)
        for member in entries:
            with self.subTest(name=member.name, kind=member.type):
                self.archive = source_archive([member])
                with self.assertRaises(ValueError):
                    sync.main()
                self.assert_no_updater()

    def test_missing_or_wrong_source_module_is_rejected(self):
        with tempfile.TemporaryDirectory() as temporary:
            directory = Path(temporary)
            archive = directory / "source.tar.gz"
            for files in ({"root/other": ""}, {"root/go.mod": "module wrong", "root/cmd/capability-catalog/main.go": "package main"}):
                with self.subTest(files=files):
                    with tarfile.open(archive, "w:gz") as tar:
                        for name, body in files.items():
                            data = body.encode()
                            member = tarfile.TarInfo(name)
                            member.size = len(data)
                            tar.addfile(member, io.BytesIO(data))
                    with self.assertRaises(ValueError):
                        sync.extract_source(archive, directory / "extract")

    def test_updater_failure_propagates_and_cleans_assets(self):
        self.updater_fails = True
        with self.assertRaises(subprocess.CalledProcessError):
            sync.main()
        self.assertFalse(Path(self.calls[-1][1]["cwd"]).exists())
        self.report.assert_not_called()


if __name__ == "__main__":
    unittest.main()
