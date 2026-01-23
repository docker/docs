#!/usr/bin/env node

/**
 * Flattens markdown directory structure and resolves all links to absolute paths.
 *
 * This script:
 * 1. Moves index.md files up one level (ai/model-runner/index.md -> ai/model-runner.md)
 * 2. Fixes _index.md and index.md references in links
 * 3. Strips /manuals/ prefix from paths (Hugo config removes this)
 * 4. Resolves all relative links to absolute HTML paths for RAG ingestion
 *
 * Usage: node flatten-and-resolve.js [public-dir]
 */

const fs = require('fs');
const path = require('path');

const PUBLIC_DIR = path.resolve(process.argv[2] || 'public');

if (!fs.existsSync(PUBLIC_DIR)) {
  console.error(`Error: Directory ${PUBLIC_DIR} does not exist`);
  process.exit(1);
}

/**
 * Recursively find all files matching a predicate
 */
function findFiles(dir, predicate) {
  const results = [];
  const entries = fs.readdirSync(dir, { withFileTypes: true });

  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      results.push(...findFiles(fullPath, predicate));
    } else if (entry.isFile() && predicate(entry.name)) {
      results.push(fullPath);
    }
  }

  return results;
}

/**
 * Step 1: Flatten index.md files
 * Move path/to/section/index.md -> path/to/section.md
 * Before moving, rewrite sibling links (e.g., "get-started.md" -> "section/get-started.md")
 */
function flattenIndexFiles() {
  const indexFiles = findFiles(PUBLIC_DIR, name => name === 'index.md');
  let count = 0;

  for (const file of indexFiles) {
    // Skip root index.md
    if (file === path.join(PUBLIC_DIR, 'index.md')) {
      continue;
    }

    const dir = path.dirname(file);
    const dirname = path.basename(dir);

    // Read content and fix sibling links
    let content = fs.readFileSync(file, 'utf8');

    // Rewrite relative links that don't start with /, ../, or http
    // These are sibling files that will become children after flattening
    content = content.replace(
      /\[([^\]]+)\]\(([a-zA-Z0-9][^):]*)\)/g,
      (match, text, link) => {
        // Skip if it's a URL or starts with special chars
        if (link.startsWith('http://') || link.startsWith('https://') ||
            link.startsWith('#')) {
          return match;
        }
        return `[${text}](${dirname}/${link})`;
      }
    );

    // Also fix reference-style links
    content = content.replace(
      /^\[([^\]]+)\]:\s+([a-zA-Z0-9][^: ]*\.md)$/gm,
      (match, ref, link) => `[${ref}]: ${dirname}/${link}`
    );

    fs.writeFileSync(file, content, 'utf8');

    // Move file up one level
    const parentDir = path.dirname(dir);
    const newPath = path.join(parentDir, `${dirname}.md`);
    fs.renameSync(file, newPath);

    count++;
  }

  console.log(`Flattened ${count} index.md files`);
  return count;
}

/**
 * Step 2: Fix _index.md and index.md references in all files
 * Also strip /manuals/ prefix from paths
 */
function fixIndexReferences() {
  const mdFiles = findFiles(PUBLIC_DIR, name => name.endsWith('.md'));
  let count = 0;

  for (const file of mdFiles) {
    const dir = path.dirname(file);
    const dirname = path.basename(dir);
    const parentDir = path.dirname(dir);
    const parentDirname = path.basename(parentDir);

    let content = fs.readFileSync(file, 'utf8');
    const original = content;

    // Fix path/_index.md or path/index.md -> path.md
    content = content.replace(/([a-zA-Z0-9_/-]+)\/_?index\.md/g, '$1.md');

    // Fix bare _index.md or index.md -> ../dirname.md
    content = content.replace(/_?index\.md/g, `../${dirname}.md`);

    // Fix ../_index.md that became ...md -> ../../parentdirname.md
    if (parentDir !== PUBLIC_DIR) {
      content = content.replace(/\.\.\.md/g, `../../${parentDirname}.md`);
    }

    // Strip /manuals/ prefix (both /manuals/ and manuals/)
    content = content.replace(/\/?manuals\//g, '/');

    if (content !== original) {
      fs.writeFileSync(file, content, 'utf8');
      count++;
    }
  }

  console.log(`Fixed _index.md references in ${count} files`);
  return count;
}

/**
 * Step 3: Resolve all relative links to absolute HTML paths
 */
function resolveLinks() {
  const mdFiles = findFiles(PUBLIC_DIR, name => name.endsWith('.md'));
  let count = 0;

  for (const file of mdFiles) {
    let content = fs.readFileSync(file, 'utf8');
    const original = content;

    // Process inline links: [text](path)
    content = content.replace(/\[([^\]]+)\]\(([^)]+)\)/g, (match, text, link) => {
      const resolved = resolveLinkPath(link, file);
      return `[${text}](${resolved})`;
    });

    // Process reference-style links: [ref]: path
    content = content.replace(/^\[([^\]]+)\]:\s+(.+)$/gm, (match, ref, link) => {
      const resolved = resolveLinkPath(link, file);
      return `[${ref}]: ${resolved}`;
    });

    if (content !== original) {
      fs.writeFileSync(file, content, 'utf8');
      count++;
    }
  }

  console.log(`Resolved links in ${count} files`);
  return count;
}

/**
 * Resolve a link path to absolute HTML format
 */
function resolveLinkPath(linkPath, currentFile) {
  // Skip external URLs and mailto
  if (linkPath.startsWith('http://') || linkPath.startsWith('https://') ||
      linkPath.startsWith('mailto:')) {
    return linkPath;
  }

  // Skip same-page anchors
  if (linkPath.startsWith('#')) {
    return linkPath;
  }

  // Split path and anchor
  const hashIndex = linkPath.indexOf('#');
  const pathPart = hashIndex >= 0 ? linkPath.substring(0, hashIndex) : linkPath;
  const anchorPart = hashIndex >= 0 ? linkPath.substring(hashIndex) : '';

  if (!pathPart) {
    // Just an anchor
    return linkPath;
  }

  // Handle absolute paths - just convert to HTML format
  if (pathPart.startsWith('/')) {
    return toHtmlPath(pathPart) + anchorPart;
  }

  // Resolve relative path to absolute
  const currentDir = path.dirname(currentFile);
  const absolutePath = path.resolve(currentDir, pathPart);
  const relativePath = path.relative(PUBLIC_DIR, absolutePath);

  // Convert to URL path (forward slashes)
  const urlPath = '/' + relativePath.split(path.sep).join('/');

  return toHtmlPath(urlPath) + anchorPart;
}

/**
 * Convert a path to HTML format (strip .md, add trailing /)
 */
function toHtmlPath(urlPath) {
  if (urlPath.endsWith('.md')) {
    return urlPath.slice(0, -3) + '/';
  }
  return urlPath;
}

// Main execution
console.log('Starting markdown flattening and link resolution...');
console.log('');

const flattenCount = flattenIndexFiles();
const fixCount = fixIndexReferences();
const resolveCount = resolveLinks();

console.log('');
console.log('Done!');
console.log(`- Flattened: ${flattenCount} files`);
console.log(`- Fixed references: ${fixCount} files`);
console.log(`- Resolved links: ${resolveCount} files`);
