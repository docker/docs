// Regenerates the permission tables in the custom-role permissions reference
// from the IAM source of truth (`custom-role-metadata.yaml` in
// docker/iam-services). Run by the "Sync custom role permissions" workflow when
// that file changes; the result is opened as a PR for review.
//
// Only the region between the generation markers is rewritten — the frontmatter,
// intro, and "Next steps" are hand-maintained. Table columns are left
// unpadded here; prettier aligns them afterwards.
//
// Usage:
//   node scripts/generate-custom-role-permissions.mjs \
//     --metadata path/to/custom-role-metadata.yaml \
//     --doc content/manuals/security/roles-and-permissions/custom-roles/permissions-reference.md
//   add --check to exit non-zero (without writing) when the doc is out of date.

import { readFileSync, writeFileSync } from 'node:fs';

const BEGIN = '<!-- BEGIN GENERATED: custom-role-permissions -->';
const END = '<!-- END GENERATED: custom-role-permissions -->';

function parseArgs(argv) {
  const args = {};
  for (let i = 2; i < argv.length; i += 1) {
    const arg = argv[i];
    if (arg === '--check') args.check = true;
    else if (arg === '--metadata') args.metadata = argv[(i += 1)];
    else if (arg === '--doc') args.doc = argv[(i += 1)];
  }
  if (!args.metadata || !args.doc) {
    throw new Error('Usage: --metadata <yaml> --doc <md> [--check]');
  }
  return args;
}

// Minimal parser for the fixed shape of custom-role-metadata.yaml:
//   groups:
//     - name: <group>
//       relations:
//         - name: <permission>
//           description: <text>
//           relation: <relation>
// Commented (`#`) and blank lines are ignored. Values after the first ": " are
// taken verbatim (descriptions may contain colons and commas) with surrounding
// quotes stripped. A dedicated parser avoids adding a YAML dependency to the
// docs repo for this single, controlled file.
function parseMetadata(yaml) {
  const groups = [];
  let group = null;
  let relation = null;

  for (const rawLine of yaml.split('\n')) {
    const line = rawLine.replace(/\s+$/, '');
    const trimmed = line.trim();
    if (trimmed === '' || trimmed.startsWith('#')) continue;

    const indent = line.length - line.trimStart().length;
    const field = (key) => {
      const prefix = `${key}:`;
      const idx = trimmed.indexOf(prefix);
      if (idx !== 0) return null;
      let value = trimmed.slice(prefix.length).trim();
      if (
        (value.startsWith('"') && value.endsWith('"')) ||
        (value.startsWith("'") && value.endsWith("'"))
      ) {
        value = value.slice(1, -1);
      }
      return value;
    };

    // A group starts at `  - name:` (indent 2); a relation at `      - name:`
    // (indent 6). Distinguish by indentation.
    if (trimmed.startsWith('- name:')) {
      const name = field.call(null, 'name') ?? trimmed.slice('- name:'.length).trim();
      if (indent <= 2) {
        group = { name, relations: [] };
        groups.push(group);
        relation = null;
      } else {
        relation = { name };
        if (group) group.relations.push(relation);
      }
      continue;
    }

    const description = field('description');
    if (description !== null && relation) {
      relation.description = description;
      continue;
    }
  }

  return groups;
}

// Render a GFM table with left-aligned, space-padded columns to match the
// hand-formatted style already in the docs (markdown is excluded from the
// repo's prettier run, so there is no formatter to align it afterwards).
function renderTable(rows) {
  const headers = ['Permission', 'Description'];
  const widths = headers.map((header, col) =>
    Math.max(header.length, ...rows.map((row) => row[col].length)),
  );
  const line = (cells) =>
    `| ${cells.map((cell, col) => cell.padEnd(widths[col])).join(' | ')} |`;
  const separator = `| ${widths
    .map((width) => `:${'-'.repeat(Math.max(width - 1, 2))}`)
    .join(' | ')} |`;
  return [line(headers), separator, ...rows.map(line)].join('\n');
}

function renderTables(groups) {
  const sections = groups
    .filter((group) => group.relations.length > 0)
    .map((group) => {
      const rows = group.relations.map((rel) => [
        rel.name,
        rel.description ?? '',
      ]);
      return `## ${group.name}\n\n${renderTable(rows)}`;
    });
  return sections.join('\n\n');
}

function main() {
  const args = parseArgs(process.argv);
  const groups = parseMetadata(readFileSync(args.metadata, 'utf8'));
  const doc = readFileSync(args.doc, 'utf8');

  const beginIdx = doc.indexOf(BEGIN);
  const endIdx = doc.indexOf(END);
  if (beginIdx === -1 || endIdx === -1 || endIdx < beginIdx) {
    throw new Error(
      `Generation markers not found in ${args.doc}. Expected:\n${BEGIN}\n...\n${END}`,
    );
  }

  const before = doc.slice(0, beginIdx + BEGIN.length);
  const after = doc.slice(endIdx);
  const next = `${before}\n\n${renderTables(groups)}\n\n${after}`;

  if (next === doc) {
    console.log('Custom role permissions reference is up to date.');
    return;
  }

  if (args.check) {
    console.error(
      'Custom role permissions reference is out of date. Run the generator to update it.',
    );
    process.exit(1);
  }

  writeFileSync(args.doc, next);
  console.log(`Updated ${args.doc}.`);
}

main();
