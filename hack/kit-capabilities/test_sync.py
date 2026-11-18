"""Failure cases for the pinned capability catalog and support metadata."""

import copy
import importlib.util
import json
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch


MODULE_PATH = Path(__file__).with_name("sync.py")
SPEC = importlib.util.spec_from_file_location("kit_capability_sync", MODULE_PATH)
sync = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(sync)


class CatalogTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.schema_bytes = sync.SCHEMA.read_bytes()
        cls.source = sync.read_json(sync.SOURCE)
        cls.support = sync.read_json(sync.SUPPORT)

    def test_generated_data_preserves_classified_support(self):
        data = sync.make_data(self.schema_bytes, self.source, self.support)
        self.assertEqual(json.loads(sync.OUTPUT.read_bytes()), data)
        by_id = {entry["id"]: entry for entry in data["capabilities"]}
        self.assertEqual(set(by_id), set(self.support))
        self.assertEqual(by_id["com.docker.sandbox/agent-context@2"]["cloud"],
                         {"status": "limited", "since": "v0.51.0"})
        self.assertEqual(by_id["com.docker.sandbox/agent-sessions@1"]["local"],
                         {"status": "declaration-only", "since": "v0.45.0"})
        self.assertEqual(by_id["com.docker.sandbox/host-port@1"]["local"],
                         {"status": "unverified"})

    def test_new_spec_type_requires_explicit_classification(self):
        schema = json.loads(self.schema_bytes)
        branch = copy.deepcopy(schema["definitions"]["capability"]["allOf"][0])
        identifier = "com.docker.sandbox/future@1"
        branch["if"]["properties"]["type"]["const"] = identifier
        branch["then"]["properties"]["config"]["$ref"] = f"capabilities/{identifier}.schema.json"
        schema["definitions"]["capability"]["allOf"].append(branch)
        source = copy.deepcopy(self.source)
        source["documentation"].append(f"docs/spec/capabilities/{identifier}.md")
        source["documentation"].sort()
        source["schema_sha256"] = sync.hashlib.sha256(sync.json_bytes(schema)).hexdigest()
        with self.assertRaisesRegex(sync.CatalogError, "future@1"):
            sync.make_data(sync.json_bytes(schema), source, self.support)
        support = copy.deepcopy(self.support)
        support[identifier] = {
            "local": {"status": "unverified"},
            "cloud": {"status": "unverified"},
            "restrictions": "Docker Sandboxes support has not been verified.",
        }
        data = sync.make_data(sync.json_bytes(schema), source, support)
        self.assertIn(identifier, {entry["id"] for entry in data["capabilities"]})

    def test_rejects_missing_docs_and_duplicate_type(self):
        docs = self.source["documentation"][1:]
        with self.assertRaisesRegex(sync.CatalogError, "spec documentation mismatch"):
            sync.inventory(self.schema_bytes, docs)
        schema = json.loads(self.schema_bytes)
        schema["definitions"]["capability"]["allOf"].append(
            copy.deepcopy(schema["definitions"]["capability"]["allOf"][0]))
        with self.assertRaisesRegex(sync.CatalogError, "duplicate capability ID"):
            sync.inventory(sync.json_bytes(schema), self.source["documentation"])

    def test_rejects_integrity_and_invalid_status_metadata(self):
        with self.assertRaisesRegex(sync.CatalogError, "checksum mismatch"):
            sync.make_data(self.schema_bytes + b" ", self.source, self.support)
        support = copy.deepcopy(self.support)
        support["com.docker.sandbox/port@1"]["cloud"]["status"] = "maybe"
        with self.assertRaisesRegex(sync.CatalogError, "invalid status"):
            sync.make_data(self.schema_bytes, self.source, support)
        support = copy.deepcopy(self.support)
        support["com.docker.sandbox/port@1"]["cloud"]["since"] = "v3.0"
        with self.assertRaisesRegex(sync.CatalogError, "invalid sbx version"):
            sync.make_data(self.schema_bytes, self.source, support)
        support = copy.deepcopy(self.support)
        support["com.docker.sandbox/port@1"]["restrictions"] = "first | second"
        with self.assertRaisesRegex(sync.CatalogError, "one Markdown table cell"):
            sync.make_data(self.schema_bytes, self.source, support)

    def test_check_rejects_stale_output(self):
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory) / "kit_capabilities.json"
            output.write_text("{}\n")
            with patch.object(sync, "OUTPUT", output), patch.object(sync.sys, "argv", ["sync.py", "check"]):
                with self.assertRaisesRegex(SystemExit, "1"):
                    sync.main()


if __name__ == "__main__":
    unittest.main()
