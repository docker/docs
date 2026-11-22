import assert from "node:assert/strict";
import test from "node:test";
import { renderGordonCard } from "../../assets/js/gordon-cards/render.mjs";

// Record DOM operations and reject HTML insertion while exercising the real parser and renderer.
globalThis.document = {
  createElement(tagName) {
    return {
      tagName,
      children: [],
      textContent: "",
      attributes: {},
      append(...children) {
        this.children.push(...children);
      },
      setAttribute(name, value) {
        this.attributes[name] = value;
      },
      set innerHTML(value) {
        assert.fail("Card text must be inserted through DOM text nodes");
      },
      querySelectorAll(tagName) {
        return descendants(this, tagName);
      },
    };
  },
};

function descendants(node, tagName) {
  const matches = [];
  for (const child of node.children) {
    if (typeof child === "string") continue;
    if (child.tagName === tagName) matches.push(child);
    matches.push(...descendants(child, tagName));
  }
  return matches;
}

function textContent(node) {
  return typeof node === "string"
    ? node
    : node.textContent + node.children.map(textContent).join("");
}

function description(text) {
  return renderGordonCard(`root = Card("Example", [], ${JSON.stringify(text)})`)
    .children[1];
}

test("card links preserve valid URL syntax and sentence punctuation", () => {
  for (const [input, href, text] of [
    [
      "Download https://download.docker.com/linux/static/stable/.",
      "https://download.docker.com/linux/static/stable/",
      "Download https://download.docker.com/linux/static/stable/.",
    ],
    [
      "https://example.com/search?q=O'Reilly",
      "https://example.com/search?q=O%27Reilly",
      "https://example.com/search?q=O'Reilly",
    ],
    [
      "[nested](https://example.com/a_(b_(c)))",
      "https://example.com/a_(b_(c))",
      "nested",
    ],
    [
      '[Docs](https://docs.docker.com/ "Title")',
      "https://docs.docker.com/",
      "Docs",
    ],
    [
      "(see https://en.wikipedia.org/wiki/Docker_(software)).",
      "https://en.wikipedia.org/wiki/Docker_(software)",
      "(see https://en.wikipedia.org/wiki/Docker_(software)).",
    ],
  ]) {
    const node = description(input);
    const links = descendants(node, "a");
    assert.equal(links.length, 1, input);
    assert.equal(links[0].href, href, input);
    assert.equal(textContent(node), text, input);
  }
});

test("links inside emphasis preserve the surrounding literal Markdown", () => {
  for (const marker of ["*", "**", "***", "~", "~~"]) {
    const node = description(
      `${marker}[Docs](https://docs.docker.com/)${marker}`,
    );
    assert.equal(descendants(node, "a")[0].href, "https://docs.docker.com/");
    assert.equal(textContent(node), `${marker}Docs${marker}`);
  }
});

test("card prose never inserts HTML and rejects non-HTTP link schemes", () => {
  for (const input of [
    "<img src=x onerror=alert(1)>",
    "[bad](javascript:alert(1))",
    "[bad](data:text/html,abc)",
    "[bad](ftp://example.com/file)",
    "`curl https://docs.docker.com/`",
    "![image](https://example.com/image.png)",
  ]) {
    const node = description(input);
    assert.equal(descendants(node, "a").length, 0, input);
    assert.equal(textContent(node), input);
  }
  const node = description(
    "[<img src=x onerror=alert(1)>](https://example.com/)",
  );
  assert.equal(descendants(node, "img").length, 0);
  assert.equal(
    descendants(node, "a")[0].textContent,
    "<img src=x onerror=alert(1)>",
  );
});

test("descriptions, step details, callouts and table cells link URLs", () => {
  const url = "https://docs.docker.com/";
  const card = renderGordonCard(`
step = Step("Title ${url}", "Detail ${url}", "curl ${url}")
steps = Steps([step])
callout = Callout("info", "Title ${url}", "Text ${url}")
table = Table(["Header ${url}"], [["Cell ${url}", 4]])
code = Code("curl ${url}", "bash")
followUps = FollowUps(["Explain ${url}"])
root = Card("Title ${url}", [steps, callout, table, code, followUps], "Description ${url}")
`);
  assert.equal(descendants(card, "a").length, 4);
  for (const tag of ["h3", "th", "code", "button"]) {
    for (const node of descendants(card, tag)) {
      assert.equal(descendants(node, "a").length, 0);
    }
  }
  const step = descendants(card, "li")[0];
  assert.equal(descendants(step.children[0], "a").length, 0);
  assert.equal(descendants(step.children[1], "a").length, 1);
  assert.equal(textContent(descendants(card, "td")[1]), "4");
});

test("incomplete links can be rendered repeatedly as an answer streams", () => {
  for (const input of [
    "https://",
    "[Docs](",
    "[Docs](https://docs.docker.com",
    "**[Docs](https://docs.docker.com/)",
    "[Docs](https://docs.docker.com/)",
  ]) {
    for (let repeat = 0; repeat < 2; repeat++) {
      assert.ok(description(input));
    }
  }
});
