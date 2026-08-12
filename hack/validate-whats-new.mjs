#!/usr/bin/env node

import fs from "node:fs";

const [path, expectedStart, expectedEnd] = process.argv.slice(2);
if (!expectedEnd) {
  throw new Error(`usage: ${process.argv[1]} FILE START END`);
}

const data = JSON.parse(fs.readFileSync(path, "utf8"));
const requiredItemKeys = [
  "product",
  "title",
  "description",
  "url",
  "published",
  "source_prs",
  "featured",
];

if (data.period_start !== expectedStart) {
  throw new Error(`period_start must be ${expectedStart}`);
}
if (data.period_end !== expectedEnd) {
  throw new Error(`period_end must be ${expectedEnd}`);
}
if (!Array.isArray(data.items)) {
  throw new Error("items must be an array");
}
for (const [index, item] of data.items.entries()) {
  const missing = requiredItemKeys.filter((key) => !(key in item));
  if (missing.length > 0) {
    throw new Error(`item ${index + 1} is missing: ${missing.join(", ")}`);
  }

  for (const key of requiredItemKeys.slice(0, 4)) {
    if (typeof item[key] !== "string" || item[key].trim() === "") {
      throw new Error(`item ${index + 1} has an empty ${key}`);
    }
  }

  if (!item.url.startsWith("/") || item.url.startsWith("/manuals/")) {
    throw new Error(`item ${index + 1} url must be an internal published URL`);
  }
  if (
    typeof item.published !== "string" ||
    !/^\d{4}-\d{2}-\d{2}$/.test(item.published) ||
    new Date(`${item.published}T00:00:00Z`).toISOString().slice(0, 10) !==
      item.published
  ) {
    throw new Error(`item ${index + 1} has an invalid published date`);
  }
  if (item.published < data.period_start || item.published > data.period_end) {
    throw new Error(`item ${index + 1} published date is outside the period`);
  }
  if (
    !Array.isArray(item.source_prs) ||
    item.source_prs.length === 0 ||
    !item.source_prs.every((number) => Number.isInteger(number) && number > 0)
  ) {
    throw new Error(
      `item ${index + 1} source_prs must contain pull request numbers`,
    );
  }
  if (typeof item.featured !== "boolean") {
    throw new Error(`item ${index + 1} featured must be a boolean`);
  }
}

const featuredCount = data.items.filter((item) => item.featured).length;
const expectedFeaturedCount = Math.min(5, data.items.length);
if (featuredCount !== expectedFeaturedCount) {
  throw new Error(
    `items must contain ${expectedFeaturedCount} featured highlights`,
  );
}

const publishedDates = data.items.map((item) => item.published);
const sortedDates = [...publishedDates].sort().reverse();
if (publishedDates.join() !== sortedDates.join()) {
  throw new Error("items must be sorted by published date, newest first");
}

const titles = data.items.map((item) => item.title);
if (new Set(titles).size !== titles.length) {
  throw new Error("items must not contain duplicate titles");
}

const urls = data.items.map((item) => item.url);
if (new Set(urls).size !== urls.length) {
  throw new Error("items must not contain duplicate URLs");
}
