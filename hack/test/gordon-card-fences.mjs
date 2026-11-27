import assert from "node:assert/strict";
import test from "node:test";
import { separateCardFences } from "../../assets/js/gordon-cards/fences.mjs";

const card = "```openui-lang\nroot = Card()\n```";

test("a card fence after a sentence moves to its own line", () => {
  assert.equal(
    separateCardFences(`Requires KVM support.${card}\n\nSources:`),
    `Requires KVM support.\n\n${card}\n\nSources:`,
  );
  assert.equal(
    separateCardFences(`Requires KVM support. \t${card}`),
    `Requires KVM support.\n\n${card}`,
  );
});

test("a card fence already at the start of a line is unchanged", () => {
  for (const text of [
    card,
    `Requires KVM support.\n${card}`,
    `Requires KVM support.\n\n${card}`,
    `1. Step\n   ${card}`,
  ]) {
    assert.equal(separateCardFences(text), text);
  }
});

test("text without a card fence is unchanged", () => {
  const text = "Run `docker ps`.\n\n```bash\ndocker ps\n```";
  assert.equal(separateCardFences(text), text);
});

test("a longer backtick run is not split", () => {
  const text = "````openui-lang\nroot = Card()\n````";
  assert.equal(separateCardFences(text), text);
});
