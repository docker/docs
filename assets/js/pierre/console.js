import shell from "shiki/langs/shellscript.mjs";

// Console transcripts contain output as well as commands. Only prompt-prefixed
// commands and their backslash continuations should be parsed as shell syntax.
export default [
  ...shell,
  {
    name: "docs-console",
    scopeName: "text.docs-console",
    patterns: [
      {
        begin: "^\\s*([$>]) ",
        beginCaptures: { 1: { name: "keyword.control.prompt" } },
        end: "(?<!\\\\)$",
        patterns: [{ include: "source.shell" }],
      },
    ],
  },
];
