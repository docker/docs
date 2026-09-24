---
title: "Get an SSH certificate"
linkTitle: "Get an SSH certificate"
description: "Have the service sign your SSH public key for one sandbox and return the host, port, username, and host keys your SSH client needs."
keywords: "cloud sandboxes, sandboxes api, get an ssh certificate"
weight: 506
params:
  sidebar:
    group: "Connect and configure"
---

Request a short-lived SSH certificate for access to a sandbox. Keep the private key on your machine; send only the public key to Docker.

You need an [authenticated client](connect-to-cloud-with-a-bearer-token.md), a sandbox name, and an SSH public key. Generate a key pair with your SSH tooling if you do not have one.

## Request the certificate {#1-request-the-certificate}

Pass the sandbox name and public key. Save the returned certificate using your SSH client's expected certificate-file format, alongside the matching private key.

The response supplies the connection details for the sandbox. Use them rather than guessing a hostname or port. The example returns those details; it does not start an SSH process.

A certificate expires. Request a new one when needed, and keep the private key and certificate out of your repository. For programmatic command execution without an SSH client, use [processes](run-your-first-command.md).

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.ssh.issueCertificate({
  publicKey,
  ttlMs: lifeSeconds * 1_000,
});
```

<details>
<summary>Complete TypeScript example: ssh/certificate.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function issueSSHCert(
  client: Sandboxes,
  sandboxName: string,
  publicKey: string,
  lifeSeconds: number,
) {
  const sandbox = await client.get(sandboxName);
  return sandbox.ssh.issueCertificate({
    publicKey,
    ttlMs: lifeSeconds * 1_000,
  });
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
