import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

const source = await readFile(
  new URL("../../assets/js/quick-install.js", import.meta.url),
  "utf8",
);
const { detectPlatform, quickInstalls } = await import(
  `data:text/javascript;base64,${Buffer.from(source).toString("base64")}`
);

test("desktop platforms select the corresponding Docker product and sbx installer", () => {
  for (const [platform, os, product] of [
    ["MacIntel", "macos", "Docker Desktop"],
    ["Win32", "windows", "Docker Desktop"],
    ["Linux x86_64", "linux", "Docker Engine"],
  ]) {
    assert.equal(detectPlatform(platform), os);
    assert.equal(quickInstalls[os][1].product, product);
    assert.equal(quickInstalls[os][0].product, "Docker Sandboxes");
  }
  for (const os of ["macos", "windows"]) {
    assert.equal(quickInstalls[os][1].command, undefined);
    assert.equal(quickInstalls[os][1].downloads.length, 2);
    for (const download of quickInstalls[os][1].downloads)
      assert.equal(new URL(download.url).hostname, "desktop.docker.com");
  }
  for (const installs of Object.values(quickInstalls)) {
    const engine = installs.at(-1);
    assert.equal(engine.product, "Docker Engine");
    assert.match(engine.command, /curl.*sudo sh/);
  }
  assert.match(quickInstalls.linux[0].label, /Ubuntu/);
});

test("mobile and unknown platforms keep installation links instead of desktop commands", () => {
  assert.equal(detectPlatform("Linux armv8l", "Android"), null);
  assert.equal(detectPlatform("MacIntel", "Macintosh", 5), null);
  assert.equal(detectPlatform("iPhone", "iPhone"), null);
  assert.equal(detectPlatform(""), null);
  assert.equal(detectPlatform("FreeBSD"), null);
});
