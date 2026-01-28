# Release Notes Generator

Automatically fetch and generate release notes from GitHub repositories.

## Overview

This tool fetches releases from configured GitHub repositories using the `gh`
CLI and generates Hugo-compatible markdown files using Go templates.

## Requirements

- [gh CLI](https://cli.github.com/) - Must be installed and authenticated
- Go 1.23 or later (for building/running)

## Installation

The tool doesn't need to be installed. Run it directly:

```console
$ cd hack/release-notes
$ go run . <repo> [version]
```

Or build it:

```console
$ cd hack/release-notes
$ go build -o release-notes
$ ./release-notes <repo> [version]
```

## Usage

Fetch all releases for a repository:

```console
$ go run . buildx
```

Fetch a specific version:

```console
$ go run . buildx v0.18.0
```

Clean and refetch all releases:

```console
$ go run . --clean buildx
```

Fetch all configured repositories:

```console
$ for repo in buildx compose buildkit; do go run . $repo; done
```

## Configuration

Edit `config.json` to add or modify repositories. Each entry includes:

- `owner` - GitHub repository owner
- `repo` - GitHub repository name
- `content_path` - Where to generate files (relative to project root)
- `title_prefix` - Product name for titles
- `description_template` - SEO description (use `{{version}}` placeholder)
- `keywords_base` - Base keywords for SEO
- `fetch_limit` - Maximum releases to fetch (default: 20)
- `include_prereleases` - Include pre-releases (default: true)

## Template

The markdown template is in `templates/release-note.md.tmpl`. It uses Go's
`text/template` syntax (same as Hugo).

Edit the template to customize the generated markdown structure, headings, or
formatting.

## Generated Files

The tool generates:

- Individual release notes: `{version}.md` in the configured content path
- Section index: `_index.md` (if it doesn't exist)

Each release note includes:

- Hugo front matter with metadata
- Pre-release badge (if applicable)
- Release body from GitHub
- Downloads section with binaries table
- Checksums section
- Link to view on GitHub

## Automation

### GitHub Actions

Create a workflow to automatically fetch new releases:

```yaml
name: Update Release Notes

on:
  schedule:
    - cron: '0 12 * * *'  # Daily at noon
  workflow_dispatch:      # Manual trigger

jobs:
  update-release-notes:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Fetch release notes
        run: |
          cd hack/release-notes
          for repo in buildx compose buildkit; do
            go run . $repo
          done

      - name: Create Pull Request
        uses: peter-evans/create-pull-request@v5
        with:
          title: "Update release notes"
          body: "Automated update of release notes from GitHub releases"
          branch: release-notes-update
          commit-message: "docs: update release notes"
```

## Development

Project structure:

```
hack/release-notes/
├── main.go              # Main program logic
├── config.json          # Repository configuration
├── templates/
│   └── release-note.md.tmpl  # Markdown template
├── go.mod               # Go module definition
└── README.md            # This file
```

The program:

1. Reads `config.json` to get repository configurations
2. Uses `gh` CLI to fetch release data as JSON
3. Parses the JSON into Go structs
4. Processes assets into categories (binaries, checksums, metadata)
5. Executes the markdown template with the data
6. Writes the generated markdown to the content directory

## Updating Existing Content

The tool skips files that already exist unless you use `--clean`. To update a
single release:

```console
$ go run . buildx v0.18.0
```

This overwrites the existing file with fresh content from GitHub.
