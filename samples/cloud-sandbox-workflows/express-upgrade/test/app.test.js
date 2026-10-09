import assert from "node:assert/strict";
import { after, before, test } from "node:test";
import app from "../app.js";

let server;
let base;

before(async () => {
  await new Promise((resolve) => {
    server = app.listen(0, "127.0.0.1", resolve);
  });
  base = `http://127.0.0.1:${server.address().port}`;
});

after(async () => {
  if (server) await new Promise((resolve) => server.close(resolve));
});

test("root route returns the client app", async () => {
  const response = await fetch(`${base}/`);
  assert.equal(response.status, 200);
  assert.match(response.headers.get("content-type"), /text\/html/);
  assert.match(await response.text(), /<h1>Client app<\/h1>/);
});

test("nested client route returns the client app", async () => {
  const response = await fetch(`${base}/projects/42/settings`);
  assert.equal(response.status, 200);
  assert.match(await response.text(), /<h1>Client app<\/h1>/);
});

test("API route returns JSON before the fallback", async () => {
  const response = await fetch(`${base}/api/status`);
  assert.equal(response.status, 200);
  assert.match(response.headers.get("content-type"), /application\/json/);
  assert.deepEqual(await response.json(), { status: "ok" });
});
