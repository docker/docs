import { parse, preprocess, postprocess } from "micromark";
import { decodeString } from "micromark-util-decode-string";
import { parseFragment } from "parse5";

export function linkDestinations(content) {
  const events = postprocess(
    parse()
      .document()
      .write(preprocess()(content, "utf8", true)),
  );
  const destinations = [];
  const stack = [];
  for (const [event, token] of events) {
    if (event === "exit") {
      stack.pop();
      continue;
    }
    stack.push(token.type);
    if (
      ["resourceDestinationString", "definitionDestinationString"].includes(
        token.type,
      )
    ) {
      destinations.push({
        start: token.start.offset,
        end: token.end.offset,
        destination: decodeString(
          content.slice(token.start.offset, token.end.offset),
        ),
        kind: stack.includes("image") ? "images" : "links",
        definition: token.type === "definitionDestinationString",
      });
    }
    if (["htmlFlow", "htmlText"].includes(token.type)) {
      const html = content.slice(token.start.offset, token.end.offset);
      const visit = (node) => {
        for (const attr of node.attrs || []) {
          if (!["href", "src"].includes(attr.name)) continue;
          const location = node.sourceCodeLocation?.attrs?.[attr.name];
          if (!location) continue;
          const attribute = html.slice(
            location.startOffset,
            location.endOffset,
          );
          const prefix = attribute.match(/^[^\s=]+\s*=\s*(["']?)/);
          if (!prefix) continue;
          destinations.push({
            start: token.start.offset + location.startOffset + prefix[0].length,
            end: token.start.offset + location.endOffset - (prefix[1] ? 1 : 0),
            destination: attr.value,
            attribute: attr.name,
          });
        }
        for (const child of node.childNodes || []) visit(child);
      };
      visit(parseFragment(html, { sourceCodeLocationInfo: true }));
    }
  }
  return destinations.sort((a, b) => b.start - a.start);
}
