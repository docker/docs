import assert from "node:assert/strict";
import test from "node:test";

import { detectPlatform } from "../../assets/js/quick-install.js";

test("desktop platforms select the corresponding operating system", () => {
  for (const [platform, os] of [
    ["MacIntel", "macos"],
    ["Win32", "windows"],
    ["Linux x86_64", "linux"],
  ])
    assert.equal(detectPlatform(platform), os);
});

test("mobile and unknown platforms keep installation links instead of desktop commands", () => {
  assert.equal(detectPlatform("Linux armv8l", "Android"), null);
  assert.equal(detectPlatform("MacIntel", "Macintosh", 5), null);
  assert.equal(detectPlatform("iPhone", "iPhone"), null);
  assert.equal(detectPlatform(""), null);
  assert.equal(detectPlatform("FreeBSD"), null);
});
