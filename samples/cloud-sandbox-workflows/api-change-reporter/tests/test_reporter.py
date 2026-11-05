import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

from reporter import compare, load_spec, render_report

ROOT = Path(__file__).resolve().parents[1]


class ReporterTests(unittest.TestCase):
    def test_example_removals(self):
        result = compare(
            load_spec(ROOT / "examples/before.yaml"),
            load_spec(ROOT / "examples/after.json"),
        )
        self.assertEqual(
            result,
            {
                "removed_paths": ["/legacy"],
                "removed_operations": {"/legacy": ["GET"], "/pets/{id}": ["DELETE"]},
            },
        )

    def test_unchanged_and_added_operations(self):
        self.assertEqual(
            compare(
                {"paths": {"/pets": {"get": {}}}},
                {"paths": {"/pets": {"get": {}, "post": {}}, "/added": {"get": {}}}},
            ),
            {"removed_paths": [], "removed_operations": {}},
        )

    def test_groups_methods_and_ignores_metadata(self):
        result = compare(
            {
                "paths": {
                    "/pets": {
                        "get": {},
                        "delete": {},
                        "summary": "Pets",
                        "parameters": [],
                    }
                }
            },
            {"paths": {}},
        )
        self.assertEqual(result["removed_operations"], {"/pets": ["DELETE", "GET"]})

    def test_json_and_yaml_load_equally(self):
        expected = load_spec(ROOT / "examples/before.yaml")
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "before.json"
            path.write_text(json.dumps(expected))
            self.assertEqual(load_spec(path), expected)

    def test_invalid_or_referenced_paths(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "bad.yaml"
            for value in [
                "[]",
                "paths: []",
                "paths: {pets: {}}",
                "paths: {/pets: null}",
                "paths: {/pets: {$ref: 'other.yaml'}}",
            ]:
                with self.subTest(value=value):
                    path.write_text(value)
                    with self.assertRaises(ValueError):
                        load_spec(path)

    def test_html_escapes_spec_content(self):
        html = render_report(
            {
                "removed_paths": ["/<script>"],
                "removed_operations": {"/<script>": ["GET"]},
            }
        )
        self.assertNotIn("<script>", html)
        self.assertIn("/&lt;script&gt;", html)

    def test_cli_writes_report(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "output/index.html"
            run = subprocess.run(
                [
                    sys.executable,
                    str(ROOT / "reporter.py"),
                    str(ROOT / "examples/before.yaml"),
                    str(ROOT / "examples/after.json"),
                    "--html",
                    str(path),
                ],
                capture_output=True,
                text=True,
            )
            self.assertEqual(run.returncode, 0, run.stderr)
            self.assertIn("Removed operations: 2", run.stdout)
            self.assertIn("/pets/{id}", path.read_text())

    def test_cli_rejects_missing_input(self):
        with tempfile.TemporaryDirectory() as directory:
            run = subprocess.run(
                [
                    sys.executable,
                    str(ROOT / "reporter.py"),
                    str(Path(directory) / "missing.yaml"),
                    str(ROOT / "examples/after.json"),
                    "--html",
                    str(Path(directory) / "index.html"),
                ],
                capture_output=True,
                text=True,
            )
            self.assertEqual(run.returncode, 2)
            self.assertIn("error:", run.stderr)


if __name__ == "__main__":
    unittest.main()
