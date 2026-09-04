Choose the authentication method for your agent:

- Codex with a ChatGPT subscription: No secret configuration is needed. `sbx`
  opens the OpenAI sign-in flow on your host before creating the sandbox.
- Claude Code with a Claude subscription: No secret configuration is needed.
  After Claude Code starts, enter `/login` inside the sandbox.
- Codex with an OpenAI API key: Run `sbx secret set openai`.
- Claude Code with an Anthropic API key: Run `sbx secret set anthropic`.
- OpenCode with an API key: Run `sbx secret set` with the service name
  `anthropic`, `openai`, `openrouter`, or `google`.
