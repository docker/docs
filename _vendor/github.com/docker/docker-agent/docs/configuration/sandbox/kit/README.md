# Docker Agent workload kit (v3)

This source kit layers the Docker Agent sandbox image with sbx-native network,
credential, shared-skills, context, and session declarations. `version` is the
example kit's version, not the Docker Agent binary version. Pin the base image
to a release or digest when publishing your own kit.

From the repository root, with a v3-capable sbx:

```console
docker agent run --sandbox \
  --sandbox-kit ./docs/configuration/sandbox/kit default
```

Add v3 mixins with repeated `--kit REF`, and arguments with `--kit-arg KEY=VALUE`.
The workload must keep `docker-agent` on PATH: Docker Agent executes its own
`run` command rather than the workload's default entrypoint.

For cloud execution, use a cloud-enabled sbx and change `--sandbox` to `--cloud`.
Also pass `--model openai/gpt-5.6` when using the built-in `default` agent.
The workspace is remote, not a copy of the current directory. Configure cloud
provider secrets through `sbx --cloud secret`; for this example, store an OpenAI
key with `sbx --cloud secret set openai` (read from stdin). No host API key is
forwarded. The image supplies the non-secret `OPENAI_API_KEY=proxy-managed`
sentinel so Docker Agent can discover OpenAI even when uploaded-kit injection
does not derive an environment variable from the credential capability.
Current sbx uploaded/assembled-image paths drop optional kit credential bindings
and reject required ones. An optional credential declaration alone is **not** a
guarantee of authenticated model access: provision cloud secrets/policies that
work with your workload and backend. The shared-skills capability is optional
because cloud sandboxes cannot mount the host's skills store.

The example permits OpenAI, models.dev, and Docker Hub egress. Extend the policy
for other providers, tool downloads, git hosts, or private registries. An agent
must enable local skills and/or `add_prompt_files: [AGENTS.md]` to consume the
corresponding sbx-provided resources.
